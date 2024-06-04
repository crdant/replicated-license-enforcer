package events

import (
    "fmt"
    "strings"
    "time"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestExpiredEvent(t *testing.T) {
    client := NewMockEventClient() 
    podRef := GetObjectReference()
    application := "slackernews-mackerel"

    past := time.Now().Add(-24 * time.Hour)
  
    err := client.CreateLicenseEvent(application, past)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    event, err := client.GetLicenseEvent(application, past)
    assert.NoError(t, err)
    assert.Equal(t, fmt.Sprintf("%s.", strings.ToLower(application)), event.ObjectMeta.GenerateName)
    assert.Equal(t, podRef.Namespace, event.ObjectMeta.Namespace)
    assert.Equal(t, "Warning", event.Type)
    assert.Equal(t, "Expired", event.Reason)
    assert.Equal(t, application, event.ObjectMeta.Labels["replicated.com/application"])
    assert.Equal(t, past.Format(time.DateOnly), event.ObjectMeta.Labels["replicated.com/expires-at"])

    assert.Equal(t, fmt.Sprintf("%s license is not valid, expired %v", application, past), event.Message)

    assert.Equal(t, podRef.APIVersion, event.InvolvedObject.APIVersion)
    assert.Equal(t, podRef.Kind, event.InvolvedObject.Kind)
    assert.Equal(t, podRef.Name, event.InvolvedObject.Name)
    assert.Equal(t, podRef.Namespace, event.InvolvedObject.Namespace)
    assert.Equal(t, podRef.UID, event.InvolvedObject.UID)

    assert.NotEmpty(t, event.FirstTimestamp)
    assert.Equal(t, application, event.Source.Component)
}

func TestSecondExpiredEvent(t *testing.T) {
    client := NewMockEventClient() 
    application := "slackernews-mackerel"
    past := time.Now().Add(-24 * time.Hour)
  
    err := client.CreateLicenseEvent(application, past)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    err = client.CreateLicenseEvent(application, past)
    event, err := client.GetLicenseEvent(application, past)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)
    assert.Equal(t, int32(2), event.Count)
}

func TestValidEvent(t *testing.T) {
    client := NewMockEventClient() 
    application := "slackernews-mackerel"
    future := time.Now().Add(24 * time.Hour)
  
    err := client.CreateLicenseEvent(application, future)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    event, err := client.GetLicenseEvent(application, future)
    assert.Equal(t, "Normal", event.Type)
    assert.Equal(t, "Valid", event.Reason)
}

func TestSecondValidEvent(t *testing.T) {
    client := NewMockEventClient() 
    application := "slackernews-mackerel"
    future := time.Now().Add(24 * time.Hour)
  
    err := client.CreateLicenseEvent(application, future)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    err = client.CreateLicenseEvent(application, future)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    event, err := client.GetLicenseEvent(application, future)
    assert.NoError(t, err)
    assert.Equal(t, int32(1), event.Count)
}

func TestValidNewExpiration(t *testing.T) {
    client := NewMockEventClient() 
    application := "slackernews-mackerel"
    future := time.Now().Add(24 * time.Hour)
    renewal := time.Now().Add(48 * time.Hour)
  
    err := client.CreateLicenseEvent(application, future)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    err = client.CreateLicenseEvent(application, renewal)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 2)
}
