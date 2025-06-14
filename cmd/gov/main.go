package main

import (
	"fmt"
	"os"

	"gov/internal/config"
	"gov/internal/logging"
	"gov/internal/validation"

	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	if err := logging.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logging.SyncLogger()

	cfg := config.LoadConfig()

	logging.Logger.Info("gov starting",
		zap.String("namespace", cfg.Namespace),
		zap.String("source", cfg.Source),
		zap.String("kustomization", cfg.Kustomization),
	)

	// Create Kubernetes client
	var kubeConfig *rest.Config
	var err error
	if isInCluster() {
		kubeConfig, err = rest.InClusterConfig()
		if err != nil {
			logging.Logger.Error("Failed to get in-cluster config", zap.Error(err))
			os.Exit(1)
		}
	} else {
		kubeconfigPath := os.Getenv("KUBECONFIG")
		if kubeconfigPath == "" {
			kubeconfigPath = clientcmd.RecommendedHomeFile
		}
		kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			logging.Logger.Error("Failed to get kubeconfig from context", zap.Error(err))
			os.Exit(1)
		}
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logging.Logger.Error("Failed to create Kubernetes clientset", zap.Error(err))
		os.Exit(1)
	}

	// Validate namespace exists
	if err := validation.ValidateNamespaceExists(clientset, cfg.Namespace); err != nil {
		logging.Logger.Error("Namespace validation failed", zap.Error(err))
		os.Exit(1)
	}
}

func isInCluster() bool {
	// Kubernetes injects this file in pods
	_, err := os.Stat("/var/run/secrets/kubernetes.io/serviceaccount/token")
	return err == nil
}

func mask(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}
