package cache

import (
	"context"
	"log"
	"sync"
	"time"
	"weatherAPI/internal/models"
)

type WeatherCacheItem struct {
	Data      models.WeatherResponse
	ExpiresAt time.Time
}

type WeatherCache struct {
	mu    sync.RWMutex
	items map[string]WeatherCacheItem
	ttl   time.Duration
}

func NewWeatherCache(ttl time.Duration) *WeatherCache {
	return &WeatherCache{
		items: make(map[string]WeatherCacheItem),
		ttl:   ttl,
	}
}

// Запускаем фоновую процедуру - которая зависит от контекста
func (c *WeatherCache) StartCleanup(ctx context.Context) {
	go c.cleanupExpired(ctx)
}

func (c *WeatherCache) Get(city string) (models.WeatherResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[city]
	if !exists || time.Now().After(item.ExpiresAt) {
		return models.WeatherResponse{}, false
	}
	return item.Data, true
}

func (c *WeatherCache) Set(city string, data models.WeatherResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[city] = WeatherCacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(c.ttl),
	}
}

func (c *WeatherCache) cleanupExpired(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.performCleanup()
		case <-ctx.Done():
			log.Println("Cache cleanup stopped")
			return
		}
	}
}

func (c *WeatherCache) performCleanup() {
	now := time.Now()
	var expiredKeys []string

	c.mu.RLock()
	for city, item := range c.items {
		if now.After(item.ExpiresAt) {
			expiredKeys = append(expiredKeys, city)
		}
	}
	c.mu.RUnlock()

	if len(expiredKeys) > 0 {
		c.mu.Lock()
		for _, city := range expiredKeys {
			delete(c.items, city)
		}
		c.mu.Unlock()
		log.Printf("Cache cleanup: removed %d expired items", len(expiredKeys))
	}
}
