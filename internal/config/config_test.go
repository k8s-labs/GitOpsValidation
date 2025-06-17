package config

import (
	"os"
	"testing"
)

func TestGetEnvBool(t *testing.T) {
	os.Setenv("TEST_BOOL", "true")
	if !getEnvBool("TEST_BOOL", false) {
		t.Error("Expected true from env bool")
	}
	os.Setenv("TEST_BOOL", "false")
	if getEnvBool("TEST_BOOL", true) {
		t.Error("Expected false from env bool")
	}
	os.Unsetenv("TEST_BOOL")
	if !getEnvBool("TEST_BOOL", true) {
		t.Error("Expected default true from env bool")
	}
}

func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	if getEnvInt("TEST_INT", 0) != 42 {
		t.Error("Expected 42 from env int")
	}
	os.Unsetenv("TEST_INT")
	if getEnvInt("TEST_INT", 7) != 7 {
		t.Error("Expected default 7 from env int")
	}
}

func TestGetEnvString(t *testing.T) {
	os.Setenv("TEST_STR", "foo")
	if getEnvString("TEST_STR", "bar") != "foo" {
		t.Error("Expected foo from env string")
	}
	os.Unsetenv("TEST_STR")
	if getEnvString("TEST_STR", "bar") != "bar" {
		t.Error("Expected default bar from env string")
	}
}
