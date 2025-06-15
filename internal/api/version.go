package api

import (
	"fmt"
	"net/http"
	"gov/internal/version"
	"gov/internal/logging"
	"go.uber.org/zap"
)

// VersionHandler responds with 200 OK and the version string.
func VersionHandler(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Info(r.Method,
		zap.String("path", "/version"),
		zap.String("remote", r.RemoteAddr),
		zap.String("userAgent", r.UserAgent()),
	)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "%s", version.Version)
}
