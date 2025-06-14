package main

import (
	"time"

	"k8s-labs/gov/internal/config"
	"k8s-labs/gov/internal/logger"
	"k8s-labs/gov/internal/validator"
)

func main() {
	cfg := config.ParseConfig()
	if err := cfg.Validate(); err != nil {
		logger.Fatal("Startup validation failed", map[string]any{"error": err.Error()})
	}
	logger.Info("gov application started", map[string]any{
		"repo":     cfg.RepoURL,
		"user":     cfg.User,
		"branch":   cfg.Branch,
		"path":     cfg.Path,
		"waitTime": cfg.WaitTime,
	})
	for {
		manifests, err := validator.ParseManifests(".")
		if err != nil {
			logger.Error("Failed to parse manifests", map[string]any{"error": err.Error()})
		} else {
			validator.ValidateManifests(manifests)
		}
		logger.Info("Sleeping before next validation", map[string]any{"waitTime": cfg.WaitTime})
		time.Sleep(time.Duration(cfg.WaitTime) * time.Second)
	}
}
