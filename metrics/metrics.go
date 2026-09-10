package metrics

import (
	"database/sql"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	ResultSuccess  = "success"
	ResultError    = "error"
	ResultNotFound = "not_found"
	ResultHit      = "hit"
	ResultMiss     = "miss"
)

var (
	HTTPRequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total de requests HTTP recebidas.",
	}, []string{"method", "route", "status"})

	HTTPRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Duracao das requests HTTP.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route", "status"})

	HTTPRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Numero de requests HTTP em andamento.",
	})

	URLsCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "url_shortener_urls_created_total",
		Help: "Total de URLs encurtadas criadas.",
	})

	CreateURLDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "url_shortener_create_url_duration_seconds",
		Help:    "Tempo para criar uma URL encurtada.",
		Buckets: prometheus.DefBuckets,
	}, []string{"result"})

	RedirectsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "url_shortener_redirects_total",
		Help: "Total de tentativas de redirecionamento.",
	}, []string{"result"})

	URLLookupTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "url_shortener_url_lookup_total",
		Help: "Total de buscas por URL encurtada.",
	}, []string{"result"})

	CacheOperationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "url_shortener_cache_operations_total",
		Help: "Total de operacoes no cache.",
	}, []string{"operation", "result"})

	CacheOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "url_shortener_cache_operation_duration_seconds",
		Help:    "Duracao das operacoes no cache.",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation", "result"})

	CacheHitsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "url_shortener_cache_hits_total",
		Help: "Total de hits no cache.",
	})

	CacheMissesTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "url_shortener_cache_misses_total",
		Help: "Total de misses no cache.",
	})

	DBOperationsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "url_shortener_db_operations_total",
		Help: "Total de operacoes no banco de dados.",
	}, []string{"operation", "result"})

	DBOperationDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "url_shortener_db_operation_duration_seconds",
		Help:    "Duracao das operacoes no banco de dados.",
		Buckets: prometheus.DefBuckets,
	}, []string{"operation", "result"})
)

var registerDBStatsOnce sync.Once

func ObserveCreateURL(result string, duration time.Duration) {
	CreateURLDuration.WithLabelValues(result).Observe(duration.Seconds())
}

func ObserveCacheOperation(operation string, result string, duration time.Duration) {
	CacheOperationsTotal.WithLabelValues(operation, result).Inc()
	CacheOperationDuration.WithLabelValues(operation, result).Observe(duration.Seconds())
}

func ObserveDBOperation(operation string, result string, duration time.Duration) {
	DBOperationsTotal.WithLabelValues(operation, result).Inc()
	DBOperationDuration.WithLabelValues(operation, result).Observe(duration.Seconds())
}

func RegisterDBStats(db *sql.DB) {
	registerDBStatsOnce.Do(func() {
		prometheus.MustRegister(
			prometheus.NewGaugeFunc(prometheus.GaugeOpts{
				Name: "db_open_connections",
				Help: "Numero de conexoes abertas no pool database/sql.",
			}, func() float64 {
				return float64(db.Stats().OpenConnections)
			}),
			prometheus.NewGaugeFunc(prometheus.GaugeOpts{
				Name: "db_in_use_connections",
				Help: "Numero de conexoes em uso no pool database/sql.",
			}, func() float64 {
				return float64(db.Stats().InUse)
			}),
			prometheus.NewGaugeFunc(prometheus.GaugeOpts{
				Name: "db_idle_connections",
				Help: "Numero de conexoes ociosas no pool database/sql.",
			}, func() float64 {
				return float64(db.Stats().Idle)
			}),
			prometheus.NewCounterFunc(prometheus.CounterOpts{
				Name: "db_wait_count_total",
				Help: "Total de vezes que uma operacao aguardou uma conexao no pool database/sql.",
			}, func() float64 {
				return float64(db.Stats().WaitCount)
			}),
			prometheus.NewCounterFunc(prometheus.CounterOpts{
				Name: "db_wait_duration_seconds_total",
				Help: "Tempo total aguardando conexoes no pool database/sql.",
			}, func() float64 {
				return db.Stats().WaitDuration.Seconds()
			}),
		)
	})
}
