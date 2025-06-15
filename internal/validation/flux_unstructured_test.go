package validation

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic/fake"
)

func TestRetrieveFluxSourceUnstructured_NotFound(t *testing.T) {
	dynClient := fake.NewSimpleDynamicClient(unstructured.UnstructuredJSONScheme)
	_, err := RetrieveFluxSourceUnstructured(context.Background(), dynClient, "default", "notfound")
	if err == nil {
		t.Fatal("expected error for missing resource, got nil")
	}
}

func TestRetrieveFluxSourceUnstructured_Found(t *testing.T) {
	gvr := metav1.GroupVersionResource{
		Group:    "source.toolkit.fluxcd.io",
		Version:  "v1beta2",
		Resource: "gitrepositories",
	}
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvr.GroupVersion().WithKind("GitRepository"))
	obj.SetName("my-source")
	obj.SetNamespace("default")
	dynClient := fake.NewSimpleDynamicClient(unstructured.UnstructuredJSONScheme, obj)
	res, err := RetrieveFluxSourceUnstructured(context.Background(), dynClient, "default", "my-source")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.GetName() != "my-source" {
		t.Errorf("expected name 'my-source', got '%s'", res.GetName())
	}
}
