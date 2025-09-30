package lot

import (
	"encoding/json"
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
		http.Error(w, "Only sellers can create lots", http.StatusForbidden)
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

	id, err := database.CreateLot(lot)
	if err != nil {
		http.Error(w, "Failed to create lot", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.IDResponse{ID: id})
}