package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ObiaNzk/LTK-test-manu/cmd/api/handlers/mocks"
	"github.com/ObiaNzk/LTK-test-manu/internal"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type HandlerTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockService *mocks.MockeventsService
	handler     *Handler
}

func (s *HandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockService = mocks.NewMockeventsService(s.ctrl)
	s.handler = NewHandler(s.mockService)
}

func (s *HandlerTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func (s *HandlerTestSuite) TestNewHandler() {
	handler := NewHandler(s.mockService)
	require.NotNil(s.T(), handler)
	require.NotNil(s.T(), handler.eventsService)
}

func (s *HandlerTestSuite) TestHelloWorld_Success() {
	expectedMessage := "Hello, World!"
	s.mockService.EXPECT().HelloWorld().Return(expectedMessage).Times(1)

	reqBody := `{"message":"test message"}`
	req := httptest.NewRequest(http.MethodPost, "/hello", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.HelloWorld(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	require.Equal(s.T(), "text/plain", resp.Header.Get("Content-Type"))
	require.Equal(s.T(), expectedMessage, w.Body.String())
}

func (s *HandlerTestSuite) TestHelloWorld_CustomMessage() {
	customMessage := "Custom Hello!"
	s.mockService.EXPECT().HelloWorld().Return(customMessage).Times(1)

	reqBody := `{"message":"custom message"}`
	req := httptest.NewRequest(http.MethodPost, "/hello", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.HelloWorld(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	require.Equal(s.T(), customMessage, w.Body.String())
}

func (s *HandlerTestSuite) TestHelloWorld_EmptyMessage() {
	s.mockService.EXPECT().HelloWorld().Return("").Times(1)

	reqBody := `{"message":"empty test"}`
	req := httptest.NewRequest(http.MethodPost, "/hello", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.HelloWorld(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	require.Empty(s.T(), w.Body.String())
}

func (s *HandlerTestSuite) TestHelloWorld_InvalidJSON() {
	reqBody := `{"message": invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/hello", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.HelloWorld(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "Invalid JSON format")
}

func (s *HandlerTestSuite) TestHelloWorld_EmptyBody() {
	req := httptest.NewRequest(http.MethodPost, "/hello", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.HelloWorld(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "Invalid JSON format")
}

func (s *HandlerTestSuite) TestGetEventByID_Success() {
	eventID := "123e4567-e89b-12d3-a456-426614174000"
	expectedEvent := internal.CreateEventResponse{
		ID:          eventID,
		Title:       "Test Event",
		Description: "Test Description",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		CreatedAt:   time.Now(),
	}

	s.mockService.EXPECT().
		GetEventByID(gomock.Any(), eventID).
		Return(expectedEvent, nil).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/events/"+eventID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", eventID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	s.handler.GetEventByID(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusOK, resp.StatusCode)
	require.Equal(s.T(), "application/json", resp.Header.Get("Content-Type"))
	require.Contains(s.T(), w.Body.String(), eventID)
	require.Contains(s.T(), w.Body.String(), "Test Event")
}

func (s *HandlerTestSuite) TestGetEventByID_NotFound() {
	eventID := "nonexistent-id"

	s.mockService.EXPECT().
		GetEventByID(gomock.Any(), eventID).
		Return(internal.CreateEventResponse{}, internal.ErrNotFound).
		Times(1)

	req := httptest.NewRequest(http.MethodGet, "/events/"+eventID, nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", eventID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	s.handler.GetEventByID(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusNotFound, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "event not found")
}

func (s *HandlerTestSuite) TestGetEventByID_EmptyID() {
	req := httptest.NewRequest(http.MethodGet, "/events/", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	w := httptest.NewRecorder()

	s.handler.GetEventByID(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "empty event")
}

func (s *HandlerTestSuite) TestCreateEvent_Success() {
	now := time.Now()
	longTitle := strings.Repeat("a", 101)

	requestBody := map[string]interface{}{
		"title":       longTitle,
		"description": "Test Description",
		"start_time":  now.Format(time.RFC3339),
		"end_time":    now.Add(time.Hour).Format(time.RFC3339),
	}

	expectedEvent := internal.CreateEventResponse{
		ID:          "123e4567-e89b-12d3-a456-426614174000",
		Title:       longTitle,
		Description: "Test Description",
		StartTime:   now,
		EndTime:     now.Add(time.Hour),
		CreatedAt:   now,
	}

	s.mockService.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(expectedEvent, nil).
		Times(1)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.CreateEvent(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusCreated, resp.StatusCode)
	require.Equal(s.T(), "application/json", resp.Header.Get("Content-Type"))
	require.Contains(s.T(), w.Body.String(), expectedEvent.ID)
	require.Contains(s.T(), w.Body.String(), longTitle)
}

func (s *HandlerTestSuite) TestCreateEvent_InvalidJSON() {
	reqBody := `{"title": invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.CreateEvent(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "Invalid JSON format")
}

func (s *HandlerTestSuite) TestCreateEvent_EmptyTitle() {
	now := time.Now()
	requestBody := map[string]interface{}{
		"title":       "",
		"description": "Test Description",
		"start_time":  now.Format(time.RFC3339),
		"end_time":    now.Add(time.Hour).Format(time.RFC3339),
	}

	s.mockService.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(internal.CreateEventResponse{}, nil).
		Times(1)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.CreateEvent(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "empty title")
}

func (s *HandlerTestSuite) TestCreateEvent_TitleTooShort() {
	now := time.Now()
	requestBody := map[string]interface{}{
		"title":       "Short",
		"description": "Test Description",
		"start_time":  now.Format(time.RFC3339),
		"end_time":    now.Add(time.Hour).Format(time.RFC3339),
	}

	s.mockService.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(internal.CreateEventResponse{}, nil).
		Times(1)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.CreateEvent(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "title should have more than 100 words")
}

func (s *HandlerTestSuite) TestCreateEvent_MissingStartTime() {
	now := time.Now()
	longTitle := strings.Repeat("a", 101)
	requestBody := map[string]interface{}{
		"title":       longTitle,
		"description": "Test Description",
		"end_time":    now.Add(time.Hour).Format(time.RFC3339),
	}

	s.mockService.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(internal.CreateEventResponse{}, nil).
		Times(1)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.CreateEvent(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "start time and end time should be set")
}

func (s *HandlerTestSuite) TestCreateEvent_MissingEndTime() {
	now := time.Now()
	longTitle := strings.Repeat("a", 101)
	requestBody := map[string]interface{}{
		"title":       longTitle,
		"description": "Test Description",
		"start_time":  now.Format(time.RFC3339),
	}

	s.mockService.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(internal.CreateEventResponse{}, nil).
		Times(1)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.CreateEvent(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "start time and end time should be set")
}

func (s *HandlerTestSuite) TestCreateEvent_StartTimeAfterEndTime() {
	now := time.Now()
	longTitle := strings.Repeat("a", 101)
	requestBody := map[string]interface{}{
		"title":       longTitle,
		"description": "Test Description",
		"start_time":  now.Add(time.Hour).Format(time.RFC3339),
		"end_time":    now.Format(time.RFC3339),
	}

	s.mockService.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(internal.CreateEventResponse{}, nil).
		Times(1)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.CreateEvent(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusBadRequest, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "start time should be before end time")
}

func (s *HandlerTestSuite) TestCreateEvent_ServiceErrorConflict() {
	now := time.Now()
	longTitle := strings.Repeat("a", 101)
	requestBody := map[string]interface{}{
		"title":       longTitle,
		"description": "Test Description",
		"start_time":  now.Format(time.RFC3339),
		"end_time":    now.Add(time.Hour).Format(time.RFC3339),
	}

	s.mockService.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(internal.CreateEventResponse{}, internal.ErrPepito).
		Times(1)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.CreateEvent(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusConflict, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "pepito")
}

func (s *HandlerTestSuite) TestCreateEvent_ServiceErrorGeneric() {
	now := time.Now()
	longTitle := strings.Repeat("a", 101)
	requestBody := map[string]interface{}{
		"title":       longTitle,
		"description": "Test Description",
		"start_time":  now.Format(time.RFC3339),
		"end_time":    now.Add(time.Hour).Format(time.RFC3339),
	}

	s.mockService.EXPECT().
		CreateEvent(gomock.Any(), gomock.Any()).
		Return(internal.CreateEventResponse{}, internal.ErrInput).
		Times(1)

	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(jsonBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handler.CreateEvent(w, req)

	resp := w.Result()
	require.Equal(s.T(), http.StatusInternalServerError, resp.StatusCode)
	require.Contains(s.T(), w.Body.String(), "Error creating event")
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
