package events

import (
    "fmt"
    "strings"
    "time"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestFirstLicenseEvent(t *testing.T) {
    client := NewMockEventClient() 
    podRef := GetObjectReference()
    application := "Slackernews"
    past := time.Now().Add(-24 * time.Hour)
  
    err := client.CreateLicenseEvent(false, application, past)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    event, err := client.GetLicenseEvent(false, application, past)
    assert.NoError(t, err)
    assert.Equal(t, fmt.Sprintf("%s.", strings.ToLower(application)), event.ObjectMeta.GenerateName)
    assert.Equal(t, podRef.Namespace, event.ObjectMeta.Namespace)
    assert.Equal(t, "Warning", event.Type)
    assert.Equal(t, "Expired", event.Reason)
    assert.Equal(t, application, event.ObjectMeta.Annotations["application"])
    assert.Equal(t, past.Format(time.RFC3339), event.ObjectMeta.Annotations["expiration"])

    assert.Equal(t, fmt.Sprintf("%s license is not valid, expired %v", application, past), event.Message)

    assert.Equal(t, podRef.APIVersion, event.InvolvedObject.APIVersion)
    assert.Equal(t, podRef.Kind, event.InvolvedObject.Kind)
    assert.Equal(t, podRef.Name, event.InvolvedObject.Name)
    assert.Equal(t, podRef.Namespace, event.InvolvedObject.Namespace)
    assert.Equal(t, podRef.UID, event.InvolvedObject.UID)

    assert.NotEmpty(t, event.FirstTimestamp)
    assert.Equal(t, "replicated", event.Source.Component)
}

func TestSecondLicenseEvent(t *testing.T) {
    client := NewMockEventClient() 
    podRef := GetObjectReference()
    application := "Slackernews"
    past := time.Now().Add(-24 * time.Hour)
  
    err := client.CreateLicenseEvent(false, application, past)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    err = client.CreateLicenseEvent(false, application, past)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    event, err := client.GetLicenseEvent(false, application, past)
    assert.NoError(t, err)
    assert.Equal(t, fmt.Sprintf("%s.", strings.ToLower(application)), event.ObjectMeta.GenerateName)
    assert.Equal(t, podRef.Namespace, event.ObjectMeta.Namespace)
    assert.Equal(t, "Warning", event.Type)
    assert.Equal(t, "Expired", event.Reason)
    assert.Equal(t, int32(2), event.Count)
    assert.Equal(t, application, event.ObjectMeta.Annotations["application"])
    assert.Equal(t, past.Format(time.RFC3339), event.ObjectMeta.Annotations["expiration"])

    assert.Equal(t, fmt.Sprintf("%s license is not valid, expired %v", application, past), event.Message)

    assert.Equal(t, podRef.APIVersion, event.InvolvedObject.APIVersion)
    assert.Equal(t, podRef.Kind, event.InvolvedObject.Kind)
    assert.Equal(t, podRef.Name, event.InvolvedObject.Name)
    assert.Equal(t, podRef.Namespace, event.InvolvedObject.Namespace)
    assert.Equal(t, podRef.UID, event.InvolvedObject.UID)

    assert.NotEmpty(t, event.FirstTimestamp)
    assert.Equal(t, "replicated", event.Source.Component)
}

func TestFirstValidEvent(t *testing.T) {
    client := NewMockEventClient() 
    podRef := GetObjectReference()
    application := "Slackernews"
    past := time.Now().Add(-24 * time.Hour)
  
    err := client.CreateLicenseEvent(true, application, past)
    assert.NoError(t, err)
    assert.Len(t, client.Events, 1)

    event, err := client.GetLicenseEvent(true, application, past)
    assert.NoError(t, err)
    assert.Equal(t, fmt.Sprintf("%s.", strings.ToLower(application)), event.ObjectMeta.GenerateName)
    assert.Equal(t, podRef.Namespace, event.ObjectMeta.Namespace)
    assert.Equal(t, "Normal", event.Type)
    assert.Equal(t, "Valid", event.Reason)
    assert.Equal(t, application, event.ObjectMeta.Annotations["application"])
    assert.Equal(t, past.Format(time.RFC3339), event.ObjectMeta.Annotations["expiration"])

    assert.Equal(t, fmt.Sprintf("%s license is valid, expires %v", application, past), event.Message)

    assert.Equal(t, podRef.APIVersion, event.InvolvedObject.APIVersion)
    assert.Equal(t, podRef.Kind, event.InvolvedObject.Kind)
    assert.Equal(t, podRef.Name, event.InvolvedObject.Name)
    assert.Equal(t, podRef.Namespace, event.InvolvedObject.Namespace)
    assert.Equal(t, podRef.UID, event.InvolvedObject.UID)

    assert.NotEmpty(t, event.FirstTimestamp)
    assert.Equal(t, "replicated", event.Source.Component)
}
