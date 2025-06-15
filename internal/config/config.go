package config

import (
	"flag"
	"fmt"
	"os"
	"gov/internal/version"
)

// Config holds all configuration for the gov application
// Command line/environment: namespace, source, kustomization
// Populated from k8s: repo, userId, pat, branch, path

type Config struct {
	Namespace     string // CLI/env: Kubernetes namespace to validate
	Source        string // CLI/env: Flux Source name
	Kustomization string // CLI/env: Flux Kustomization name
	Sleep         int    // CLI/env: Sleep duration in seconds between validations (default 60)
	Repo          string // Populated from k8s: Git repository URL
	UserName      string // Populated from k8s: GitHub user name
	Password      string // Populated from k8s: GitHub password or token
	Branch        string // Populated from k8s: Git branch to use
	Path          string // Populated from k8s: Path within the repo for manifests
}

func getEnvOrDefault(envKey, def string) string {
	if val, ok := os.LookupEnv(envKey); ok && val != "" {
		return val
	}
	return def
}

// LoadConfig parses CLI flags, then environment variables, then uses defaults
func LoadConfig() *Config {
	var (
		nsFlag   = flag.String("namespace", getEnvOrDefault("GOV_NAMESPACE", "flux-system"), "Kubernetes namespace to validate")
		srcFlag  = flag.String("source", getEnvOrDefault("GOV_SOURCE", "gitops"), "Flux Source name")
		kustFlag = flag.String("kustomization", getEnvOrDefault("GOV_KUSTOMIZATION", "flux-listeners"), "Flux Kustomization name")
		sleepFlag = flag.Int("sleep", 60, "Sleep duration in seconds between validations")
		versionFlag = flag.Bool("version", false, "Print version and exit")
		versionShortFlag = flag.Bool("v", false, "Print version and exit (shorthand)")
	)
	flag.Parse()

	if *versionFlag || *versionShortFlag {
		fmt.Println(version.Version)
		os.Exit(0)
	}

	cfg := &Config{
		Namespace:     *nsFlag,
		Source:        *srcFlag,
		Kustomization: *kustFlag,
		Sleep:         *sleepFlag,

		// The following fields are populated after reading from k8s:
		Repo:   "",
		UserName: "",
		Password: "",
		Branch: "",
		Path:   "",
	}

	return cfg
}
