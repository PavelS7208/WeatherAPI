package service

import (
	"context"
	"testing"
	"time"
	"weatherAPI/internal/models"
)

type mockCache struct {
	data      map[string]models.WeatherResponse
	getCalled int
	setCalled int
}

func (m *mockCache) Get(city string) (models.WeatherResponse, bool) {
	m.getCalled++
	data, exists := m.data[city]
	return data, exists
}

func (m *mockCache) Set(city string, data models.WeatherResponse) {
	m.setCalled++
	if m.data == nil {
		m.data = make(map[string]models.WeatherResponse)
	}
	m.data[city] = data
}

type mockClient struct {
	fetchCalled int
	returnError error
}

func (m *mockClient) FetchWeather(_ context.Context, city string) (models.WeatherResponse, error) {
	m.fetchCalled++
	if m.returnError != nil {
		return models.WeatherResponse{}, m.returnError
	}
	return models.WeatherResponse{
		City:        city,
		Temperature: 18.5,
		Description: "Clear",
		Timestamp:   time.Now().Unix(),
	}, nil
}

func TestWeatherService_CacheHit(t *testing.T) {
	cache := &mockCache{
		data: map[string]models.WeatherResponse{
			"Paris": {
				City:        "Paris",
				Temperature: 22.0,
				Description: "Sunny",
			},
		},
	}
	client := &mockClient{}

	service := NewWeatherService(cache, client)

	weather, err := service.GetWeather(context.Background(), "Paris")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if weather.City != "Paris" {
		t.Errorf("Expected city Paris, got %s", weather.City)
	}

	if client.fetchCalled != 0 {
		t.Error("Client should not be called on cache hit")
	}

	if cache.getCalled != 1 {
		t.Errorf("Expected cache.Get to be called once, called %d times", cache.getCalled)
	}
}

func TestWeatherService_CacheMiss(t *testing.T) {
	cache := &mockCache{}
	client := &mockClient{}

	service := NewWeatherService(cache, client)

	weather, err := service.GetWeather(context.Background(), "Berlin")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if weather.City != "Berlin" {
		t.Errorf("Expected city Berlin, got %s", weather.City)
	}

	if client.fetchCalled != 1 {
		t.Error("Client should be called on cache miss")
	}

	if cache.setCalled != 1 {
		t.Error("Cache.Set should be called after fetching from client")
	}
}
