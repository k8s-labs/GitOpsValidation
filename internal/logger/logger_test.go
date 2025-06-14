package logger

import (
	"testing"
)

func TestInfoLog(t *testing.T) {
	Info("test info", map[string]any{"foo": "bar"})
}

func TestWarnLog(t *testing.T) {
	Warn("test warn", nil)
}

func TestErrorLog(t *testing.T) {
	Error("test error", map[string]any{"err": "fail"})
}
