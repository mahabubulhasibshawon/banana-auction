package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"
	"strconv"

	"banana-auction/db"
	"banana-auction/models"
	"banana-auction/utils"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" || req.Name == "" || req.Role == "" {
		http.Error(w, "Username, password, name, and role are required", http.StatusBadRequest)
		return
	}

	if req.Role != "seller" && req.Role != "buyer" {
		http.Error(w, "Invalid role: must be 'seller' or 'buyer'", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Name:         req.Name,
		Role:         req.Role,
	}

	id, err := db.CreateUser(user)
	if err != nil {
		if errors.Is(err, db.ErrDuplicateUsername) {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.IDResponse{ID: id})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.LoginResponse{Token: token})
}

func CreateLot(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := db.GetUser(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user.Role != "seller" {
		http.Error(w, "Unauthorized: only sellers can create lots", http.StatusForbidden)
		return
	}

	var req models.CreateLotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.TotalWeightKG < 1000 {
		http.Error(w, "Minimum weight allowed is 1000 kg", http.StatusBadRequest)
		return
	}

	lot := models.Lot{
		SellerID:       userID,
		Cultivar:       req.Cultivar,
		PlantedCountry: req.PlantedCountry,
		HarvestDate:    req.HarvestDate,
		TotalWeightKG:  req.TotalWeightKG,
	}

	id, err := db.CreateLot(lot)
	if err != nil {
		http.Error(w, "Failed to create lot", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.IDResponse{ID: id})
}

func UpdateLot(w http.ResponseWriter, r *http.Request) {
	idStr := path.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid lot ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := db.GetUser(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user.Role != "seller" {
		http.Error(w, "Unauthorized: only sellers can update lots", http.StatusForbidden)
		return
	}

	var req models.UpdateLotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lot, err := db.GetLot(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Lot not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if lot.SellerID != userID {
		http.Error(w, "Unauthorized to update this lot", http.StatusForbidden)
		return
	}

	lot.HarvestDate = req.HarvestDate
	if err := db.UpdateLot(lot); err != nil {
		http.Error(w, "Failed to update lot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func CreateAuction(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := db.GetUser(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user.Role != "seller" {
		http.Error(w, "Unauthorized: only sellers can create auctions", http.StatusForbidden)
		return
	}

	var req models.CreateAuctionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lot, err := db.GetLot(req.LotID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
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

	exists, err := db.AuctionExistsForLot(req.LotID)
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

	id, err := db.CreateAuction(auction)
	if err != nil {
		http.Error(w, "Failed to create auction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.IDResponse{ID: id})
}

func CreateBid(w http.ResponseWriter, r *http.Request) {
	auctionIDStr := path.Base(path.Dir(r.URL.Path))
	auctionID, err := strconv.Atoi(auctionIDStr)
	if err != nil {
		http.Error(w, "Invalid auction ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := db.GetUser(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user.Role != "buyer" {
		http.Error(w, "Unauthorized: only buyers can place bids", http.StatusForbidden)
		return
	}

	var req models.CreateBidRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err = db.GetAuction(auctionID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
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

	id, err := db.CreateBid(bid)
	if err != nil {
		http.Error(w, "Failed to create bid", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.IDResponse{ID: id})
}

func ListBids(w http.ResponseWriter, r *http.Request) {
	auctionIDStr := path.Base(path.Dir(r.URL.Path))
	auctionID, err := strconv.Atoi(auctionIDStr)
	if err != nil {
		http.Error(w, "Invalid auction ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := db.GetUser(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user.Role != "seller" {
		http.Error(w, "Unauthorized: only sellers can list bids", http.StatusForbidden)
		return
	}

	auction, err := db.GetAuction(auctionID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Auction not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	lot, err := db.GetLot(auction.LotID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if lot.SellerID != userID {
		http.Error(w, "Unauthorized to view bids for this auction", http.StatusForbidden)
		return
	}

	bids, err := db.ListBids(auctionID)
	if err != nil {
		http.Error(w, "Failed to list bids", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bids)
}

func DeleteLot(w http.ResponseWriter, r *http.Request) {
	idStr := path.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid lot ID", http.StatusBadRequest)
		return
	}

	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := db.GetUser(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user.Role != "seller" {
		http.Error(w, "Unauthorized: only sellers can delete lots", http.StatusForbidden)
		return
	}

	lot, err := db.GetLot(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			http.Error(w, "Lot not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if lot.SellerID != userID {
		http.Error(w, "Unauthorized to delete this lot", http.StatusForbidden)
		return
	}

	if err := db.DeleteLot(id); err != nil {
		http.Error(w, "Failed to delete lot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}