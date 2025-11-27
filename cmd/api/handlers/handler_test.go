package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ObiaNzk/LTK-test-manu/internal/mocks"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type HandlerTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	mockService *mocks.MockServiceInterface
	handler     *Handler
}

func (s *HandlerTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockService = mocks.NewMockServiceInterface(s.ctrl)
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

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
