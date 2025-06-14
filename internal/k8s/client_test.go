package k8s

import "testing"

func TestHealthCheck(t *testing.T) {
	if !HealthCheck() {
		t.Error("expected HealthCheck to return true")
	}
}
