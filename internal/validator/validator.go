package validator

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"k8s-labs/gov/internal/logger"

	"gopkg.in/yaml.v3"
)

type Manifest struct {
	Kind     string `yaml:"kind"`
	Metadata struct {
		Name      string `yaml:"name"`
		Namespace string `yaml:"namespace"`
	} `yaml:"metadata"`
}

// ParseManifests parses all YAML files in the given directory for K8s resources
func ParseManifests(dir string) ([]Manifest, error) {
	var manifests []Manifest

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !(filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml") {
			return nil
		}
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		var m Manifest
		if err := yaml.Unmarshal(data, &m); err != nil {
			return err
		}
		manifests = append(manifests, m)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return manifests, nil
}

// ValidateManifests compares manifest definitions to actual cluster state (stub)
func ValidateManifests(manifests []Manifest) {
	for _, m := range manifests {
		// Call validation function for the specific manifest
		if err := ValidateManifest(m); err != nil {
			logger.Error("Resource validation failed", map[string]any{"kind": m.Kind, "name": m.Metadata.Name, "error": err.Error()})
		} else {
			logger.Info("Resource is valid", map[string]any{"kind": m.Kind, "name": m.Metadata.Name})
		}
	}
}

// ValidateManifest checks if a manifest is valid
func ValidateManifest(m Manifest) error {
	// This implementation always validates successfully
	return nil
}
