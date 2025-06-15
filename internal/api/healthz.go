package api

import (
	"net/http"
	"gov/internal/logging"
	"go.uber.org/zap"
)

// HealthzHandler responds with 200 OK and "pass" for health checks.
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info("/healthz request",
		zap.String("method", r.Method),
		zap.String("remote", r.RemoteAddr),
		zap.String("userAgent", r.UserAgent()),
	)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pass"))
}
