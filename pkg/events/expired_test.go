package events

import (
    "fmt"
    "strings"
    "time"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestExpiredEvent(t *testing.T) {
    setupMockEnvironment()
    podRef := GetObjectReference()
    application := "Slackernews"
    past := time.Now().Add(-24 * time.Hour)
  
    client := new(MockEventClient) 
    err := client.CreateExpiredEvent(application, past)
    assert.NoError(t, err)
    assert.Len(t, client.EventsCreated, 1)

    event := client.EventsCreated[0]
    assert.Equal(t, fmt.Sprintf("%s-license-expired.", strings.ToLower(application)), event.ObjectMeta.GenerateName)
    assert.Equal(t, podRef.Namespace, event.ObjectMeta.Namespace)
    assert.Equal(t, "Warning", event.Type)
    assert.Equal(t, "Expired", event.Reason)

    assert.Equal(t, fmt.Sprintf("%s license is not valid, expired %v", application, past), event.Message)

    assert.Equal(t, podRef.APIVersion, event.InvolvedObject.APIVersion)
    assert.Equal(t, podRef.Kind, event.InvolvedObject.Kind)
    assert.Equal(t, podRef.Name, event.InvolvedObject.Name)
    assert.Equal(t, podRef.Namespace, event.InvolvedObject.Namespace)
    assert.Equal(t, podRef.UID, event.InvolvedObject.UID)

    assert.NotEmpty(t, event.FirstTimestamp)
    assert.Equal(t, "replicated", event.Source.Component)
}
