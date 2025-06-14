package validation

import (
	"context"
	"fmt"

	"gov/internal/logging"

	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ValidateNamespaceExists checks if the given namespace exists in the cluster
func ValidateNamespaceExists(clientset *kubernetes.Clientset, namespace string) error {
	ctx := context.Background()
	_, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		logging.Logger.Error("Namespace not found or error", zap.String("namespace", namespace), zap.Error(err))
		return fmt.Errorf("namespace '%s' not found or error: %w", namespace, err)
	}
	logging.Logger.Info("Namespace exists", zap.String("namespace", namespace))
	return nil
}
