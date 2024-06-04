package client

import (
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestNewClient(t *testing.T) {
    baseURL := "http://replicated:3000"
    client := NewClient(baseURL)
    if client.BaseURL != baseURL {
        t.Errorf("Expected baseURL to be %s, got %s", baseURL, client.BaseURL)
    }
    if client.HTTPClient == nil {
        t.Errorf("Expected HTTPClient to be initialized, got nil")
    }
}

func TestValidGetRequest(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            t.Errorf("Expected method GET, got %s", r.Method)
        }
        if r.URL.String() != "/test" {
            t.Errorf("Expected URL to be /test, got %s", r.URL.String())
        }
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("response"))
    }))
    defer server.Close()

    client := NewClient(server.URL)
    resp, err := client.makeRequest(http.MethodGet, "/test", nil)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Expected status code 200, got %d", resp.StatusCode)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    if string(body) != "response" {
        t.Errorf("Expected response body to be 'response', got '%s'", string(body))
    }
}

func TestNotFoundGetRequest(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            t.Errorf("Expected method GET, got %s", r.Method)
        }
        if r.URL.String() == "/test" {
            t.Errorf("Expected non-existing URL, got %s", r.URL.String())
        }
        w.WriteHeader(http.StatusNotFound)
    }))
    defer server.Close()

    client := NewClient(server.URL)
    resp, err := client.makeRequest(http.MethodGet, "/invalid", nil)
    if err != nil {
        t.Errorf("Expected no error, got %v", err)
    }
    if resp.StatusCode != http.StatusNotFound {
        t.Errorf("Expected status code 401, got %d", resp.StatusCode)
    }

    body, _ := ioutil.ReadAll(resp.Body)
    if string(body) != "" {
        t.Errorf("Expected response body to be empty, got '%s'", string(body))
    }
}
