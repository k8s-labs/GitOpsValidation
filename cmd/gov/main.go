package main

import (
	"os"
	"time"

	"k8s-labs/gov/internal/config"
	"k8s-labs/gov/internal/git"
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

	// Clone or verify repository
	// Extract repo name from URL (last segment)
	repoName := git.ExtractRepoName(cfg.RepoURL)
	repoDir := "./" + repoName

	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		logger.Info("Cloning repository", map[string]any{"repo": cfg.RepoURL, "branch": cfg.Branch})
		err := git.CloneRepo(git.CloneOptions{
			RepoURL: cfg.RepoURL,
			User:    cfg.User,
			PAT:     cfg.PAT,
			Branch:  cfg.Branch,
			Dir:     repoDir,
		})
		if err != nil {
			logger.Fatal("Failed to clone repository", map[string]any{"error": err.Error()})
		}
	} else {
		logger.Info("Repository directory exists, verifying", map[string]any{"dir": repoDir})
		if err := git.VerifyRepo(repoDir, cfg.RepoURL); err != nil {
			logger.Fatal("Repository verification failed", map[string]any{"error": err.Error()})
		}

		if err := git.CheckoutBranch(cfg.Branch); err != nil {
			logger.Fatal("Failed to checkout branch", map[string]any{"error": err.Error()})
		}

		if err := git.PullLatest(repoDir); err != nil {
			logger.Fatal("Failed to pull latest changes", map[string]any{"error": err.Error()})
		}
	}

	// Change to the specified path within the repo if needed
	if err := git.ChangeToPath(repoDir, cfg.Path); err != nil {
		logger.Fatal("Failed to change to repository path", map[string]any{"error": err.Error()})
	}

	for {
		logger.Info("Pulling latest changes", map[string]any{"repo": cfg.RepoURL})
		if err := git.PullLatest(repoDir); err != nil {
			logger.Error("Failed to pull latest changes", map[string]any{"error": err.Error()})
		}

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
