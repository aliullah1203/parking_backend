package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"parking_management_system_backend/controllers"
	"parking_management_system_backend/helpers"
	"parking_management_system_backend/services"
)

func SetupRoutes() {
	// Public routes
	http.Handle("/api/signup", http.HandlerFunc(controllers.Signup))
	http.Handle("/api/login", http.HandlerFunc(controllers.Login))

	// Protected routes (JWT)
	http.Handle("/api/update-profile", helpers.AuthMiddleware(http.HandlerFunc(controllers.UpdateProfile)))
	http.Handle("/api/profile", helpers.AuthMiddleware(http.HandlerFunc(controllers.GetProfile)))

	// Google OAuth2 login
	http.HandleFunc("/api/oauth/google/login", func(w http.ResponseWriter, r *http.Request) {
		state := "random_state_string" // You can generate a secure random string
		url := services.GetGoogleLoginURL(state)
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	})

	// Google OAuth2 callback
	http.HandleFunc("/api/oauth/google/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		user, accessToken, refreshToken, err := services.HandleGoogleCallback(code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Serialize user to JSON string
		userJSON, _ := json.Marshal(user)
		// URL-encode the JSON string
		userEncoded := url.QueryEscape(string(userJSON))

		// Redirect to frontend with tokens and user info
		frontendURL := fmt.Sprintf(
			"http://localhost:3000/dashboard?access_token=%s&refresh_token=%s&user=%s",
			accessToken, refreshToken, userEncoded,
		)
		http.Redirect(w, r, frontendURL, http.StatusSeeOther)
	})
}
