package controllers

import (
	"encoding/json"
	"net/http"
	"parking_management_system_backend/config"
	models "parking_management_system_backend/user"
)

// GetUser returns a single user by ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("userID") // set in the route from URL
	if userID == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	var user models.User
	err := config.DB.Get(&user, `
		SELECT id, name, email, phone, address, role, status, subscription_status, created_at, updated_at
		FROM users
		WHERE id=$1
	`, userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
