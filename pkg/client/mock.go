package client 

import (
    "time"

    license "github.com/replicatedhq/replicated-sdk/pkg/license/types"
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

func NewMockAPIClient(name string, slug string, expiration time.Time, fields ...*license.LicenseField) *MockAPIClient {
    mock := &MockAPIClient{
        name: name,
          
        slug: slug,
        expiration: expiration,
    }
    
    mock.On("GetAppName").Return(mock.name, nil)
    mock.On("GetAppSlug").Return(mock.slug, nil)
    mock.On("GetExpirationDate").Return(mock.expiration, nil)

    for _, field := range fields {
      mock.On("GetLicenseField", field.Name).Return(field, nil)
    }

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

func (m *MockAPIClient) GetLicenseField(field string) (*license.LicenseField, error) {
    args := m.Called()
    return args.Get(0).(*license.LicenseField), args.Error(1)
}
