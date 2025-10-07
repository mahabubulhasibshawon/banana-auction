
package api

import (
	"net/http"

	"banana-auction/api/handlers"
	"banana-auction/api/middlewares"
	"banana-auction/internal/domain/auction"
	"banana-auction/internal/domain/bid"
	"banana-auction/internal/domain/lot"
	"banana-auction/internal/domain/user"
	"banana-auction/internal/infrastructure/persistence/postgres"
)

func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	userSvc := user.NewService(postgres.NewUserRepo(postgres.GetDB()))
	userHandler := handlers.NewUserHandler(userSvc)

	lotSvc := lot.NewService(postgres.NewLotRepo(postgres.GetDB()))
	lotHandler := handlers.NewLotHandler(lotSvc, userSvc)

	auctionSvc := auction.NewService(postgres.NewAuctionRepo(postgres.GetDB()))
	bidSvc := bid.NewService(postgres.NewBidRepo(postgres.GetDB()))
	auctionHandler := handlers.NewAuctionHandler(auctionSvc, lotSvc, bidSvc)

	bidHandler := handlers.NewBidHandler(bidSvc)

	// Public routes
	mux.Handle("POST /signup",http.HandlerFunc(userHandler.Signup))
	mux.Handle("POST /login",http.HandlerFunc(userHandler.Login))

	// Protected routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/lots", lotHandler.Create)
	protectedMux.HandleFunc("/lots/{id}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPatch {
			lotHandler.Update(w, r)
		} else if r.Method == http.MethodDelete {
			lotHandler.Delete(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	protectedMux.HandleFunc("/auctions", auctionHandler.Create)
	protectedMux.HandleFunc("/auctions/{id}", auctionHandler.GetAuction)
	protectedMux.HandleFunc("/auctions/{id}/bids", auctionHandler.ListBids)
	protectedMux.HandleFunc("/auctions/{id}/bids/", bidHandler.PlaceBid) 

	protectedHandler := middlewares.JwtAuthMiddleware(protectedMux)

	mux.Handle("/", protectedHandler)

	return middlewares.CorsMiddleware(mux)
}