package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"weatherAPI/internal/models"
)

type WeatherService interface {
	GetWeather(ctx context.Context, city string) (models.WeatherResponse, error)
}

type WeatherHandler struct {
	service WeatherService
}

func NewWeatherHandler(service WeatherService) *WeatherHandler {
	return &WeatherHandler{service: service}
}

func (h *WeatherHandler) HandleGetWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	weather, err := h.service.GetWeather(r.Context(), city)
	if err != nil {
		log.Printf("Error fetching weather: %v", err)
		http.Error(w, "Failed to fetch weather data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(weather); err != nil {
		log.Printf("Error encoding response: %v", err)
	}

}
