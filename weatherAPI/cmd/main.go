package main

import (
	"log"
	"weatherAPI/internal/app_server"
)

func main() {
	if err := app_server.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
