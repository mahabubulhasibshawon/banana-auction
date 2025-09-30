package lot

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

func Update(w http.ResponseWriter, r *http.Request) {
	idStr := path.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid lot ID", http.StatusBadRequest)
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
		http.Error(w, "Only sellers can update lots", http.StatusForbidden)
		return
	}

	var req models.UpdateLotRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	lot, err := database.GetLot(id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
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
	if err := database.UpdateLot(lot); err != nil {
		http.Error(w, "Failed to update lot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
