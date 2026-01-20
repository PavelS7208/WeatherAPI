package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"weatherAPI/internal/models"
)

type mockWeatherService struct {
	returnError error
}

func (m *mockWeatherService) GetWeather(ctx context.Context, city string) (models.WeatherResponse, error) {
	if m.returnError != nil {
		return models.WeatherResponse{}, m.returnError
	}
	return models.WeatherResponse{
		City:        city,
		Temperature: 20.0,
		Description: "Cloudy",
		Timestamp:   time.Now().Unix(),
	}, nil
}

func TestWeatherHandler_Success(t *testing.T) {
	service := &mockWeatherService{}
	handler := NewWeatherHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/weather?city=Tokyo", nil)
	w := httptest.NewRecorder()

	handler.HandleGetWeather(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.WeatherResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.City != "Tokyo" {
		t.Errorf("Expected city Tokyo, got %s", response.City)
	}
}

func TestWeatherHandler_MissingCity(t *testing.T) {
	service := &mockWeatherService{}
	handler := NewWeatherHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/weather", nil)
	w := httptest.NewRecorder()

	handler.HandleGetWeather(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestWeatherHandler_MethodNotAllowed(t *testing.T) {
	service := &mockWeatherService{}
	handler := NewWeatherHandler(service)

	req := httptest.NewRequest(http.MethodPost, "/weather?city=Rome", nil)
	w := httptest.NewRecorder()

	handler.HandleGetWeather(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestWeatherHandler_ServiceError(t *testing.T) {
	service := &mockWeatherService{
		returnError: errors.New("service error"),
	}
	handler := NewWeatherHandler(service)

	req := httptest.NewRequest(http.MethodGet, "/weather?city=Madrid", nil)
	w := httptest.NewRecorder()

	handler.HandleGetWeather(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 500, got %d", w.Code)
	}
}
