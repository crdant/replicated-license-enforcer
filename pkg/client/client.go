package client

import (
    "io"
    "net/http"
    "time"

    license "github.com/replicatedhq/replicated-sdk/pkg/license/types"
)

// A simplified client for the Replicated SDK designed for the purposes of
// validating licenses and producing sensible error messages, log entries,
// monitoring events, etc.
type Client struct {
    HTTPClient *http.Client
    BaseURL    string
}

// Returns a new client with that will access the Replicated SDK at the
// provided URL
func NewClient(baseURL string) *Client {
    return &Client{
        HTTPClient: &http.Client{Timeout: time.Second * 30},
        BaseURL:    baseURL,
    }
}

// Common method to make requests
func (c *Client) makeRequest(method, url string, body io.Reader) (*http.Response, error) {
    req, err := http.NewRequest(method, c.BaseURL+url, body)
    if err != nil {
        return nil, err
    }
    return c.HTTPClient.Do(req)
}

// A client for a subset of the Replicated SDK as required for validating
// license and descrbing the application when discussing it in errors messages
// or recording information about the license
type ReplicatedClient interface {
  GetAppName() (string, error)
  GetAppSlug() (string, error)
  GetExpirationDate() (time.Time, error) 

  GetLicenseField(string) (*license.LicenseField, error) 
}
