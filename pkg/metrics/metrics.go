package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// controller
var (
	TasksSubmitted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orchestrator_tasks_submitted_total",
		Help: "Total number of tasks submitted to the controller",
	})

	TasksPublishErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "orchestrator_task_publish_errors_total",
		Help: "Total number of Kafka publish failures",
	})
)

// worker
var (
	TasksProcessed = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "orchestrator_tasks_processed_total",
		Help: "Total tasks processed by workers",
	}, []string{"status"})

	TaskDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "orchestrator_task_duration_seconds",
		Help:    "Time taked to execute a task end-to-end",
		Buckets: prometheus.DefBuckets,
	})
)
