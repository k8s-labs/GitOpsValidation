package validation

import (
	"context"
	"testing"

	"k8s.io/client-go/kubernetes/fake"
)

func TestRetrieveFluxSource_NotImplemented(t *testing.T) {
	client := fake.NewSimpleClientset()
	_, err := RetrieveFluxSource(context.Background(), client, "default", "my-source")
	if err == nil {
		t.Fatal("expected error for not implemented, got nil")
	}
}
