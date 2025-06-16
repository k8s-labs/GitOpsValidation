package validation

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

// RetrieveFluxSourceUnstructured fetches a Flux GitRepository resource in the given namespace using the dynamic client.
func RetrieveFluxSourceUnstructured(ctx context.Context, dynClient dynamic.Interface, namespace, sourceName string) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "source.toolkit.fluxcd.io",
		Version:  "v1",
		Resource: "gitrepositories",
	}
	resource, err := dynClient.Resource(gvr).Namespace(namespace).Get(ctx, sourceName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Flux GitRepository '%s' in namespace '%s': %w", sourceName, namespace, err)
	}
	return resource, nil
}

// RetrieveFluxKustomizationUnstructured fetches a Flux Kustomization resource in the given namespace using the dynamic client.
func RetrieveFluxKustomizationUnstructured(ctx context.Context, dynClient dynamic.Interface, namespace, kustomizationName string) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "kustomize.toolkit.fluxcd.io",
		Version:  "v1",
		Resource: "kustomizations",
	}
	resource, err := dynClient.Resource(gvr).Namespace(namespace).Get(ctx, kustomizationName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get Flux Kustomization '%s' in namespace '%s': %w", kustomizationName, namespace, err)
	}
	return resource, nil
}
