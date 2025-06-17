package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func StartServer(version string) {
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pass")
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, version)
	})

	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("GOV_API_PORT")
	if port == "" {
		port = "8080"
	}
	go http.ListenAndServe(":"+port, nil)
}
