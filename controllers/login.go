package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"parking_management_system_backend/config"
	"parking_management_system_backend/helpers"
	models "parking_management_system_backend/user"

	"golang.org/x/crypto/bcrypt"
)

// Login handles user login
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode request body
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Login decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch user by email
	var user models.User
	query := `SELECT id, name, email, phone, role, password FROM users WHERE email=$1 LIMIT 1`
	err := config.DB.Get(&user, query, req.Email)
	if err != nil {
		log.Println("Login DB query error:", err)
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := helpers.GenerateTokens(user.ID, user.Role)
	if err != nil {
		log.Println("Login token generation error:", err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	// Return user info with tokens
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":      user.ID,
			"name":    user.Name,
			"email":   user.Email,
			"contact": user.Phone,
			"role":    user.Role,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
