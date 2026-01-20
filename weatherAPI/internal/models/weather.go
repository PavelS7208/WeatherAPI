package models

type WeatherResponse struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
	Description string  `json:"description"`
	Timestamp   int64   `json:"timestamp"`
}
