package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mohammadrabetian/ports/domain"
	"github.com/mohammadrabetian/ports/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPortService struct {
	mock.Mock
}

func (m *MockPortService) GetPortByID(id string) (*domain.Port, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Port), args.Error(1)
}

func (m *MockPortService) AddOrUpdatePort(port *domain.Port) error {
	args := m.Called(port)
	return args.Error(0)
}

func (m *MockPortService) ListPorts(page int, limit int) ([]*domain.Port, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Port), args.Error(1)
}

var testPort = &domain.Port{
	ID:      "AEAJM",
	Name:    "Jebel Ali",
	City:    "Jebel Ali",
	Country: "United Arab Emirates",
}

var mockPortService = new(MockPortService)
var router *gin.Engine

func init() {
	handlers.InitPortHandlers(mockPortService)
	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.GET("/api/v1/ports/:id", handlers.GetPortByID)
	router.GET("/api/v1/ports", handlers.ListPorts)
}

func TestGetPortByID_Success(t *testing.T) {
	mockPortService.On("GetPortByID", testPort.ID).Return(testPort, nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ports/"+testPort.ID, nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]*domain.Port
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Compare only relevant fields
	assert.Equal(t, testPort.ID, response["port"].ID)
	assert.Equal(t, testPort.Name, response["port"].Name)
	assert.Equal(t, testPort.City, response["port"].City)
	assert.Equal(t, testPort.Country, response["port"].Country)

	mockPortService.AssertExpectations(t)
}

func TestGetPortByID_NotFound(t *testing.T) {
	mockPortService.On("GetPortByID", "unknown").Return((*domain.Port)(nil), errors.New("Port not found")).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ports/unknown", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockPortService.AssertExpectations(t)
}

func TestListPorts_Success(t *testing.T) {
	mockPortService.On("ListPorts", handlers.DefaultPage, handlers.DefaultLimit).Return([]*domain.Port{testPort}, nil).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ports", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string][]*domain.Port
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Compare only relevant fields
	assert.Equal(t, 1, len(response["ports"]))
	assert.Equal(t, testPort.ID, response["ports"][0].ID)
	assert.Equal(t, testPort.Name, response["ports"][0].Name)
	assert.Equal(t, testPort.City, response["ports"][0].City)
	assert.Equal(t, testPort.Country, response["ports"][0].Country)

	mockPortService.AssertExpectations(t)
}

func TestListPorts_Error(t *testing.T) {
	mockPortService.On("ListPorts", handlers.DefaultPage, handlers.DefaultLimit).Return(([]*domain.Port)(nil), errors.New("Failed to retrieve the ports")).Once()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/ports", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockPortService.AssertExpectations(t)
}
