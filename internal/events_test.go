package internal

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/suite"
)

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks github.com/ObiaNzk/LTK-test-manu/internal ServiceInterface

type ServiceTestSuite struct {
	suite.Suite
	service *Service
}

func (s *ServiceTestSuite) SetupTest() {
	s.service = NewService()
}

func (s *ServiceTestSuite) TestNewService() {
	service := NewService()
	require.NotNil(s.T(), service)
}

func (s *ServiceTestSuite) TestHelloWorld() {
	result := s.service.HelloWorld()
	require.Equal(s.T(), "Hello, World!", result)
	require.NotEmpty(s.T(), result)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
