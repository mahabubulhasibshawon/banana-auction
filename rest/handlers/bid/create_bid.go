package bid

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"

	"banana-auction/database"
	"banana-auction/models"
	"banana-auction/rest/middlewares"
)

func Create(w http.ResponseWriter, r *http.Request) {
	auctionIDStr := path.Base(path.Dir(r.URL.Path))
	auctionID, err := strconv.Atoi(auctionIDStr)
	if err != nil {
		http.Error(w, "Invalid auction ID", http.StatusBadRequest)
		return
	}

	userID, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := database.GetUser(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user.Role != "buyer" {
		http.Error(w, "Only buyers can place bids", http.StatusForbidden)
		return
	}

	var req models.CreateBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = database.GetAuction(auctionID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Auction not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	bid := models.Bid{
		AuctionID:     auctionID,
		BuyerID:       userID,
		BidPricePerKG: req.BidPricePerKG,
	}

	id, err := database.CreateBid(bid)
	if err != nil {
		http.Error(w, "Failed to create bid", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.IDResponse{ID: id})
}
