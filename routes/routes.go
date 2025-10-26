package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"parking_management_system_backend/controllers"
	"parking_management_system_backend/helpers"
	"parking_management_system_backend/services"

	"github.com/gorilla/mux"
)

func SetupRoutes() {
	r := mux.NewRouter()

	// ----------------------
	// Public routes
	// ----------------------
	r.HandleFunc("/api/signup", controllers.Signup).Methods("POST")
	r.HandleFunc("/api/login", controllers.Login).Methods("POST")

	// ----------------------
	// Protected routes (JWT)
	// ----------------------
	r.Handle("/api/profile", helpers.AuthMiddleware(http.HandlerFunc(controllers.GetProfile))).Methods("GET")
	r.Handle("/api/update-profile", helpers.AuthMiddleware(http.HandlerFunc(controllers.UpdateProfile))).Methods("POST")

	// ----------------------
	// Google OAuth2
	// ----------------------
	r.HandleFunc("/api/oauth/google/login", func(w http.ResponseWriter, r *http.Request) {
		state := "random_state_string" // TODO: generate secure random state
		url := services.GetGoogleLoginURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}).Methods("GET")

	r.HandleFunc("/api/oauth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		user, accessToken, refreshToken, err := services.HandleGoogleCallback(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Serialize user to JSON string
		userJSON, _ := json.Marshal(user)
		userEncoded := url.QueryEscape(string(userJSON))

		// 🔹 Redirect to React OAuthRedirect component
		frontendURL := fmt.Sprintf(
			"http://localhost:3000/oauth2/redirect?access_token=%s&refresh_token=%s&user=%s",
			accessToken, refreshToken, userEncoded,
		)
		http.Redirect(w, r, frontendURL, http.StatusSeeOther)
	}).Methods("GET")

	// ----------------------
	// Slot management routes
	// All users can view slots
	// r.Handle("/api/slots", helpers.AuthMiddleware(http.HandlerFunc(controllers.GetSlots))).Methods("GET")

	// // Only users with role "user" can book or release a slot
	// r.Handle("/api/book-slot", helpers.AuthMiddleware(
	// 	helpers.AuthorizeRole(http.HandlerFunc(controllers.BookSlot), "user"),
	// )).Methods("POST")

	// r.Handle("/api/release-slot", helpers.AuthMiddleware(
	// 	helpers.AuthorizeRole(http.HandlerFunc(controllers.ReleaseSlot), "user"),
	// )).Methods("POST")

	// // Admin-only notification update
	// r.Handle("/api/slots/{id}/notify", helpers.AuthMiddleware(
	// 	helpers.AuthorizeRole(http.HandlerFunc(controllers.NotifySlot), "admin"),
	// )).Methods("PUT")

	// Register router
	http.Handle("/", r)
}
