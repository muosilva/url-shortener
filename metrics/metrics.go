package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var URLsCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "url_shortener_urls_created_total",
	Help: "Total de URLs encurtadas criadas.",
})
