package enforce

import (
    "time"
    "testing"

    "github.com/crdant/replicated-license-enforcer/pkg/client"
    "github.com/crdant/replicated-license-enforcer/pkg/events"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)



func TestIsActiveLicenseValid(t *testing.T) {
    mockClient := client.DefaultMockAPIClient()

    enforcer := NewEnforcer(mockClient, nil)
    valid, err := enforcer.isValid(mockClient)
    if err != nil {
        t.Fatalf("Expected license check to succeed and got %v", err)
    }
    if !valid {
      t.Fatalf("Expected license to be valid and got invalid")
    }
}

func TestIsExpiredLicenseNotValid(t *testing.T) {
    past := time.Now().Add(-24 * time.Hour)
    name := "Slackernews"
    slug := "slackernews-mackerel"

    mockClient := client.NewMockAPIClient(name, slug, past)

    enforcer := NewEnforcer(mockClient, nil)
    valid, err := enforcer.isValid(mockClient)
    if err != nil {
        t.Fatalf("Expected license check to succeed and got %v", err)
    }
    if valid {
        t.Fatalf("Expected license to be invalid and got valid")
    }
}

func TestCheckValidLicense(t *testing.T) {
    future := time.Now().Add(24 * time.Hour)
    name := "Slackernews"
    slug := "slackernews-mackerel"

    sdkClient := client.NewMockAPIClient(name, slug, future)
    k8sClient := events.NewMockEventClient()
    enforcer := NewEnforcer(sdkClient, k8sClient)
    err := enforcer.Check()

    assert.NoError(t, err)
    assert.Len(t, k8sClient.Events, 1)

    event, err := k8sClient.GetLicenseEvent(slug, future)
    require.NoError(t, err)
    require.NotNil(t, event)

    assert.Equal(t, "Normal", event.Type)
    assert.Equal(t, "Valid", event.Reason)
}

func TestCheckExpiredLicense(t *testing.T) {
    past := time.Now().Add(-24 * time.Hour)
    name := "Slackernews"
    slug := "slackernews-mackerel"

    sdkClient := client.NewMockAPIClient(name, slug, past)
    k8sClient := events.NewMockEventClient()
    enforcer := NewEnforcer(sdkClient, k8sClient)
    err := enforcer.Check()

    assert.Error(t, err)
    assert.Len(t, k8sClient.Events, 1)

    event, err := k8sClient.GetLicenseEvent(slug, past)
    require.NoError(t, err)
    require.NotNil(t, event)

    assert.Equal(t, "Warning", event.Type)
    assert.Equal(t, "Expired", event.Reason)
}

func TestCheckMonitorValid(t *testing.T) {
    future := time.Now().Add(24 * time.Hour)
    name := "Slackernews"
    slug := "slackernews-mackerel"
    interval := 3 * time.Second

    sdkClient := client.NewMockAPIClient(name, slug, future)
    k8sClient := events.NewMockEventClient()
    enforcer := NewEnforcer(sdkClient, k8sClient)

    enforcer.Monitor(interval);
    time.Sleep(interval)
    enforcer.Stop()

    assert.Len(t, k8sClient.Events, 1)

    event, err := k8sClient.GetLicenseEvent(slug, future)
    require.NoError(t, err)
    require.NotNil(t, event)

    assert.Equal(t, "Normal", event.Type)
    assert.Equal(t, "Valid", event.Reason)
    assert.Equal(t, int32(1), event.Count)
}

func TestCheckMonitorExpired(t *testing.T) {
    past := time.Now().Add(-24 * time.Hour)
    name := "Slackernews"
    slug := "slackernews-mackerel"
    interval := 1 * time.Second

    sdkClient := client.NewMockAPIClient(name, slug, past)
    k8sClient := events.NewMockEventClient()
    enforcer := NewEnforcer(sdkClient, k8sClient)

    enforcer.Monitor(interval);
    time.Sleep(3 * interval)
    enforcer.Stop()

    assert.Len(t, k8sClient.Events, 1)

    event, err := k8sClient.GetLicenseEvent(slug, past)
    require.NoError(t, err)
    require.NotNil(t, event)

    assert.Equal(t, "Warning", event.Type)
    assert.Equal(t, "Expired", event.Reason)
    // there should be at least one check per second in our interval, but
    // it is non-deterministic with exponential backoff
    assert.GreaterOrEqual(t, event.Count, int32(interval.Seconds()))
}
