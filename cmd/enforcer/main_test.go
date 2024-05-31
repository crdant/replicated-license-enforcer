package main

import (
    "time"
    "testing"

    "github.com/stretchr/testify/mock"
)

type MockAPIClient struct {
    mock.Mock
}

func (m *MockAPIClient) GetExpirationDate() (time.Time, error) {
    args := m.Called()
    return args.Get(0).(time.Time), args.Error(1)
}

func TestLicenseCheckValidLicense(t *testing.T) {
    future := time.Now().Add(24 * time.Hour)
    mockClient := new(MockAPIClient)
    mockClient.On("GetExpirationDate", mock.Anything).Return(future, nil)

    err := checkLicense(mockClient)
    if err != nil {
        t.Fatalf("Expected license to be valid and got %v", err)
    }
}


func TestLicenseCheckExpriredLicense(t *testing.T) {
    future := time.Now().Add(-24 * time.Hour)
    mockClient := new(MockAPIClient)
    mockClient.On("GetExpirationDate", mock.Anything).Return(future, nil)

    err := checkLicense(mockClient)
    if err == nil {
        t.Fatalf("Expected license to be invalid and got %v", err)
    }
}

