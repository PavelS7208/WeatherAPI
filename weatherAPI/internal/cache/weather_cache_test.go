package cache

import (
	"testing"
	"time"
	"weatherAPI/internal/models"
)

func TestWeatherCache_SetAndGet(t *testing.T) {
	cache := NewInMemoryWeatherCache(5 * time.Second)

	weather := models.WeatherResponse{
		City:        "Moscow",
		Temperature: 20.5,
		Description: "Sunny",
		Timestamp:   time.Now().Unix(),
	}

	// Сохраняем в кэш
	cache.Set("Moscow", weather)

	// Получаем из кэша
	cached, found := cache.Get("Moscow")
	if !found {
		t.Fatal("Expected to find cached weather data")
	}

	if cached.City != weather.City {
		t.Errorf("Expected city %s, got %s", weather.City, cached.City)
	}
	if cached.Temperature != weather.Temperature {
		t.Errorf("Expected temperature %f, got %f", weather.Temperature, cached.Temperature)
	}
}

func TestWeatherCache_Expiration(t *testing.T) {
	cache := NewInMemoryWeatherCache(100 * time.Millisecond)

	weather := models.WeatherResponse{
		City:        "London",
		Temperature: 15.0,
		Description: "Rainy",
		Timestamp:   time.Now().Unix(),
	}

	cache.Set("London", weather)

	// Проверяем что данные есть
	if _, found := cache.Get("London"); !found {
		t.Fatal("Expected to find cached data")
	}

	// Ждем истечения TTL
	time.Sleep(150 * time.Millisecond)

	// Проверяем что данные истекли
	if _, found := cache.Get("London"); found {
		t.Fatal("Expected cache to expire")
	}
}

func TestWeatherCache_NotFound(t *testing.T) {
	cache := NewInMemoryWeatherCache(5 * time.Second)

	if _, found := cache.Get("NonExistent"); found {
		t.Fatal("Expected not to find non-existent city")
	}
}
