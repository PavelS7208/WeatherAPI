package service

import (
	"context"
	"log"
	"weatherAPI/internal/models"

	"golang.org/x/sync/singleflight"
)

type Cache interface {
	Get(city string) (models.WeatherResponse, bool)
	Set(city string, data models.WeatherResponse)
}

type Client interface {
	FetchWeather(ctx context.Context, city string) (models.WeatherResponse, error)
}

type WeatherService struct {
	cache  Cache
	client Client
	sf     singleflight.Group
}

func NewWeatherService(cache Cache, client Client) *WeatherService {
	return &WeatherService{
		cache:  cache,
		client: client,
	}
}

func (s *WeatherService) GetWeather(ctx context.Context, city string) (models.WeatherResponse, error) {
	// Проверяем кэш
	if cached, found := s.cache.Get(city); found {
		log.Printf("Cache HIT for city: %s", city)
		return cached, nil
	}

	log.Printf("Cache MISS for city: %s", city)

	// Используем singleflight чтобы избежать дублирующих запросов к API
	v, err, _ := s.sf.Do(city, func() (interface{}, error) {
		// Двойная проверка кэша (другая горутина могла уже загрузить)
		if cached, found := s.cache.Get(city); found {
			return cached, nil
		}

		weather, err := s.client.FetchWeather(ctx, city)
		if err != nil {
			return models.WeatherResponse{}, err
		}

		s.cache.Set(city, weather)
		return weather, nil
	})

	if err != nil {
		return models.WeatherResponse{}, err
	}

	return v.(models.WeatherResponse), nil
}
