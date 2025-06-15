package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"time"

	"gov/internal/config"
	"gov/internal/logging"
	"gov/internal/validation"
	"gov/internal/api"
	"gov/internal/version"

	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
//    "k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

func main() {
	versionFlag := flag.Bool("version", false, "Print version and exit")
	versionShortFlag := flag.Bool("v", false, "Print version and exit (shorthand)")
	flag.Parse()
	if *versionFlag || *versionShortFlag {
		fmt.Println(version.Version)
		os.Exit(0)
	}

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
		if err := ValidateNamespaceExists(cfg.Namespace); err != nil {
			logging.Logger.Error("Namespace validation failed", zap.Error(err))
		} else {
			logging.Logger.Info("Namespace validation succeeded", zap.String("namespace", cfg.Namespace))
		}

		logging.Logger.Info("Sleeping", zap.Int("seconds", cfg.Sleep))

		// Sleep for cfg.Sleep seconds
		select {
			case <-time.After(time.Duration(cfg.Sleep) * time.Second):
		}
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
