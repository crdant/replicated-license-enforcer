package events

import (
  "os"
  "time"

  v1 "k8s.io/api/core/v1"
)

type MockEventClient struct {
      EventsCreated []*v1.Event
}

// Mock environment setup for testing
func setupMockEnvironment() {
    os.Setenv("POD_NAME", "slackernews-84d6df6674-c4vg7")
    os.Setenv("POD_NAMESPACE", "slackernews")
    os.Setenv("POD_UID", "0e8d56c7-6277-4a79-9847-bdcb3b4e3184")
}

func (c *MockEventClient) CreateExpiredEvent(application string, date time.Time) error {
    event := NewExpiredEvent(application, date)
    c.EventsCreated = append(c.EventsCreated, event)
    return nil
}

func (c *MockEventClient) CreateValidEvent(application string, date time.Time) error {
    event := NewExpiredEvent(application, date)
    c.EventsCreated = append(c.EventsCreated, event)
    return nil
}


