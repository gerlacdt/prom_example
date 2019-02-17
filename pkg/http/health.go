package http

import (
	"encoding/json"
	"net/http"
)

// HealthResponse struct of /health
type HealthResponse struct {
	Status string `json:"status"`
}

// HealthHandler ...
type HealthHandler struct{}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{Status: "OK"}
	_ = json.NewEncoder(w).Encode(resp)
}
