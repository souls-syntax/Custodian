package service

import (
	"context"
	"fmt"
	"time"

	"github.com/souls-syntax/Templates/internal/cache"
	"github.com/souls-syntax/Templates/internal/database"
	"github.com/souls-syntax/Templates/internal/metrics"
)

type AnalysisJob struct {
	Hash					 string
	QueryText			 string
}

type AsyncProcessor struct {
	JobQueue chan AnalysisJob
	Llm            *LlmClient
	DB             *database.Store
	Cache          *cache.RedisCache
}

func NewAsyncProcessor(llm *LlmClient, db *database.Store, c *cache.RedisCache) *AsyncProcessor {
	worker := &AsyncProcessor{
		JobQueue: make(chan AnalysisJob, 10000),
		Llm: llm,
		DB: db,
		Cache: c,
	}

	go worker.start()
	
	return worker

}

func (w *AsyncProcessor) start() {
	fmt.Println("🐣 Async Intelligence worker started.") 
	fmt.Println("😴 Waiting for jobs.....")

	for job := range w.JobQueue {
		fmt.Println("📥 Job received")

		ctx := context.Background()
		fmt.Printf("⚙️ Processing Async job: %s....\n",job.Hash[:8])

		decision, err := w.Llm.GetAnalysis(ctx, job.QueryText)
		if err != nil {
			fmt.Printf("LLM Failed to process %s with error %v..\n",job.Hash,err)
		}

		decision.QueryText = job.QueryText
		decision.QueryHash = job.Hash

		w.DB.SaveDecision(decision, "LLM-ASYNC")

		w.Cache.Set(ctx, job.Hash, decision.ToCacheModel(), 2*time.Hour)
		fmt.Printf("Async Upgrade Complete for %s. New Verdict: %s, New Confidence: %v\n", job.Hash, decision.Verdict, decision.Confidence)
	}
}

func (w *AsyncProcessor) Enqueue(hash string, text string) {
	select {
	case w.JobQueue <- AnalysisJob{Hash:hash,QueryText:text}:
			 metrics.AsyncJobsEnqueuedTotal.Inc()
	default:
		fmt.Printf("Job Queue full dropping async analysis query hash: %v",hash)
		metrics.AsyncJobsDroppedTotal.Inc()
	}
}
