package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration parameters for gov
// Parameters: repo URL, user, PAT, branch, path, wait time

type Config struct {
	RepoURL  string // GitOps repo URL (required)
	User     string // User ID for the repo (default: "gitops")
	PAT      string // Personal Access Token (optional)
	Branch   string // Branch to use (default: "main")
	Path     string // Path within the repo (default: "./")
	WaitTime int    // Wait time in seconds between validations (default: 60)
}

// ParseFlags parses command line arguments into a Config struct
func ParseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.RepoURL, "gov-repo", "", "GitOps repo URL (required)")
	flag.StringVar(&config.RepoURL, "r", "", "GitOps repo URL (required, shorthand)")
	flag.StringVar(&config.User, "user", "gitops", "User ID for the repo")
	flag.StringVar(&config.PAT, "pat", "", "Personal Access Token (optional)")
	flag.StringVar(&config.Branch, "branch", "main", "Branch to use")
	flag.StringVar(&config.Path, "path", "./", "Path within the repo")
	flag.IntVar(&config.WaitTime, "wait", 60, "Wait time in seconds between validations")

	flag.Parse()
	return config
}

// ParseConfig parses command line arguments and environment variables into a Config struct
func ParseConfig() *Config {
	config := &Config{}

	// Set defaults from environment variables if present
	if v := os.Getenv("GOV_REPO"); v != "" {
		config.RepoURL = v
	}
	if v := os.Getenv("GOV_USER"); v != "" {
		config.User = v
	}
	if v := os.Getenv("GOV_PAT"); v != "" {
		config.PAT = v
	}
	if v := os.Getenv("GOV_BRANCH"); v != "" {
		config.Branch = v
	}
	if v := os.Getenv("GOV_PATH"); v != "" {
		config.Path = v
	}
	if v := os.Getenv("GOV_WAIT"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			config.WaitTime = i
		}
	}

	// Override with command line flags if provided
	flag.StringVar(&config.RepoURL, "repo", config.RepoURL, "GitOps repo URL (required)")
	flag.StringVar(&config.RepoURL, "r", config.RepoURL, "GitOps repo URL (required, shorthand)")
	flag.StringVar(&config.User, "user", config.User, "User ID for the repo")
	flag.StringVar(&config.PAT, "pat", config.PAT, "Personal Access Token (optional)")
	flag.StringVar(&config.Branch, "branch", config.Branch, "Branch to use")
	flag.StringVar(&config.Path, "path", config.Path, "Path within the repo")
	flag.IntVar(&config.WaitTime, "wait", config.WaitTime, "Wait time in seconds between validations")

	flag.Parse()
	return config
}

// Validate checks that all required parameters are set
func (c *Config) Validate() error {
	if c.RepoURL == "" {
		return fmt.Errorf("missing required parameter: repo (or GOV_REPO)")
	}
	return nil
}
