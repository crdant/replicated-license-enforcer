package client 

import (
    "time"

    "github.com/stretchr/testify/mock"
)

type MockAPIClient struct {
    name string;
    slug string;
    expiration time.Time;

    mock.Mock;
}


func DefaultMockAPIClient() *MockAPIClient {
    return NewMockAPIClient("Slackernews", "slackernews-mackerel", time.Now().Add(24 * time.Hour))
}

func NewMockAPIClient(name string, slug string, expiration time.Time) *MockAPIClient {
    mock := &MockAPIClient{
        name: name,
          
        slug: slug,
        expiration: expiration,
    }
    
    mock.On("GetAppName").Return(mock.name, nil)
    mock.On("GetAppSlug").Return(mock.slug, nil)
    mock.On("GetExpirationDate").Return(mock.expiration, nil)

    return mock
}

func (m *MockAPIClient) GetExpirationDate() (time.Time, error) {
    args := m.Called()
    return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockAPIClient) GetAppName() (string, error) {
    args := m.Called()
    return args.Get(0).(string), args.Error(1)
}

func (m *MockAPIClient) GetAppSlug() (string, error) {
    args := m.Called()
    return args.Get(0).(string), args.Error(1)
}


