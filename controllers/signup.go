package controllers

import (
	"encoding/json"
	"net/http"
	"parking_management_system_backend/config"
	"parking_management_system_backend/helpers"
	models "parking_management_system_backend/user"
	"time"

	"github.com/google/uuid"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check duplicate
	var count int
	err := config.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1 OR phone=$2", user.Email, user.Phone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Email or phone already exists", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New()
	user.Role = "CUSTOMER"
	user.Status = "ACTIVE"
	user.SubscriptionStatus = "SUBSCRIBED"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err = config.DB.NamedExec(`INSERT INTO users 
	(id, name, email, phone, address, role, status, subscription_status, password, created_at, updated_at) 
	VALUES (:id,:name,:email,:phone,:address,:role,:status,:subscription_status,:password,:created_at,:updated_at)`, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, _ := helpers.GenerateTokens(user.ID.String(), user.Role)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "User created successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
