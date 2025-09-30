package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"banana-auction/config"
	"banana-auction/database"
	"banana-auction/rest/handlers"
	"banana-auction/rest/middlewares"
)

func Serve() {
	cfg := config.GetConfig()
	if err := database.InitDB(cfg); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /signup", handlers.SignupHandler)
	mux.HandleFunc("POST /login", handlers.LoginHandler)

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /lots", handlers.CreateLotHandler)
	protectedMux.HandleFunc("PATCH /lots/", handlers.UpdateLotHandler)
	protectedMux.HandleFunc("DELETE /lots/", handlers.DeleteLotHandler)
	protectedMux.HandleFunc("GET /lots", handlers.ListLotHandler)
	protectedMux.HandleFunc("POST /auctions", handlers.CreateAuctionHandler)
	protectedMux.HandleFunc("GET /auctions/{id}/bids", handlers.ListBidsHandler)
	protectedMux.HandleFunc("POST /auctions/{id}/bids", handlers.CreateBidHandler)

	// Apply middlewares
	handler := middlewares.CorsMiddleware(middlewares.JwtAuthMiddleware(protectedMux))
	mux.Handle("/", handler)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HttpPort),
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting %s (v%s) on port %d", cfg.ServiceName, cfg.Version, cfg.HttpPort)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
