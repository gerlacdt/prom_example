package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
)

// PromMiddleware type for all http routes
type PromMiddleware struct {
	requestCounter prometheus.Counter
	handler        http.Handler
}

// NewMiddleware creates a new http Prometheus middleware
func NewMiddleware(requestCounter prometheus.Counter, handler http.Handler) *PromMiddleware {
	return &PromMiddleware{requestCounter: requestCounter, handler: handler}
}

func (h *PromMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.requestCounter.Inc()
	h.handler.ServeHTTP(w, r)
}
