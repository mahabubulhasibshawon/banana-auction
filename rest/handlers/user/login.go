package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"banana-auction/config"
	"banana-auction/database"
	"banana-auction/models"
	"banana-auction/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := database.GetUserByUsername(req.Username)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
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

	token, err := utils.GenerateJWT(user.ID, config.GetConfig().JwtSecretKey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.LoginResponse{Token: token})
}
