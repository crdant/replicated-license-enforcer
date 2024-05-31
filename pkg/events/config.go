package events

import (
    "os"
    "fmt"
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
    fmt.Println("Running in cluster")
    return rest.InClusterConfig()
  } else {
    fmt.Println("Running outside cluster")
	  return clientcmd.BuildConfigFromFlags("", GetPathToKubeConfig())
  }  
}
