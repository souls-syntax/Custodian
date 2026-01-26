package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/souls-syntax/Templates/internal/cache"
	"github.com/souls-syntax/Templates/internal/database"
	"github.com/souls-syntax/Templates/internal/handlers"
	"github.com/souls-syntax/Templates/internal/metrics"
	"github.com/souls-syntax/Templates/internal/service"
)

func main() {

	prometheus.MustRegister(
		metrics.VerifyRequestTotal,
		metrics.VerifyLatencyMs,
		metrics.CacheHitsTotal,
		metrics.CacheMissTotal,
		metrics.AsyncJobsEnqueuedTotal,
		metrics.AsyncJobsDroppedTotal,
		metrics.RaceWinnerTotal,
	)

	godotenv.Load(".env")

	portString := os.Getenv("PORT")
	redisUrl := os.Getenv("REDIS_URL")
	bertUrl := os.Getenv("BERT_URL")
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	if redisUrl == "" {
		log.Fatal("Redis url not set")
	}
	opt, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Setting up redis
	client := redis.NewClient(opt)
	myCache := cache.NewRedisCache(client)

	// Setting up BERT
	bertClient := service.NewBertClient(bertUrl)

	// Setting up DB
	store, err := database.NewStore(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	store.Init()

	// Setting up LLM
	llmClient := service.NewLlmClient(bertUrl)

	// Setting up the Worker Queue
	myWorker := service.NewAsyncProcessor(llmClient, store, myCache)

	myVerifier := service.NewVerifier(myCache, bertClient, store, myWorker)

	apiCfg := &handlers.ApiConfig{
		Verifier: myVerifier,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlers.HandlerReadiness)
	v1Router.Get("/err", handlers.HandlerErr)

	//End Exposed for prometheous
	v1Router.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})

	// Exposed End for verification endpoint
	v1Router.Post("/verify", apiCfg.HandlerVerify)

	// Mounting a route for API control
	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server Starting on port %v", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
