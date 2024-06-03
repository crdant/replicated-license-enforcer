package client

import (
    "encoding/json"

    app "github.com/replicatedhq/replicated-sdk/pkg/handlers"
)

type AppInfo = app.GetCurrentAppInfoResponse

type AppClient interface {
  GetAppName() (string, error) 
  GetAppSlug() (string, error) 
}

func (c *Client) GetAppName() (string, error) {
    info, err := c.GetAppInfo()
    if err != nil {
      return "", err
    }
    return info.AppName, nil 
}

func (c *Client) GetAppSlug() (string, error) {
    info, err := c.GetAppInfo()
    if err != nil {
      return "", err
    }
    return info.AppSlug, nil 
}

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
