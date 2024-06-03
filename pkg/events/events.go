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
    GetLicenseEvent(application string, date time.Time) (*v1.Event, error)
    CreateLicenseEvent(application string, date time.Time) error
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


func PrepareLicenseEvent(client EventClient, application string, date time.Time) (*v1.Event, error) {
  valid := date.After(time.Now())
  event, err := client.GetLicenseEvent(application, date)
  if err != nil {
    log.Error("Error getting existing event", "error", err)
    return nil, err
  }
  if event != nil {
    log.Debug("Event already exists")
    if !valid {
      log.Debug("Expired event, incrementing count", "previous", event.Count)
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
      Labels: map[string]string{
        "replicated.com/application": application,
        "replicated.com/expires-at": date.Format(time.DateOnly),
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

func (c *KubernetesEventClient) GetLicenseEvent(application string, date time.Time) (*v1.Event, error) {
    podRef := GetObjectReference()
    listOptions := metav1.ListOptions{
        FieldSelector: getFieldSelector(date),
        LabelSelector: getLabelSelector(application, date),
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


func (c *KubernetesEventClient) CreateLicenseEvent(application string, date time.Time) error {
    event, err := PrepareLicenseEvent(c, application, date)
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

func getFieldSelector(date time.Time) string {
  valid := date.After(time.Now())
  reason := "Valid"
  if !valid {
      reason = "Expired"
  }
  log.Debug("Creating field selector", "involvedObject.name", os.Getenv("POD_NAME"), "involvedObject.namespace", os.Getenv("POD_NAMESPACE"), "reason", reason)
  return fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s,reason=%s", os.Getenv("POD_NAME"), os.Getenv("POD_NAMESPACE"), reason)
}

func getLabelSelector(application string, date time.Time) string {
  log.Debug("Creating label selector", "replicated.com/application", application, "replicated.com/expires-at", date)
  return fmt.Sprintf("replicated.com/application=%s,replicated.com/expires-at=%s", application, date.Format(time.DateOnly))
}
