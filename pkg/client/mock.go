package client 

import (
    "time"

    "github.com/stretchr/testify/mock"
)

type MockAPIClient struct {
    mock.Mock
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


