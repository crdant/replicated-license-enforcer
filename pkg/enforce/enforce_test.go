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

    valid, err := isValid(mockClient)
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

    valid, err := isValid(mockClient)
    if err != nil {
        t.Fatalf("Expected license check to succeed and got %v", err)
    }
    if valid {
        t.Fatalf("Expected license to be invalid and got valid")
    }
}

func TestCheckExpiredLicense(t *testing.T) {
    past := time.Now().Add(-24 * time.Hour)
    name := "Slackernews"
    slug := "slackernews-mackerel"

    sdkClient := client.NewMockAPIClient(name, slug, past)
    k8sClient := events.NewMockEventClient()
    err := check(sdkClient, k8sClient)

    assert.Error(t, err)
    assert.Len(t, k8sClient.Events, 1)

    event, err := k8sClient.GetLicenseEvent(slug, past)
    require.NoError(t, err)
    require.NotNil(t, event)

    assert.Equal(t, "Warning", event.Type)
    assert.Equal(t, "Expired", event.Reason)
}
