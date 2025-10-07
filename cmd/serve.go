package cmd

import (
	"banana-auction/config"
	"banana-auction/api"
	"banana-auction/internal/infrastructure/persistence/postgres"
	"fmt"
	"log"
	"net/http"
)

func Serve() {
	cfg := config.GetConfig()
	if err := postgres.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// handler := routes.SetupRoutes()
	handler := api.SetupRoutes()

	log.Printf("Starting server on :%d", cfg.HttpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.HttpPort), handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
