package internal

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type mockStorage struct{}

func (m *mockStorage) CreateEvent(ctx context.Context, event CreateEventRequest) (CreateEventResponse, error) {
	return CreateEventResponse{}, nil
}

func (m *mockStorage) GetEvents(ctx context.Context) ([]CreateEventResponse, error) {
	return []CreateEventResponse{}, nil
}

func (m *mockStorage) GetEventByID(ctx context.Context, id string) (CreateEventResponse, error) {
	return CreateEventResponse{}, nil
}

type ServiceTestSuite struct {
	suite.Suite
	mockStorage *mockStorage
	service     *Service
}

func (s *ServiceTestSuite) SetupTest() {
	s.mockStorage = &mockStorage{}
	s.service = NewService(s.mockStorage)
}

func (s *ServiceTestSuite) TestNewService() {
	service := NewService(s.mockStorage)
	require.NotNil(s.T(), service)
	require.NotNil(s.T(), service.storage)
}

func (s *ServiceTestSuite) TestHelloWorld() {
	result := s.service.HelloWorld()
	require.Equal(s.T(), "Hello, World!", result)
	require.NotEmpty(s.T(), result)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
