package controllers

import (
	"encoding/json"
	"net/http"
	"parking_management_system_backend/config"
	models "parking_management_system_backend/user"
)

// GetUsers returns all users (ADMIN/SUPER_ADMIN only)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	// Role validation is handled by middleware, so we can just fetch users
	var users []models.User
	err := config.DB.Select(&users, `
		SELECT id, name, email, phone, address, role, status, subscription_status, created_at, updated_at
		FROM users
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
