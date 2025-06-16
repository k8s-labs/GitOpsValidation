package main

import (
	"fmt"
	"os"
	"net/http"
	"time"
	"context"

	"gov/internal/config"
	"gov/internal/logging"
	"gov/internal/validation"
	"gov/internal/api"

	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
)

func main() {
	cfg := config.LoadConfig()

	if err := logging.InitLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logging.SyncLogger()

	logging.Logger.Info("gov starting",
		zap.String("namespace", cfg.Namespace),
		zap.String("source", cfg.Source),
		zap.String("kustomization", cfg.Kustomization),
		zap.Int("sleep", cfg.Sleep),
		zap.Bool("daemonMode", cfg.DaemonMode),
	)

	getConfigFromK8s(cfg)

	logging.Logger.Info("config",
		zap.String("repo", cfg.Repo),
		zap.String("user", cfg.UserName),
		zap.String("PAT", cfg.Password),
		zap.String("path", cfg.Path),
	)

	if cfg.DaemonMode {
		// Start HTTP server for /healthz and /version endpoints
		mux := http.NewServeMux()
		mux.HandleFunc("/healthz", api.HealthzHandler)
		mux.HandleFunc("/version", api.VersionHandler)

		go func() {
			addr := ":8080"
			logging.Logger.Info("Starting HTTP server", zap.String("addr", addr))
			if err := http.ListenAndServe(addr, mux); err != nil {
				logging.Logger.Fatal("HTTP server failed", zap.Error(err))
			}
		}()
		
		for {
			validate(cfg)

			logging.Logger.Info("Sleeping", zap.Int("seconds", cfg.Sleep))

			// Sleep for cfg.Sleep seconds
			select {
				case <-time.After(time.Duration(cfg.Sleep) * time.Second):
			}
		}
	} else {
		validate(cfg)
	}
}

func GetRestConfig() (*rest.Config, error) {
	restConfig, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
		}
	}

	return restConfig, nil
}

func GetDynamicClient() (*dynamic.DynamicClient, error) {
	restConfig, err := GetRestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get rest config: %w", err)
	}

	dynClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return dynClient, nil
}

func getConfigFromK8s(cfg *config.Config) {
	dynClient, err := GetDynamicClient()
	if err != nil {
		logging.Logger.Error("Failed to create dynamic client", zap.Error(err))
		os.Exit(1)
	}

	if err := validation.PopulateConfigFromFluxSource(context.Background(), dynClient, cfg); err != nil {
		logging.Logger.Error("Failed to populate config from Flux Source", zap.Error(err))
		os.Exit(1)
	}

	if err := validation.PopulateConfigFromFluxKustomization(context.Background(), dynClient, cfg); err != nil {
		logging.Logger.Error("Failed to populate config from Flux Kustomization", zap.Error(err))
		os.Exit(1)
	}
}

func validate(cfg *config.Config) {
	if err := ValidateNamespaceExists(cfg.Namespace); err != nil {
		logging.Logger.Error("Namespace validation failed", zap.Error(err))
	} else {
		logging.Logger.Info("Namespace validation succeeded", zap.String("namespace", cfg.Namespace))
	}
}

func mask(s string) string {
	if len(s) <= 4 {
		return "****"
	}
	return s[:2] + "****" + s[len(s)-2:]
}

func GetClientSet() (*kubernetes.Clientset, error) {
	config, err := GetRestConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	return clientset, nil
}

func ValidateNamespaceExists(namespace string) error {
	clientset, err := GetClientSet()
	if err != nil {
		return fmt.Errorf("failed to create kubernetes clientset: %w", err)
	}

	if err := validation.ValidateNamespaceExists(clientset, namespace); err != nil {
		return fmt.Errorf("namespace validation failed: %w", err)
	}

	if err := validation.ValidateServiceInNamespace(clientset, "default", "kubernetes", "ClusterIP", 443); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := validation.ValidateServiceInNamespace(clientset, namespace, "notification-controller", "ClusterIP", 80); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := validation.ValidateServiceInNamespace(clientset, namespace, "source-controller", "ClusterIP", 80); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := validation.ValidateServiceInNamespace(clientset, namespace, "webhook-receiver", "ClusterIP", 80); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := validation.ValidatePodWithPrefixRunning(clientset, namespace, "helm-controller"); err != nil {
		return fmt.Errorf("pod validation failed: %w", err)
	}

	if err := validation.ValidatePodWithPrefixRunning(clientset, namespace, "kustomize-controller"); err != nil {
		return fmt.Errorf("pod validation failed: %w", err)
	}

	if err := validation.ValidatePodWithPrefixRunning(clientset, namespace, "notification-controller"); err != nil {
		return fmt.Errorf("pod validation failed: %w", err)
	}

	if err := validation.ValidateNamespaceExists(clientset, "kube-system"); err != nil {
		return fmt.Errorf("namespace validation failed: %w", err)
	}

	if err := validation.ValidateServiceInNamespace(clientset, "kube-system", "metrics-server", "ClusterIP", 443); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := validation.ValidateServiceInNamespace(clientset, "kube-system", "kube-dns", "ClusterIP", 53); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := validation.ValidatePodWithPrefixRunning(clientset, "kube-system", "coredns"); err != nil {
		return fmt.Errorf("pod validation failed: %w", err)
	}

	if err := validation.ValidatePodWithPrefixRunning(clientset, "kube-system", "local-path-provisioner"); err != nil {
		return fmt.Errorf("pod validation failed: %w", err)
	}

	if err := validation.ValidatePodWithPrefixRunning(clientset, "kube-system", "metrics-server"); err != nil {
		return fmt.Errorf("pod validation failed: %w", err)
	}

	if err := validation.ValidateNamespaceExists(clientset, "heartbeat"); err != nil {
		return fmt.Errorf("namespace validation failed: %w", err)
	}

	if err := validation.ValidateServiceInNamespace(clientset, "heartbeat", "heartbeat", "ClusterIP", 8080); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := validation.ValidatePodWithPrefixRunning(clientset, "heartbeat", "heartbeat"); err != nil {
		return fmt.Errorf("pod validation failed: %w", err)
	}

	if err := validation.ValidateNamespaceExists(clientset, "timeclock"); err != nil {
		return fmt.Errorf("namespace validation failed: %w", err)
	}

	if err := validation.ValidateServiceInNamespace(clientset, "timeclock", "timeclock", "ClusterIP", 8080); err != nil {
		return fmt.Errorf("service validation failed: %w", err)
	}

	if err := validation.ValidatePodWithPrefixRunning(clientset, "timeclock", "timeclock"); err != nil {
		return fmt.Errorf("pod validation failed: %w", err)
	}

	return nil
}
