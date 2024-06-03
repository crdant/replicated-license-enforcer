package client

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

var mockAppInfo = `{
  "appSlug": "slackernews-mackerel",
  "appName": "SlackerNews",
  "appStatus": "unavailable",
  "helmChartURL": "oci://registry.shortrib.io/slackernews-mackerel/stable/slackernews",
  "currentRelease": {
    "versionLabel": "1.1.0-rc.2",
    "releaseNotes": "",
    "createdAt": "2024-05-13T15:23:45Z",
    "deployedAt": "2024-05-30T14:51:14-04:00",
    "helmReleaseName": "slackernews",
    "helmReleaseRevision": 2,
    "helmReleaseNamespace": "slackernews-demo"
  }
}`

func TestGetAppName(t *testing.T) {
    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockAppInfo))
    }))
    defer server.Close()

    // Client pointing to the mock server
    c := NewClient(server.URL)

    // Execute the function
    name, err := c.GetAppName()
    if err != nil {
        t.Fatalf("Did not expect an error, got %v", err)
    }

    if name != "SlackerNews" {
        t.Errorf("Expected app name 'slackernews', got '%s'", name)
    }
}

func TestGetAppSlug(t *testing.T) {
    // Mock server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(mockAppInfo))
    }))
    defer server.Close()

    // Client pointing to the mock server
    c := NewClient(server.URL)

    // Execute the function
    slug, err := c.GetAppSlug()
    if err != nil {
        t.Fatalf("Did not expect an error, got %v", err)
    }

    if slug != "slackernews-mackerel" {
        t.Errorf("Expected app slug 'slackernews-mackerel', got '%s'", slug)
    }
}

