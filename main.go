package main

import (
	"log"
	"net/http"
	"time"

	"banana-auction/db"
	"banana-auction/handlers"
)

func main() {
	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	mux := http.NewServeMux()

	// Routes without JWT middleware
	mux.HandleFunc("POST /signup", handlers.Signup)
	mux.HandleFunc("POST /login", handlers.Login)

	// Routes with JWT middleware
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("POST /lots", handlers.CreateLot)
	protectedMux.HandleFunc("PATCH /lots/", handlers.UpdateLot)
	protectedMux.HandleFunc("POST /auctions", handlers.CreateAuction)
	protectedMux.HandleFunc("POST /auctions/{id}/bids", handlers.CreateBid)
	protectedMux.HandleFunc("GET /auctions/{id}/bids", handlers.ListBids)
	protectedMux.HandleFunc("DELETE /lots/", handlers.DeleteLot)

	// Apply middlewares to protected routes
	handler := handlers.CorsMiddleware(handlers.JwtAuthMiddleware(protectedMux))

	// Combine all routes
	mux.Handle("/", handler)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Starting server on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}