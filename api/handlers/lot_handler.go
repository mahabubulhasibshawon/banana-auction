
package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"banana-auction/internal/domain/lot"
	"banana-auction/internal/domain/user"
	"banana-auction/api/middlewares"
)

type LotHandler struct {
	svc    lot.Service
	userSvc user.Service
}

func NewLotHandler(svc lot.Service, userSvc user.Service) *LotHandler {
	return &LotHandler{svc: svc, userSvc: userSvc}
}

func (h *LotHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Cultivar       string `json:"cultivar"`
		PlantedCountry string `json:"planted_country"`
		HarvestDate    string `json:"harvest_date"`
		TotalWeightKG  int    `json:"total_weight_kg"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.svc.CreateLot(userID, req.Cultivar, req.PlantedCountry, req.HarvestDate, req.TotalWeightKG)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"lot id": id})
}

func (h *LotHandler) Update(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid lot ID", http.StatusBadRequest)
		return
	}

	userID, err := middlewares.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		HarvestDate string `json:"harvest_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateLot(id, userID, req.HarvestDate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LotHandler) Delete(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	idStr := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid lot ID", http.StatusBadRequest)
		return
	}

	userID, err := middlewares.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.svc.DeleteLot(id, userID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *LotHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, err := middlewares.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Fetch the user to check their role
	user, err := h.userSvc.GetUser(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// Check if the user is a seller
	if user.Role != "seller" {
		http.Error(w, "Only sellers can list lots", http.StatusForbidden)
		return
	}

	// List all lots
	lots, err := h.svc.ListLots()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(lots)
}