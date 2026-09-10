package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var URLsCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "url_shortener_urls_created_total",
	Help: "Total de URLs encurtadas criadas.",
})

var CreateURLDuration = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "url_shortener_create_url_duration_seconds",
	Help:    "Tempo para criar uma URL encurtada.",
	Buckets: prometheus.DefBuckets,
})
