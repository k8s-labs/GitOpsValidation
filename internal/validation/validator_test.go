package validation

import (
	"go.uber.org/zap"
	"testing"
)

func TestValidators(t *testing.T) {
	logger, _ := zap.NewProduction()
	clientset, err := GetKubeClientset(logger.Sugar())
	if err != nil {
		t.Skip("Skipping test: could not get Kubernetes clientset (need KUBECONFIG)")
	}

	tests := []struct {
		name       string
		validator  ResourceValidator
		shouldPass bool
	}{
		{"Namespace default", NamespaceValidator{Namespace: "default"}, true},
		{"Service default/kubernetes", ServiceValidator{Namespace: "default", Name: "kubernetes", Type: "ClusterIP", Port: 443}, true},
		{"Namespace kube-system", NamespaceValidator{Namespace: "kube-system"}, true},
		{"Service kube-system/kube-dns", ServiceValidator{Namespace: "kube-system", Name: "kube-dns", Type: "ClusterIP", Port: 53}, true},
		{"Service kube-system/metrics-server", ServiceValidator{Namespace: "kube-system", Name: "metrics-server", Type: "ClusterIP", Port: 443}, true},
		{"Pod kube-system/coredns", PodValidator{Namespace: "kube-system", NamePrefix: "coredns"}, true},
		{"Pod kube-system/local-path-provisioner", PodValidator{Namespace: "kube-system", NamePrefix: "local-path-provisioner"}, true},
		{"Pod kube-system/metrics-server", PodValidator{Namespace: "kube-system", NamePrefix: "metrics-server"}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.validator.Validate(clientset, logger.Sugar())
			if tc.shouldPass && err != nil {
				t.Errorf("Expected %s to pass, but got error: %v", tc.name, err)
			}
			if !tc.shouldPass && err == nil {
				t.Errorf("Expected %s to fail, but it passed", tc.name)
			}
		})
	}
}
