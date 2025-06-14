package validator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseManifests(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.yaml")
	content := []byte(`kind: Deployment
metadata:
  name: test-deploy
  namespace: default
`)
	if err := os.WriteFile(file, content, 0644); err != nil {
		t.Fatalf("failed to write test manifest: %v", err)
	}
	manifests, err := ParseManifests(dir)
	if err != nil {
		t.Fatalf("ParseManifests failed: %v", err)
	}
	if len(manifests) != 1 {
		t.Errorf("expected 1 manifest, got %d", len(manifests))
	}
	if manifests[0].Kind != "Deployment" {
		t.Errorf("expected kind Deployment, got %s", manifests[0].Kind)
	}
	if manifests[0].Metadata.Name != "test-deploy" {
		t.Errorf("expected name test-deploy, got %s", manifests[0].Metadata.Name)
	}
}
