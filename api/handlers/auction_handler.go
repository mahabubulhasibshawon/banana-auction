package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"banana-auction/api/middlewares"
	"banana-auction/internal/domain/auction"
	"banana-auction/internal/domain/bid"
	"banana-auction/internal/domain/lot"
)

type AuctionHandler struct {
	svc    auction.Service
	lotSvc lot.Service
	bidSvc bid.Service
}

func NewAuctionHandler(svc auction.Service, lotSvc lot.Service, bidSvc bid.Service) *AuctionHandler {
	return &AuctionHandler{svc: svc, lotSvc: lotSvc, bidSvc: bidSvc}
}

func (h *AuctionHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		LotID             int     `json:"lot_id"`
		StartDate         string  `json:"start_date"`
		DurationDays      int     `json:"duration_days"`
		InitialPricePerKG float64 `json:"initial_price_per_kg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch the lot to verify the seller
	lot, err := h.lotSvc.GetLot(req.LotID)
	if err != nil {
		http.Error(w, "Lot not found", http.StatusNotFound)
		return
	}

	// Check if the user is the seller of the lot
	if lot.SellerID != userID {
		http.Error(w, "Only the seller of the lot can create an auction", http.StatusForbidden)
		return
	}

	// Create the auction
	id, err := h.svc.CreateAuction(req.LotID, req.StartDate, req.DurationDays, req.InitialPricePerKG)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *AuctionHandler) ListBids(w http.ResponseWriter, r *http.Request) {
	// Extract auctionID from the path (e.g., /auctions/7/bids/)
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

	// Fetch the auction
	auction, err := h.svc.GetAuction(auctionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Fetch the associated lot
	lot, err := h.lotSvc.GetLot(auction.LotID)
	if err != nil {
		http.Error(w, "Lot not found", http.StatusInternalServerError)
		return
	}

	// Check if the user is the seller of the lot
	if lot.SellerID != userID {
		http.Error(w, "Only the seller can list bids for this auction", http.StatusForbidden)
		return
	}

	// List bids for the auction using bid service
	bids, err := h.bidSvc.ListBids(auctionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(bids)
}

func (h *AuctionHandler) GetAuction(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := pathParts[len(pathParts)-1]
	auctionID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid auction ID", http.StatusBadRequest)
		return
	}

	userID, err := middlewares.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch the auction
	auction, err := h.svc.GetAuction(auctionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Fetch the associated lot
	lot, err := h.lotSvc.GetLot(auction.LotID)
	if err != nil {
		http.Error(w, "Lot not found", http.StatusInternalServerError)
		return
	}

	// Check if the user is the seller of the lot
	if lot.SellerID != userID {
		http.Error(w, "Only the seller can view this auction", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(auction)
}
