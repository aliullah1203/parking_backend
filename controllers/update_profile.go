package controllers

import (
	"encoding/json"
	"net/http"
	"parking_management_system_backend/config"
	"parking_management_system_backend/helpers"
)

type User struct {
	ID      string `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Email   string `db:"email" json:"email"`
	Phone   string `db:"phone" json:"phone"`
	License string `db:"license" json:"license"`
	NID     string `db:"nid" json:"nid"`
	Picture string `db:"picture" json:"picture"`
	Role    string `db:"role" json:"role"`
}

// UpdateProfile only updates license, NID, and picture
func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
		return
	}

	claims := helpers.GetClaimsFromContext(r)
	if claims == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"message": "Unauthorized"})
		return
	}
	userID := claims.UserID

	var updated struct {
		License string `json:"license"`
		NID     string `json:"nid"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
		return
	}

	query := `
		UPDATE users
		SET license=$1, nid=$2, picture=$3, updated_at=NOW()
		WHERE id=$4
		RETURNING id, name, email, phone, license, nid, picture, role
	`

	var user User
	err := config.DB.Get(&user, query, updated.License, updated.NID, updated.Picture, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "Failed to update profile: " + err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Profile updated successfully",
		"user":    user,
	})
}
