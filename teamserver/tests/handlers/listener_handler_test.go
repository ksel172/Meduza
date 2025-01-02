package handler_tests

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCreateListener(t *testing.T) {
	// mockListenerDAL := &mocks.MockListenerDAL{}
	// handler := handlers.NewListenersHandler(mockListenerDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestGetListenerById(t *testing.T) {
	// mockListenerDAL := &mocks.MockListenerDAL{}
	// handler := handlers.NewListenersHandler(mockListenerDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestGetAllListeners(t *testing.T) {
	// mockListenerDAL := &mocks.MockListenerDAL{}
	// handler := handlers.NewListenersHandler(mockListenerDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestDeleteListener(t *testing.T) {
	// mockListenerDAL := &mocks.MockListenerDAL{}
	// handler := handlers.NewListenersHandler(mockListenerDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestUpdateListener(t *testing.T) {
	// mockListenerDAL := &mocks.MockListenerDAL{}
	// handler := handlers.NewListenersHandler(mockListenerDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
	}{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
