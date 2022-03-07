package k8s

import (
	"fmt"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"

	v1 "github.com/allinbits/starport-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// InClusterConfig should be called from pods running inside a cluster. It returns a Config struct that can be used with
// NewClient().
func InClusterConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}

// KubectlConfig returns a Config struct that can be used with NewClient() using the current context of kubectl.
func KubectlConfig() (*rest.Config, error) {
	home := homedir.HomeDir()
	if home == "" {
		return nil, fmt.Errorf("kubernetes homedir empty")
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	// use the current context in kubeconfig
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

// NewClient returns a Kubernetes API client using the provided Config struct.
func NewClient(config *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()

	if err := v1.SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("cannot add starport operator schemas, %w", err)
	}

	if err := corev1.SchemeBuilder.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("cannot add core schemas, %w", err)
	}

	c, err := client.New(config, client.Options{
		Scheme: scheme,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create client, %w", err)
	}

	return c, nil
}
