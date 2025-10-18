package controllers

import (
	"encoding/json"
	"net/http"
	"parking_management_system_backend/config"
	"parking_management_system_backend/helpers"
)

// GetProfile returns the logged-in user's profile
func GetProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	claims := helpers.GetClaimsFromContext(r)
	if claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}
	userID := claims.UserID

	var user User
	query := `SELECT id, name, email, phone, license, nid, picture, role FROM users WHERE id=$1`
	err := config.DB.Get(&user, query, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to fetch profile: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(user)
}
