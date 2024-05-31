package events

import (
    "os"
    "path/filepath"

  	"k8s.io/client-go/util/homedir"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

func isRunningInCluster() bool {
    // Try to create a Kubernetes client using the in-cluster configuration.
    _, err := rest.InClusterConfig()
    return err == nil
}

func GetPathToKubeConfig() string {
    if os.Getenv("KUBECONFIG") != "" {
      return os.Getenv("KUBECONFIG")
    } else {
      home := homedir.HomeDir()
      return filepath.Join(home, ".kube", "config")
    }
}

func GetKubernetesConfig() (*rest.Config, error) {
  if isRunningInCluster() {
    return rest.InClusterConfig()
  } else {
	  return clientcmd.BuildConfigFromFlags("", GetPathToKubeConfig())
  }  
}
