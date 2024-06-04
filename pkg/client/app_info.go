package client

import (
    "encoding/json"

    app "github.com/replicatedhq/replicated-sdk/pkg/handlers"
)

type AppInfo = app.GetCurrentAppInfoResponse

// Returns the name of for the application, as returned from the Replicated SDK
func (c *Client) GetAppName() (string, error) {
    info, err := c.GetAppInfo()
    if err != nil {
      return "", err
    }
    return info.AppName, nil 
}

// Returns the application slug for the application, as returned from the
// Replicated SDK
func (c *Client) GetAppSlug() (string, error) {
    info, err := c.GetAppInfo()
    if err != nil {
      return "", err
    }
    return info.AppSlug, nil 
}

// List details about an application instance, including the app name, location
// of the Helm chart in the Replicated OCI registry, and details about the
// current application release that the instance is running.
func (c *Client) GetAppInfo() (*AppInfo, error) {
    response, err := c.makeRequest("GET", "/api/v1/app/info", nil)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    var appInfo AppInfo
    if err := json.NewDecoder(response.Body).Decode(&appInfo); err != nil {
         return nil, err
    }
    return &appInfo, nil
}
