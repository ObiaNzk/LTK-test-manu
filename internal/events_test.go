package internal_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ObiaNzk/LTK-test-manu/internal"
	"github.com/ObiaNzk/LTK-test-manu/internal/mocks"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

//go:generate mockgen -source=events.go -destination=mocks/mock_storage.go -package=mocks

type ServiceTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockStorage *mocks.Mockstorage
	service     *internal.Service
}

func (s *ServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockStorage = mocks.NewMockstorage(s.ctrl)
	s.service = internal.NewService(s.mockStorage)
}

func (s *ServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *ServiceTestSuite) TestNewService() {
	service := internal.NewService(s.mockStorage)
	require.NotNil(s.T(), service)
}

func (s *ServiceTestSuite) TestCreateEvent_Success() {
	ctx := context.Background()
	now := time.Now()
	title := strings.Repeat("a", 101)

	request := internal.CreateEventRequest{
		Title:       title,
		Description: "pepito",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
	}

	expectedResponse := internal.CreateEventResponse{
		ID:          "test-id",
		Title:       title,
		Description: "pepito",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		CreatedAt:   now,
	}

	s.mockStorage.EXPECT().
		CreateEvent(gomock.Any(), request).
		Return(expectedResponse, nil)

	result, err := s.service.CreateEvent(ctx, request)

	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedResponse.ID, result.ID)
	require.Equal(s.T(), expectedResponse.Title, result.Title)
}

func (s *ServiceTestSuite) TestCreateEvent_EmptyTitle() {
	ctx := context.Background()
	now := time.Now()

	request := internal.CreateEventRequest{
		Title:       "",
		Description: "pepito",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
	}

	_, err := s.service.CreateEvent(ctx, request)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, internal.ErrInput)
	require.EqualError(s.T(), err, "title cannot be empty: missing input values")
}

func (s *ServiceTestSuite) TestCreateEvent_EmptyDescription() {
	ctx := context.Background()
	now := time.Now()
	title := strings.Repeat("a", 101)

	request := internal.CreateEventRequest{
		Title:       title,
		Description: "",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
	}

	_, err := s.service.CreateEvent(ctx, request)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, internal.ErrInput)
	require.EqualError(s.T(), err, "description cannot be empty: missing input values")
}

func (s *ServiceTestSuite) TestCreateEvent_MissingStartTime() {
	ctx := context.Background()
	title := strings.Repeat("a", 101)

	request := internal.CreateEventRequest{
		Title:       title,
		Description: "pepito",
		StartTime:   time.Time{},
		EndTime:     time.Now(),
	}

	_, err := s.service.CreateEvent(ctx, request)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, internal.ErrInput)
	require.EqualError(s.T(), err, "start time and end time should be set: missing input values")
}

func (s *ServiceTestSuite) TestCreateEvent_MissingEndTime() {
	ctx := context.Background()
	title := strings.Repeat("a", 101)

	request := internal.CreateEventRequest{
		Title:       title,
		Description: "pepito",
		StartTime:   time.Now(),
		EndTime:     time.Time{},
	}

	_, err := s.service.CreateEvent(ctx, request)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, internal.ErrInput)
	require.EqualError(s.T(), err, "start time and end time should be set: missing input values")
}

func (s *ServiceTestSuite) TestCreateEvent_TitleTooShort() {
	ctx := context.Background()
	now := time.Now()

	request := internal.CreateEventRequest{
		Title:       "Short",
		Description: "pepito",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
	}

	_, err := s.service.CreateEvent(ctx, request)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, internal.ErrInput)
	require.EqualError(s.T(), err, "title should have more than 100 words: missing input values")
}

func (s *ServiceTestSuite) TestCreateEvent_StorageError() {
	ctx := context.Background()
	now := time.Now()
	title := strings.Repeat("a", 101)

	request := internal.CreateEventRequest{
		Title:       title,
		Description: "pepito",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
	}

	storageError := errors.New("database error")
	s.mockStorage.EXPECT().
		CreateEvent(gomock.Any(), request).
		Return(internal.CreateEventResponse{}, storageError)

	_, err := s.service.CreateEvent(ctx, request)

	require.Error(s.T(), err)
	require.EqualError(s.T(), err, "creating event: database error")
	require.ErrorIs(s.T(), err, storageError)
}

func (s *ServiceTestSuite) TestGetEventByID_Success() {
	ctx := context.Background()
	eventID := "test-id-123"
	now := time.Now()

	expectedEvent := internal.CreateEventResponse{
		ID:          eventID,
		Title:       strings.Repeat("a", 101),
		Description: "pepito",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		CreatedAt:   now,
	}

	s.mockStorage.EXPECT().
		GetEventByID(gomock.Any(), eventID).
		Return(expectedEvent, nil)

	result, err := s.service.GetEventByID(ctx, eventID)

	require.NoError(s.T(), err)
	require.Equal(s.T(), expectedEvent.ID, result.ID)
	require.Equal(s.T(), expectedEvent.Title, result.Title)
	require.Equal(s.T(), expectedEvent.Description, result.Description)
	require.Equal(s.T(), expectedEvent.CreatedAt, result.CreatedAt)
	require.Equal(s.T(), expectedEvent.EndTime, result.EndTime)
	require.Equal(s.T(), expectedEvent.StartTime, result.StartTime)
}

func (s *ServiceTestSuite) TestGetEventByID_EmptyID() {
	ctx := context.Background()

	_, err := s.service.GetEventByID(ctx, "")

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, internal.ErrInput)
	require.EqualError(s.T(), err, "empty id: missing input values")
}

func (s *ServiceTestSuite) TestGetEventByID_NotFound() {
	ctx := context.Background()
	eventID := "nonexistent-id"

	s.mockStorage.EXPECT().
		GetEventByID(gomock.Any(), eventID).
		Return(internal.CreateEventResponse{}, internal.ErrNotFound)

	_, err := s.service.GetEventByID(ctx, eventID)

	require.Error(s.T(), err)
	require.ErrorIs(s.T(), err, internal.ErrNotFound)
	require.EqualError(s.T(), err, "getting event: not found")
}

func (s *ServiceTestSuite) TestGetEventByID_StorageError() {
	ctx := context.Background()
	eventID := "test-id"

	storageError := errors.New("database connection error")
	s.mockStorage.EXPECT().
		GetEventByID(gomock.Any(), eventID).
		Return(internal.CreateEventResponse{}, storageError)

	_, err := s.service.GetEventByID(ctx, eventID)

	require.Error(s.T(), err)
	require.EqualError(s.T(), err, "getting event: database connection error")
	require.ErrorIs(s.T(), err, storageError)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
