package app_server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weatherAPI/config"
	"weatherAPI/internal/cache"
	"weatherAPI/internal/client"
	"weatherAPI/internal/handler"
	"weatherAPI/internal/service"
)

const (
	shutdownTimeout = 10 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 120 * time.Second
)

type CleanableCache interface {
	service.Cache
	StartCleanup(ctx context.Context)
}

type App struct {
	server *http.Server
	cache  CleanableCache
}

func NewApp(cfg *config.Config) (*App, error) {
	weatherCache := cache.NewInMemoryWeatherCache(cfg.CacheTTL)
	weatherClient := client.NewWeatherClient(cfg.ExternalAPIURL, cfg.APIKey)
	weatherService := service.NewWeatherService(weatherCache, weatherClient)
	weatherHandler := handler.NewWeatherHandler(weatherService)

	server := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      setupRoutes(weatherHandler),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return &App{
		server: server,
		cache:  weatherCache,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	a.cache.StartCleanup(ctx)

	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting server on %s", a.server.Addr)
		serverErrors <- a.server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		log.Println("Shutdown signal received")
	}

	return a.shutdown()
}

func (a *App) shutdown() error {
	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Forced shutdown: %v", err)
		return err
	}

	log.Println("Server stopped gracefully")
	return nil
}

func setupRoutes(weatherHandler *handler.WeatherHandler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/weather", weatherHandler.HandleGetWeather)
	mux.HandleFunc("/health", healthCheck)
	return mux
}

func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		log.Printf("Error writing health check response: %v", err)
	}
}

// Run — точка входа из main
func Run() error {
	app, err := NewApp(config.Load())
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return app.Run(ctx)
}
