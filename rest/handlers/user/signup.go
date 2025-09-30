package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"banana-auction/database"
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
		http.Error(w, "Role must be 'seller' or 'buyer'", http.StatusBadRequest)
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

	id, err := database.CreateUser(user)
	if err != nil {
		if errors.Is(err, database.ErrDuplicateUsername) {
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
