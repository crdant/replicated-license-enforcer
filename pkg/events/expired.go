package events

import (
    "fmt"
    "strings"
    "time"
    "context"
    
    "github.com/charmbracelet/log"

    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PrepareExpiredEvent(client EventClient, application string, date time.Time) (*v1.Event, error) {
  event, err := client.GetExpiredEvent(application, date)
  if err != nil {
    log.Error("Error getting existing event", "error", err)
    return nil, err
  }
  if event != nil {
    log.Debug("Event already exists, incrementing count", "previous",)
    event.Count++
    return event, nil
  }

  podRef := GetObjectReference()
  event = &v1.Event{
    ObjectMeta: metav1.ObjectMeta{
      GenerateName: fmt.Sprintf("%s-license-expired.", strings.ToLower(application)),
      Namespace:   podRef.Namespace,
    },
    Type:    "Warning",
    Reason:  "Expired",
    Message: fmt.Sprintf("%s license is not valid, expired %v", application, date),
    InvolvedObject: podRef,
    FirstTimestamp: metav1.Time{Time: time.Now()},
    Source: GetEventSource(),
    Count: 1,
  }
  return event, nil
}

func (c *KubernetesEventClient) GetExpiredEvent(application string, date time.Time) (*v1.Event, error) {
    podRef := GetObjectReference()
    listOptions := metav1.ListOptions{
        FieldSelector: fmt.Sprintf("involvedObject.uid=%s,reason=%s", podRef.UID, "Expired"),
    }
    events, err := c.Clientset.CoreV1().Events(podRef.Namespace).List(context.TODO(), listOptions)
    if err != nil {
        log.Error("Error getting events from Kubernertes", "error", err)
        return nil, err
    }
    if len(events.Items) > 0 {
        log.Debug("Found events", "count", len(events.Items))
        return &events.Items[len(events.Items)-1], nil
    }
    return nil, nil
}


func (c *KubernetesEventClient) CreateExpiredEvent(application string, date time.Time) error {
    event, err := PrepareExpiredEvent(c, application, date)
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
