package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"banana-auction/api/middlewares"
	"banana-auction/internal/domain/bid"
)

type BidHandler struct {
	svc bid.Service
}

func NewBidHandler(svc bid.Service) *BidHandler {
	return &BidHandler{svc: svc}
}

func (h *BidHandler) PlaceBid(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract auctionID from the path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL format. Use /auctions/{auctionID}/bids", http.StatusBadRequest)
		return
	}

	lastSegment := pathParts[len(pathParts)-1]
	if lastSegment != "bids" {
		http.Error(w, "Invalid URL format. Expected /auctions/{auctionID}/bids", http.StatusBadRequest)
		return
	}

	auctionIDStr := pathParts[len(pathParts)-2]

	auctionID, err := strconv.Atoi(auctionIDStr)

	if err != nil {
		http.Error(w, "Invalid auction ID", http.StatusBadRequest)
		return
	}

	userID, err := middlewares.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		BidPricePerKG float64 `json:"bid_price_per_kg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.svc.PlaceBid(auctionID, userID, req.BidPricePerKG)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"Bid placed successfully!!! bid_id": id})
}
