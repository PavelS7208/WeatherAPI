package client

import (
	"context"
	"net/http"
	"time"
	"weatherAPI/internal/models"
)

type WeatherClient struct {
	apiURL     string
	apiKey     string
	httpClient *http.Client
}

func NewWeatherClient(apiURL, apiKey string) *WeatherClient {
	return &WeatherClient{
		apiURL: apiURL,
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *WeatherClient) FetchWeather(ctx context.Context, city string) (models.WeatherResponse, error) {
	// Здесь реальный запрос к внешнему API
	// Например: GET {apiURL}/weather?q={city}&appid={apiKey}

	// Симуляция для примера
	time.Sleep(500 * time.Millisecond)

	return models.WeatherResponse{
		City:        city,
		Temperature: 15.5 + float64(time.Now().Unix()%10),
		Description: "Partly cloudy",
		Timestamp:   time.Now().Unix(),
	}, nil
}
