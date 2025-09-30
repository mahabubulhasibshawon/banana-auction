package lot

import (
	"errors"
	"net/http"
	"path"
	"strconv"

	"banana-auction/database"
	"banana-auction/rest/middlewares"
)

func DeleteLot(w http.ResponseWriter, r *http.Request) {
	idStr := path.Base(r.URL.Path)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid lot ID", http.StatusBadRequest)
		return
	}

	userID, err := middlewares.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := database.GetUser(userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if user.Role != "seller" {
		http.Error(w, "Unauthorized: only sellers can delete lots", http.StatusForbidden)
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
		http.Error(w, "Unauthorized to delete this lot", http.StatusForbidden)
		return
	}

	if err := database.DeleteLot(id); err != nil {
		http.Error(w, "Failed to delete lot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
