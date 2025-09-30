package auction

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"

	"banana-auction/database"
	"banana-auction/rest/middlewares"
)

func List(w http.ResponseWriter, r *http.Request) {
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
	if user.Role != "seller" {
		http.Error(w, "Only sellers can list bids", http.StatusForbidden)
		return
	}

	auction, err := database.GetAuction(auctionID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Auction not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	lot, err := database.GetLot(auction.LotID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if lot.SellerID != userID {
		http.Error(w, "Unauthorized to view bids for this auction", http.StatusForbidden)
		return
	}

	bids, err := database.ListBids(auctionID)
	if err != nil {
		http.Error(w, "Failed to list bids", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bids)
}
