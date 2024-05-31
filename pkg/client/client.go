package client

import (
    "io"
    "net/http"
    "time"
)

type Client struct {
    HTTPClient *http.Client
    BaseURL    string
}

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

type ReplicatedClient interface {
  GetAppName() (string, error)
  GetAppSlug() (string, error)
  GetExpirationDate() (time.Time, error) 
}
