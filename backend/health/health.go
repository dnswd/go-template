// health/handler.go
package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HttpHandler struct {
	db *pgxpool.Pool
}

func NewHandler(db *pgxpool.Pool) *HttpHandler {
	return &HttpHandler{db: db}
}

type HealthResponse struct {
	Status  string            `json:"status"`
	Checks  map[string]string `json:"checks,omitempty"`
	Version string            `json:"version,omitempty"`
}

// Liveness - is the app running?
func (h *HttpHandler) Liveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status: "ok",
	})
}

// Readiness - is the app ready to serve traffic?
func (h *HttpHandler) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	checks := make(map[string]string)
	status := "ok"

	// Check database
	if err := h.db.Ping(ctx); err != nil {
		checks["database"] = "unhealthy: " + err.Error()
		status = "degraded"
	} else {
		checks["database"] = "healthy"
	}

	w.Header().Set("Content-Type", "application/json")
	if status != "ok" {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	json.NewEncoder(w).Encode(HealthResponse{
		Status: status,
		Checks: checks,
	})
}
