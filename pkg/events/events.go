package events

import (
  "os"
	"time"

  v1 "k8s.io/api/core/v1"
  types "k8s.io/apimachinery/pkg/types"
  "k8s.io/client-go/kubernetes"
)

type EventClient interface {
    CreateExpiredEvent(application string, date time.Time) error
    CreateValidEvent(application string, date time.Time) error
}

type KubernetesEventClient struct {
     Clientset *kubernetes.Clientset
}

func GetObjectReference() v1.ObjectReference {
    return v1.ObjectReference{
        APIVersion: "v1",
        Kind:       "Pod",
        Name:      os.Getenv("POD_NAME"),
        Namespace: os.Getenv("POD_NAMESPACE"),
        UID:       types.UID(os.Getenv("POD_UID")),
    }
}

func GetEventSource() v1.EventSource {
    return v1.EventSource{
      Component: "replicated",
    }
}
 
func NewKubernetesEventClient() (*KubernetesEventClient, error) {
    config, err := GetKubernetesConfig()
    if err != nil {
        return nil, err
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
      return nil, err
    }
    return &KubernetesEventClient{Clientset: clientset}, nil
}
