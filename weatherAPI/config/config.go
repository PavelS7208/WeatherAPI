package config

import (
	"os"
	"time"
)

type Config struct {
	ServerAddr     string
	CacheTTL       time.Duration
	ExternalAPIURL string
	APIKey         string
}

func Load() *Config {
	return &Config{
		ServerAddr:     getEnv("SERVER_ADDR", ":8080"),
		CacheTTL:       parseDuration(getEnv("CACHE_TTL", "3m")),
		ExternalAPIURL: getEnv("EXTERNAL_API_URL", "https://api.openweathermap.org/data/2.5"),
		APIKey:         getEnv("API_KEY", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

// Если ошиблись в дефолтном параметре, то будет 10 минут
func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 3 * time.Minute
	}
	return d
}
