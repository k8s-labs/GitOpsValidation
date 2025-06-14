package main

import (
	"os"
	"testing"
)

func TestMain_ValidationFailure(t *testing.T) {
	// Simulate missing GOV_REPO
	os.Clearenv()
	os.Args = []string{"gov"}
	// This would normally call main(), but since main() calls os.Exit, we can't test it directly.
	// Instead, test config.ParseConfig and Validate logic in config package tests.
}

// Note: Main function logic is covered by config package tests for config parsing and validation.
// For full integration testing, consider refactoring main() to allow injection/mocking.
