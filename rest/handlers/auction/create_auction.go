package auction

import (
	"encoding/json"
	"errors"
	"net/http"

	"banana-auction/database"
	"banana-auction/models"
	"banana-auction/rest/middlewares"
)

func Create(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Only sellers can create auctions", http.StatusForbidden)
		return
	}

	var req models.CreateAuctionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lot, err := database.GetLot(req.LotID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			http.Error(w, "Lot not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if lot.SellerID != userID {
		http.Error(w, "Unauthorized to create auction for this lot", http.StatusForbidden)
		return
	}

	exists, err := database.AuctionExistsForLot(req.LotID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "Auction already exists for this lot", http.StatusBadRequest)
		return
	}

	auction := models.Auction{
		LotID:             req.LotID,
		StartDate:         req.StartDate,
		DurationDays:      req.DurationDays,
		InitialPricePerKG: req.InitialPricePerKG,
	}

	id, err := database.CreateAuction(auction)
	if err != nil {
		http.Error(w, "Failed to create auction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.IDResponse{ID: id})
}
