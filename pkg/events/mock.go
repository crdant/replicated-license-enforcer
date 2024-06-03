package events

import (
  "fmt"
  "os"
  "time"

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
    return fmt.Sprintf("%s-%s", application, date.Format(time.RFC3339))
}

func NewMockEventClient() *MockEventClient {
    setupMockEnvironment()
    return &MockEventClient{
        Events: make(map[string]*v1.Event),
    }
}

func (c *MockEventClient) GetExpiredEvent(application string, date time.Time) (*v1.Event, error) {
    key := generateEventKey(application, date)
    event, ok := c.Events[key]
    if !ok { 
      return nil, nil
    }
    return event, nil
}


func (c *MockEventClient) CreateExpiredEvent(application string, date time.Time) error {
    event, err := PrepareExpiredEvent(c, application, date)
    if err != nil {
      return err
    }
    key := generateEventKey(application, date)
    c.Events[key] = event
    return nil
}

func (c *MockEventClient) CreateValidEvent(application string, date time.Time) error {
    event, err := PrepareExpiredEvent(c, application, date)
    if err != nil {
      return err
    }
    key := generateEventKey(application, date)
    c.Events[key] = event
    return nil
}


