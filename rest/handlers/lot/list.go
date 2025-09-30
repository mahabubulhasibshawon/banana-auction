package lot

import (
	"encoding/json"
	"net/http"

	"banana-auction/database"
	"banana-auction/rest/middlewares"
)

func List(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Only sellers can list lots", http.StatusForbidden)
		return
	}

	lots, err := database.ListLots()
	if err != nil {
		http.Error(w, "Failed to list lots", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lots)
}
