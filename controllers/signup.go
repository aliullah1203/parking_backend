package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"parking_management_system_backend/config"
	"parking_management_system_backend/helpers"
	models "parking_management_system_backend/user"
	"time"

	"github.com/google/uuid"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Contact         string `json:"contact"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Password != req.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	var count int
	err := config.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1 OR phone=$2", req.Email, req.Contact)
	if err != nil {
		log.Println("DB query error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Email or phone already exists", http.StatusBadRequest)
		return
	}

	user := models.User{
		ID:                 uuid.New(),
		Name:               req.Name,
		Email:              req.Email,
		Phone:              req.Contact,
		Role:               "CUSTOMER",
		Status:             "ACTIVE",
		SubscriptionStatus: "SUBSCRIBED",
		Password:           req.Password,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	_, err = config.DB.NamedExec(`
		INSERT INTO users 
		(id, name, email, phone, role, status, subscription_status, password, created_at, updated_at)
		VALUES (:id,:name,:email,:phone,:role,:status,:subscription_status,:password,:created_at,:updated_at)
	`, &user)
	if err != nil {
		log.Println("DB insert error:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := helpers.GenerateTokens(user.ID.String(), user.Role)
	if err != nil {
		log.Println("Token generation error:", err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "User created successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
