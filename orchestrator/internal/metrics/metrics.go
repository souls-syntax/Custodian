package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	VerifyRequestTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "template_verify_requests_total",
			Help: "Total number of verify requests",
		},
	)

	VerifyLatencyMs = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "template_verify_latency_ms",
			Help:    "Latency of verify requests",
			Buckets: []float64{5, 10, 25, 50, 100, 250, 500, 1000, 2500, 5000},
		},
	)

	RaceWinnerTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "template_race_winner_total",
			Help: "Which source won the race (db vs bert)",
		},
		[]string{"source"}, // Label key
	)

	CacheHitsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "template_cache_hit",
			Help: "Total cahce hits",
		},
	)

	CacheMissTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "template_cache_miss",
			Help: "Total cahce miss",
		},
	)

	AsyncJobsEnqueuedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "template_async_jobs_enqueued_total",
			Help: "Total async jobs enqueued",
		},
	)

	AsyncJobsDroppedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "template_async_jobs_dropped_total",
			Help: "Total async jobs dropped",
		},
	)
)
