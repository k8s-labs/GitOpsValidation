package main

import (
	   "flag"
	   "fmt"
	   "os"
	   "net/http"

	   "gov/internal/config"
	   "gov/internal/logging"
	   "gov/internal/validation"
	   "gov/internal/api"
	   "gov/internal/version"

	   "go.uber.org/zap"
	   "k8s.io/client-go/kubernetes"
	   "k8s.io/client-go/rest"
	   "k8s.io/client-go/tools/clientcmd"
	   "context"
	   "k8s.io/client-go/dynamic"
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
	
	select {} // Block forever so the HTTP server stays up

	return
	
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

	// Create dynamic client for CRDs
	dynClient, err := dynamic.NewForConfig(kubeConfig)
	if err != nil {
		logging.Logger.Error("Failed to create dynamic client", zap.Error(err))
		os.Exit(1)
	}

	// Retrieve and populate config from Flux Source
	if err := validation.PopulateConfigFromFluxSource(context.Background(), dynClient, cfg); err != nil {
		logging.Logger.Error("Failed to retrieve and populate config from Flux Source", zap.Error(err))
		os.Exit(1)
	}
	logging.Logger.Info("Populated config from Flux Source", zap.String("repo", cfg.Repo), zap.String("branch", cfg.Branch))
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
