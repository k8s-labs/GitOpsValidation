package config

import (
	"os"
	"testing"
)

func TestParseConfig_EnvVars(t *testing.T) {
	os.Setenv("GOV_REPO", "https://github.com/example/repo.git")
	os.Setenv("GOV_USER", "testuser")
	os.Setenv("GOV_PAT", "testpat")
	os.Setenv("GOV_BRANCH", "dev")
	os.Setenv("GOV_PATH", "manifests/")
	os.Setenv("GOV_WAIT", "42")
	defer os.Clearenv()

	cfg := ParseConfig()

	if cfg.RepoURL != "https://github.com/example/repo.git" {
		t.Errorf("expected RepoURL from env, got %s", cfg.RepoURL)
	}
	if cfg.User != "testuser" {
		t.Errorf("expected User from env, got %s", cfg.User)
	}
	if cfg.PAT != "testpat" {
		t.Errorf("expected PAT from env, got %s", cfg.PAT)
	}
	if cfg.Branch != "dev" {
		t.Errorf("expected Branch from env, got %s", cfg.Branch)
	}
	if cfg.Path != "manifests/" {
		t.Errorf("expected Path from env, got %s", cfg.Path)
	}
	if cfg.WaitTime != 42 {
		t.Errorf("expected WaitTime from env, got %d", cfg.WaitTime)
	}
}

func TestConfig_Validate(t *testing.T) {
	cfg := &Config{RepoURL: ""}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for missing RepoURL, got nil")
	}
	cfg.RepoURL = "https://github.com/example/repo.git"
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error for valid RepoURL, got %v", err)
	}
}
