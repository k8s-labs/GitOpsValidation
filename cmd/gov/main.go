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

		getConfigFromK8s(cfg)

		logging.Logger.Info("config",
			zap.String("repo", cfg.Repo),
			zap.String("user", cfg.UserName),
			zap.String("PAT", cfg.Password),
			zap.String("path", cfg.Path),
		)
	}
}

func getConfigFromK8s(cfg *config.Config) {
	var restConfig *rest.Config
	var err error

	if isInCluster() {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			logging.Logger.Error("Failed to create in-cluster config", zap.Error(err))
			os.Exit(1)
		}
	} else {
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
		restConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			logging.Logger.Error("Failed to build kubeconfig", zap.Error(err))
			os.Exit(1)
		}
	}

	dynClient, err := dynamic.NewForConfig(restConfig)
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

// isInCluster checks if the application is running inside a Kubernetes cluster
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

func ValidateNamespaceExists(namespace string) error {
	var clientset *kubernetes.Clientset

	if isInCluster() {
		logging.Logger.Info("Running in Kubernetes cluster")
		config, err := rest.InClusterConfig()
		if err != nil {
			return fmt.Errorf("failed to create in-cluster config: %w", err)
		}
		//dynClient, err := dynamic.NewForConfig(config)
		// if err != nil {
		// 	return fmt.Errorf("failed to create dynamic client: %w", err)
		// }
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("failed to create kubernetes clientset: %w", err)
		}
	} else {
		logging.Logger.Info("Running outside Kubernetes cluster, using kubeconfig")
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return fmt.Errorf("failed to build kubeconfig: %w", err)
		}
		// dynClient, err := dynamic.NewForConfig(config)
		// if err != nil {
		// 	return fmt.Errorf("failed to create dynamic client: %w", err)
		// }
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return fmt.Errorf("failed to create kubernetes clientset: %w", err)
		}
	}

	if err := validation.ValidateNamespaceExists(clientset, namespace); err != nil {
		return fmt.Errorf("namespace validation failed: %w", err)
	}

	return nil
}
