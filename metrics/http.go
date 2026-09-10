package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		HTTPRequestsInFlight.Inc()
		defer HTTPRequestsInFlight.Dec()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		status := ww.Status()
		if status == 0 {
			status = http.StatusOK
		}

		route := "unknown"
		if routeContext := chi.RouteContext(r.Context()); routeContext != nil {
			if routePattern := routeContext.RoutePattern(); routePattern != "" {
				route = routePattern
			}
		}

		statusCode := strconv.Itoa(status)
		duration := time.Since(start).Seconds()

		HTTPRequestTotal.WithLabelValues(r.Method, route, statusCode).Inc()
		HTTPRequestDuration.WithLabelValues(r.Method, route, statusCode).Observe(duration)
	})
}
