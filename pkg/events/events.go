package events

import (
    "fmt"
    "strings"
    "os"
    "time"
    "context"

    "github.com/charmbracelet/log"

    "k8s.io/client-go/kubernetes"
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    types "k8s.io/apimachinery/pkg/types"

)

type EventClient interface {
    GetLicenseEvent(valid bool, application string, date time.Time) (*v1.Event, error)
    CreateLicenseEvent(valid bool, application string, date time.Time) error
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


func PrepareLicenseEvent(client EventClient, valid bool, application string, date time.Time) (*v1.Event, error) {
  event, err := client.GetLicenseEvent(valid, application, date)
  if err != nil {
    log.Error("Error getting existing event", "error", err)
    return nil, err
  }
  if event != nil {
    log.Debug("Event already exists")
    if !valid {
      log.Debug("Invalid event, incrementing count", "previous", event.Count)
      event.Count++
    }
    return event, nil
  }

  podRef := GetObjectReference()
  eventType := "Normal" 
  reason := "Valid"
  message := fmt.Sprintf("%s license is valid, expires %v", application, date)

  if !valid {
      eventType = "Warning"
      reason = "Expired"
      message = fmt.Sprintf("%s license is not valid, expired %v", application, date)
  } 

  event = &v1.Event{
    ObjectMeta: metav1.ObjectMeta{
      GenerateName: fmt.Sprintf("%s.", strings.ToLower(application)),
      Namespace:   podRef.Namespace,
      Annotations: map[string]string{
        "application": application,
        "expiration": date.Format(time.RFC3339),
      },
    },
    Type:    eventType,
    Reason:  reason,
    Message: message,
    InvolvedObject: podRef,
    FirstTimestamp: metav1.Time{Time: time.Now()},
    Source: GetEventSource(),
    Count: 1,
  }
  return event, nil
}

func (c *KubernetesEventClient) GetLicenseEvent(valid bool, application string, date time.Time) (*v1.Event, error) {
    podRef := GetObjectReference()
    reason := "Valid"
    if !valid {
        reason = "Expired"
    }
    listOptions := metav1.ListOptions{
        FieldSelector: fmt.Sprintf("involvedObject.uid=%s,reason=%s", podRef.UID, reason),
    }
    events, err := c.Clientset.CoreV1().Events(podRef.Namespace).List(context.TODO(), listOptions)
    if err != nil {
        log.Error("Error getting events from Kubernertes", "error", err)
        return nil, err
    }
    if len(events.Items) > 0 {
        log.Debug("Found events", "events_created", len(events.Items))
        return &events.Items[len(events.Items)-1], nil
    }
    return nil, nil
}


func (c *KubernetesEventClient) CreateLicenseEvent(valid bool, application string, date time.Time) error {
    event, err := PrepareLicenseEvent(c, valid, application, date)
    if err != nil {
      log.Error("Error preparing Kubernetes event", "error", err)
      return nil
    }

    if event.Count > 1 {
      log.Debug("Updating existing event", "count", event.Count)
      _, err := c.Clientset.CoreV1().Events(event.ObjectMeta.Namespace).Update(context.TODO(), event, metav1.UpdateOptions{});
      return err
    }

    _, err = c.Clientset.CoreV1().Events(event.ObjectMeta.Namespace).Create(context.TODO(), event, metav1.CreateOptions{});
    return err
}
