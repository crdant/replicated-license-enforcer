package enforce

import (
    "time"
    "testing"

    "github.com/crdant/replicated-license-enforcer/pkg/client"

    "github.com/stretchr/testify/mock"
)

func TestIsActiveLicenseValid(t *testing.T) {
    future := time.Now().Add(24 * time.Hour)
    mockClient := client.newMockAPIClient()
    mockClient.On("GetExpirationDate", mock.Anything).Return(future, nil)

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
    mockClient := client.newMockAPIClient()
    mockClient.On("GetExpirationDate", mock.Anything).Return(past, nil)

    valid, err := isValid(mockClient)
    if err != nil {
        t.Fatalf("Expected license check to succeed and got %v", err)
    }
    if valid {
        t.Fatalf("Expected license to be invalid and got valid")
    }
}

func TestCheckValidLicense(t *testing.T) {
    past := time.Now().Add(-24 * time.Hour)
    name := "Slackernews"
    slug := "slackernews-mackerel"

    mockClient := client.newMockAPIClient()
    mockClient.On("GetExpirationDate", mock.Anything).Return(past, nil)
    mockClient.On("GetAppName", mock.Anything).Return(name, nil)
    mockClient.On("GetAppSlug", mock.Anything).Return(slug, nil)
  
}
