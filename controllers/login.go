package controllers

import (
	"encoding/json"
	"net/http"
	"parking_management_system_backend/config"
	"parking_management_system_backend/helpers"
	models "parking_management_system_backend/user"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var user models.User
	err := config.DB.Get(&user, "SELECT id, name, email, role, password FROM users WHERE email=$1 LIMIT 1", req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	if !helpers.VerifyPassword(user.Password, req.Password) {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, _ := helpers.GenerateTokens(user.ID.String(), user.Role)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
