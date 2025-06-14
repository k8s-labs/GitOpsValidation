package git

import (
	"testing"
)

func TestChangeToPath_NonExistent(t *testing.T) {
	dir := t.TempDir()
	err := ChangeToPath(dir, "doesnotexist")
	if err == nil {
		t.Error("expected error for non-existent path, got nil")
	}
}

func TestVerifyRepo_NonExistent(t *testing.T) {
	dir := t.TempDir()
	err := VerifyRepo(dir, "https://github.com/example/repo.git")
	if err == nil {
		t.Error("expected error for missing .git/config, got nil")
	}
}

// Note: More comprehensive tests for CloneRepo, CheckoutBranch, and PullLatest would require a test git repo and/or mocking exec.Command.
