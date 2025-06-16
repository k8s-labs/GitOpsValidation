package validation

import (
	"context"
	"fmt"

	"gov/internal/logging"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ValidateNamespaceExists checks if the given namespace exists in the cluster
func ValidateNamespaceExists(clientset *kubernetes.Clientset, namespace string) error {
	ctx := context.Background()

	_, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("namespace '%s' not found or error: %w", namespace, err)
	}

	return nil
}

// ValidateServiceInNamespace checks that the namespace exists, the service exists, and the service has the specified type and port.
// Logs errors and validation results. Returns an error if validation fails.
func ValidateServiceInNamespace(clientset *kubernetes.Clientset, namespace, serviceName, serviceType string, port int32) error {
	ctx := context.Background()
	// Check namespace
	_, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		msg := fmt.Sprintf("Namespace '%s' not found: %v", namespace, err)
		logging.Logger.Error(msg)
		return fmt.Errorf(msg)
	}

	// Check service
	svc, err := clientset.CoreV1().Services(namespace).Get(ctx, serviceName, metav1.GetOptions{})
	if err != nil {
		msg := fmt.Sprintf("Service '%s' not found in namespace '%s': %v", serviceName, namespace, err)
		logging.Logger.Error(msg)
		return fmt.Errorf(msg)
	}

	// Check service type
	actualType := string(svc.Spec.Type)
	if actualType != serviceType {
		msg := fmt.Sprintf("Service '%s' in namespace '%s' has type '%s', expected '%s'", serviceName, namespace, actualType, serviceType)
		logging.Logger.Error(msg)
		return fmt.Errorf(msg)
	}

	// Check port
	foundPort := false
	for _, p := range svc.Spec.Ports {
		if p.Port == port {
			foundPort = true
			break
		}
	}
	if !foundPort {
		msg := fmt.Sprintf("Service '%s' in namespace '%s' does not have port %d", serviceName, namespace, port)
		logging.Logger.Error(msg)
		return fmt.Errorf(msg)
	}

	msg := fmt.Sprintf("Validation successful: Service '%s' in namespace '%s' exists with type '%s' and port %d", serviceName, namespace, serviceType, port)
	logging.Logger.Info(msg)
	return nil
}

// ValidatePodInNamespace checks that the namespace exists and the pod exists in that namespace.
// Logs errors and validation results. Returns an error if validation fails.
func ValidatePodInNamespace(clientset *kubernetes.Clientset, namespace, podName string) error {
	ctx := context.Background()
	// Check namespace
	_, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		msg := fmt.Sprintf("Namespace '%s' not found: %v", namespace, err)
		logging.Logger.Error(msg)
		return fmt.Errorf(msg)
	}

	// Check pod
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		msg := fmt.Sprintf("Pod '%s' not found in namespace '%s': %v", podName, namespace, err)
		logging.Logger.Error(msg)
		return fmt.Errorf(msg)
	}

	msg := fmt.Sprintf("Validation successful: Pod '%s' exists in namespace '%s' (phase: %s)", podName, namespace, pod.Status.Phase)
	logging.Logger.Info(msg)
	return nil
}

// ValidatePodWithPrefixRunning checks if at least one pod with the given prefix exists and is running in the namespace.
// Logs errors and validation results. Returns an error if validation fails.
func ValidatePodWithPrefixRunning(clientset *kubernetes.Clientset, namespace, podPrefix string) error {
	ctx := context.Background()
	// Check namespace
	_, err := clientset.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{})
	if err != nil {
		msg := fmt.Sprintf("Namespace '%s' not found: %v", namespace, err)
		logging.Logger.Error(msg)
		return fmt.Errorf(msg)
	}

	podList, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		msg := fmt.Sprintf("Failed to list pods in namespace '%s': %v", namespace, err)
		logging.Logger.Error(msg)
		return fmt.Errorf(msg)
	}

	foundRunning := false
	for _, pod := range podList.Items {
		if len(pod.Name) >= len(podPrefix) && pod.Name[:len(podPrefix)] == podPrefix {
			if pod.Status.Phase == "Running" {
				msg := fmt.Sprintf("Validation successful: Pod '%s' with prefix '%s' is running in namespace '%s'", pod.Name, podPrefix, namespace)
				logging.Logger.Info(msg)
				foundRunning = true
				break
			}
		}
	}
	if !foundRunning {
		msg := fmt.Sprintf("No running pod found with prefix '%s' in namespace '%s'", podPrefix, namespace)
		logging.Logger.Error(msg)
		return fmt.Errorf(msg)
	}
	return nil
}
