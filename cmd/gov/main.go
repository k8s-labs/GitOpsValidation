package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"

	"gov/internal/api"
	"gov/internal/config"
	"gov/internal/logging"
	"gov/internal/validation"
	"k8s.io/client-go/kubernetes"
)

const version = "0.0.1"

func main() {
	logging.InitLogger()
	defer logging.SyncLogger()

	cfg := config.LoadConfig()
	clientset, err := validation.GetKubeClientset(logging.Logger)
	if err != nil {
		fatal("Failed to get Kubernetes clientset", err)
	}

	rootCmd := newRootCmd(cfg, clientset)
	if err := rootCmd.Execute(); err != nil {
		fatal("Command failed", err)
	}
}

func newRootCmd(cfg *config.Config, clientset *kubernetes.Clientset) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gov",
		Short: "Kubernetes validation tool",
		Run: func(cmd *cobra.Command, args []string) {
			if cfg.Daemon {
				logging.Logger.Infow("Starting in daemon mode",
					"namespace", cfg.Namespace,
					"source", cfg.Source,
					"kustomization", cfg.Kustomization,
					"sleep", cfg.Sleep,
				)
				api.StartServer(version)
				for {
					validation.ValidateAll(clientset, logging.Logger)
					logging.Logger.Infow("Sleeping", "seconds", cfg.Sleep)
					time.Sleep(time.Duration(cfg.Sleep) * time.Second)
				}
			} else {
				logging.Logger.Infow("Running one-time validation",
					"namespace", cfg.Namespace,
					"source", cfg.Source,
					"kustomization", cfg.Kustomization,
				)
				validation.ValidateAll(clientset, logging.Logger)
			}
		},
	}
	cmd.Flags().BoolVarP(&cfg.Daemon, "daemon", "d", cfg.Daemon, "Run as daemon")
	cmd.Flags().IntVarP(&cfg.Sleep, "sleep", "l", cfg.Sleep, "Sleep seconds between validations (daemon mode)")
	cmd.Flags().StringVarP(&cfg.Namespace, "namespace", "n", cfg.Namespace, "Kubernetes namespace Flux is deployed to")
	cmd.Flags().StringVarP(&cfg.Source, "source", "s", cfg.Source, "The Flux source repo")
	cmd.Flags().StringVarP(&cfg.Kustomization, "kustomization", "k", cfg.Kustomization, "The base Kustomization")
	cmd.Version = version
	return cmd
}

func fatal(msg string, err error) {
	logging.Logger.Errorw(msg, "error", err)
	os.Exit(1)
}
