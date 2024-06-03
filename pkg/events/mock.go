package events

import (
  "fmt"
  "os"
  "time"

  "github.com/charmbracelet/log"
  v1 "k8s.io/api/core/v1"
)

type MockEventClient struct {
    Events   map[string]*v1.Event
}

// Mock environment setup for testing
func setupMockEnvironment() {
    os.Setenv("POD_NAME", "slackernews-84d6df6674-c4vg7")
    os.Setenv("POD_NAMESPACE", "slackernews")
    os.Setenv("POD_UID", "0e8d56c7-6277-4a79-9847-bdcb3b4e3184")
}

func generateEventKey(application string, date time.Time) string {
    fieldSelector := getFieldSelector(date)
    labelSelector := getLabelSelector(application, date)
    log.Debug("Generated event key", "key", fmt.Sprintf("%s,%s", fieldSelector, labelSelector))
    return fmt.Sprintf("%s,%s", fieldSelector, labelSelector)
}

func NewMockEventClient() *MockEventClient {
    setupMockEnvironment()
    logLevel := os.Getenv("LOG_LEVEL")
    if logLevel == "" {
      logLevel = "fatal"
    }
    log.ParseLevel(logLevel)
    return &MockEventClient{
        Events: make(map[string]*v1.Event),
    }
}

func (c *MockEventClient) GetLicenseEvent(application string, date time.Time) (*v1.Event, error) {
    key := generateEventKey(application, date)
    event, ok := c.Events[key]
    if !ok { 
      log.Debug("Event not found", "key", key)
      return nil, nil
    }
    log.Debug("returning event", "event", event)
    return event, nil
}


func (c *MockEventClient) CreateLicenseEvent(application string, date time.Time) error {
    event, err := PrepareLicenseEvent(c, application, date)
    if err != nil {
      log.Error("Error preparing event", "error", err)
      return err
    }
    key := generateEventKey(application, date)
    log.Debug("adding event to store", "key", key, "event", event)
    c.Events[key] = event
    return nil
}
