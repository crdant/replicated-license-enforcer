package events

import (
    "fmt"
    "strings"
    "time"
    "context"

    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewExpiredEvent(application string, date time.Time) *v1.Event {
  podRef := GetObjectReference()
  event := &v1.Event{
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
  }

  return event
}
func (c *KubernetesEventClient) CreateExpiredEvent(application string, date time.Time) error {
    event := NewExpiredEvent(application, date)
    if _, err := c.Clientset.CoreV1().Events(event.ObjectMeta.Namespace).Create(context.TODO(), event, metav1.CreateOptions{}); err != nil {
        fmt.Printf("Error creating event: %v", err)
        return err
    }
    return nil
}
