package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TotalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_requests",
			Help: "Total number of requests",
		},
		[]string{"url"},
	)

	AnswerDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "answer_duration_seconds",
			Help:    "Histogram of response duration for requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"url"},
	)
)

var (
	AmountOfCreatedPvz = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "created_pvz_total",
			Help: "Total number of created pvz",
		},
		[]string{"url"},
	)
	AmountOfCreatedReceptions = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "created_receptions_total",
			Help: "Total number of created receptions",
		},
		[]string{"url"},
	)

	AmountOfAddedProducts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "added_products_total",
			Help: "Total number of added products",
		},
		[]string{"url"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(TotalRequests)
	prometheus.MustRegister(AnswerDuration)
	prometheus.MustRegister(AmountOfCreatedPvz)
	prometheus.MustRegister(AmountOfCreatedReceptions)
	prometheus.MustRegister(AmountOfAddedProducts)
}
