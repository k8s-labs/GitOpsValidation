package validation

import (
	"context"
	"fmt"
	"encoding/base64"
	"gov/internal/config"
	"gov/internal/logging"

	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// PopulateConfigFromFluxSource retrieves the Flux Source and fills config fields.
func PopulateConfigFromFluxSource(ctx context.Context, dynClient dynamic.Interface, cfg *config.Config) error {
	source, err := RetrieveFluxSourceUnstructured(ctx, dynClient, cfg.Namespace, cfg.Source)
	if err != nil {
		logging.Logger.Error("Failed to retrieve Flux Source", zap.Error(err))
		return err
	}
	// Extract fields from the unstructured object
	if spec, ok := source.Object["spec"].(map[string]interface{}); ok {
		if url, ok := spec["url"].(string); ok {
			cfg.Repo = url
		}
		if branch, ok := spec["ref"].(map[string]interface{}); ok {
			if b, ok := branch["branch"].(string); ok {
				cfg.Branch = b
			}
		}
		if secretRef, ok := spec["secretRef"].(map[string]interface{}); ok {
			if secretName, ok := secretRef["name"].(string); ok {
				logging.Logger.Info("Using secret for Flux Source", zap.String("secretName", secretName))

				if err := ExtractGitCredentialsFromSecret(ctx, dynClient, cfg, secretName); err != nil {
					return fmt.Errorf("failed to extract git credentials from secret %s: %w", secretName, err)
				}

				logging.Logger.Info(cfg.UserName, zap.String("pat", cfg.Password))
			}
		}
	}
	return nil
}

// ExtractGitCredentialsFromSecret retrieves the username and password from a Kubernetes secret
func ExtractGitCredentialsFromSecret(ctx context.Context, dynClient dynamic.Interface, cfg *config.Config, secretName string) (err error) {
	if secretName == "" {
		logging.Logger.Debug("No secret name provided, skipping credentials extraction")
		return nil
	}
	logging.Logger.Debug("Retrieving git credentials from secret", zap.String("namespace", cfg.Namespace), zap.String("secretName", secretName))
	
	secretsGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "secrets",
	}

	secret, err := dynClient.Resource(secretsGVR).Namespace(cfg.Namespace).Get(ctx, secretName, metav1.GetOptions{})

	if err != nil {
		logging.Logger.Error("Failed to retrieve secret", zap.String("secretName", secretName), zap.Error(err))
		return err
	}
	
	// Extract data from the secret
	data, found, err := unstructured.NestedMap(secret.Object, "data")
	if err != nil || !found || data == nil {
		logging.Logger.Error("Secret data not found", zap.String("secretName", secretName))
		return fmt.Errorf("secret data not found in %s", secretName)
	}
	
	// Extract username and password
	usernameEncoded, foundUsername := data["username"].(string)
	passwordEncoded, foundPassword := data["password"].(string)
	
	if !foundUsername || !foundPassword {
		logging.Logger.Warn("Missing credentials in secret", zap.Bool("hasUsername", foundUsername), zap.Bool("hasPassword", foundPassword))
		return fmt.Errorf("missing username or password in secret %s", secretName)
	}
	
	// Decode from base64
	usernameBytes, err := base64.StdEncoding.DecodeString(usernameEncoded)
	if err != nil {
		return fmt.Errorf("failed to decode username: %w", err)
	}
	
	passwordBytes, err := base64.StdEncoding.DecodeString(passwordEncoded)
	if err != nil {
		return fmt.Errorf("failed to decode password: %w", err)
	}

	cfg.UserName = string(usernameBytes)
	cfg.Password = string(passwordBytes)

	return nil
}
