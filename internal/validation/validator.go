package validation

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

type ResourceValidator interface {
	Validate(clientset *kubernetes.Clientset, logger *zap.SugaredLogger) error
	Description() string
}

type NamespaceValidator struct {
	Namespace string
}

func (n NamespaceValidator) Validate(clientset *kubernetes.Clientset, logger *zap.SugaredLogger) error {
	_, err := clientset.CoreV1().Namespaces().Get(context.Background(), n.Namespace, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("namespace '%s' does not exist: %w", n.Namespace, err)
	}
	return nil
}
func (n NamespaceValidator) Description() string {
	return fmt.Sprintf("Namespace: %s", n.Namespace)
}

type ServiceValidator struct {
	Namespace string
	Name      string
	Type      string
	Port      int
}

func (s ServiceValidator) Validate(clientset *kubernetes.Clientset, logger *zap.SugaredLogger) error {
	svc, err := clientset.CoreV1().Services(s.Namespace).Get(context.Background(), s.Name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("service '%s/%s' does not exist: %w", s.Namespace, s.Name, err)
	}
	if string(svc.Spec.Type) != s.Type {
		return fmt.Errorf("service '%s/%s' type mismatch: expected %s, got %s", s.Namespace, s.Name, s.Type, svc.Spec.Type)
	}
	foundPort := false
	for _, p := range svc.Spec.Ports {
		if int(p.Port) == s.Port {
			foundPort = true
			break
		}
	}
	if !foundPort {
		return fmt.Errorf("service '%s/%s' missing port %d", s.Namespace, s.Name, s.Port)
	}
	return nil
}
func (s ServiceValidator) Description() string {
	return fmt.Sprintf("Service: %s/%s type=%s port=%d", s.Namespace, s.Name, s.Type, s.Port)
}

type PodValidator struct {
	Namespace  string
	NamePrefix string
}

func (p PodValidator) Validate(clientset *kubernetes.Clientset, logger *zap.SugaredLogger) error {
	pods, err := clientset.CoreV1().Pods(p.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list pods in '%s': %w", p.Namespace, err)
	}
	found := false
	for _, pod := range pods.Items {
		if len(pod.Name) >= len(p.NamePrefix) && pod.Name[:len(p.NamePrefix)] == p.NamePrefix {
			found = true
			if pod.Status.Phase != "Running" {
				return fmt.Errorf("pod '%s/%s' not running (phase: %s)", p.Namespace, pod.Name, pod.Status.Phase)
			}
		}
	}
	if !found {
		return fmt.Errorf("no pods found with prefix '%s' in namespace '%s'", p.NamePrefix, p.Namespace)
	}
	return nil
}
func (p PodValidator) Description() string {
	return fmt.Sprintf("Pod: %s/%s*", p.Namespace, p.NamePrefix)
}

func ValidateAll(clientset *kubernetes.Clientset, logger *zap.SugaredLogger) {
	validators := []ResourceValidator{
		NamespaceValidator{Namespace: "default"},
		ServiceValidator{Namespace: "default", Name: "kubernetes", Type: "ClusterIP", Port: 443},
		NamespaceValidator{Namespace: "kube-system"},
		ServiceValidator{Namespace: "kube-system", Name: "kube-dns", Type: "ClusterIP", Port: 53},
		ServiceValidator{Namespace: "kube-system", Name: "metrics-server", Type: "ClusterIP", Port: 443},
		PodValidator{Namespace: "kube-system", NamePrefix: "coredns"},
		PodValidator{Namespace: "kube-system", NamePrefix: "local-path-provisioner"},
		PodValidator{Namespace: "kube-system", NamePrefix: "metrics-server"},
		NamespaceValidator{Namespace: "heartbeat"},
		ServiceValidator{Namespace: "heartbeat", Name: "heartbeat", Type: "ClusterIP", Port: 8080},
		PodValidator{Namespace: "heartbeat", NamePrefix: "heartbeat"},
		NamespaceValidator{Namespace: "timeclock"},
		ServiceValidator{Namespace: "timeclock", Name: "timeclock", Type: "ClusterIP", Port: 8080},
		PodValidator{Namespace: "timeclock", NamePrefix: "timeclock"},
	}
	for _, v := range validators {
		err := v.Validate(clientset, logger)
		if err != nil {
			logger.Errorw("Validation failed", "resource", v.Description(), "error", err)
		} else {
			logger.Infow("Validation success", "resource", v.Description())
		}
	}
}

func GetKubeClientset(logger *zap.SugaredLogger) (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = os.ExpandEnv("$HOME/.kube/config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logger.Errorw("Failed to get kubeconfig", "error", err)
			return nil, err
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Errorw("Failed to create Kubernetes client", "error", err)
		return nil, err
	}
	return clientset, nil
}
