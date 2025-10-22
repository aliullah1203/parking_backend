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
	"golang.org/x/crypto/bcrypt"
)

// Signup creates a new user
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Name            string `json:"name"`
		Email           string `json:"email"`
		Phone           string `json:"phone"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("Signup decode error:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Password != req.ConfirmPassword {
		http.Error(w, "Passwords do not match", http.StatusBadRequest)
		return
	}

	// Check if email or phone already exists
	var count int
	err := config.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1 OR phone=$2", req.Email, req.Phone)

	if err != nil {
		log.Println("Signup DB query error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "Email or phone already exists", http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Signup password hashing error:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	user := models.User{
		ID:                 uuid.New().String(),
		Name:               req.Name,
		Email:              req.Email,
		Phone:              req.Phone,
		Role:               "CUSTOMER",
		Status:             "ACTIVE",
		SubscriptionStatus: "SUBSCRIBED",
		Password:           string(hashedPassword),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	_, err = config.DB.NamedExec(`
		INSERT INTO users (id, name, email, phone, role, status, subscription_status, password, created_at, updated_at)
		VALUES (:id,:name,:email,:phone,:role,:status,:subscription_status,:password,:created_at,:updated_at)
	`, &user)
	if err != nil {
		log.Println("Signup DB insert error:", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := helpers.GenerateTokens(user.ID, user.Role)
	if err != nil {
		log.Println("Signup token generation error:", err)
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "User created successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
