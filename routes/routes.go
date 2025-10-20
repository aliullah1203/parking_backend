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

    // Public routes
    r.HandleFunc("/api/signup", controllers.Signup).Methods("POST")
    r.HandleFunc("/api/login", controllers.Login).Methods("POST")

    // Protected routes (JWT)
    r.Handle("/api/update-profile", helpers.AuthMiddleware(http.HandlerFunc(controllers.UpdateProfile))).Methods("PUT")
    r.Handle("/api/profile", helpers.AuthMiddleware(http.HandlerFunc(controllers.GetProfile))).Methods("GET")

    // Google OAuth2 login
    r.HandleFunc("/api/oauth/google/login", func(w http.ResponseWriter, r *http.Request) {
        state := "random_state_string" // TODO: generate a secure random state
        url := services.GetGoogleLoginURL(state)
        http.Redirect(w, r, url, http.StatusTemporaryRedirect)
    }).Methods("GET")

    // Google OAuth2 callback
    r.HandleFunc("/api/oauth/google/callback", func(w http.ResponseWriter, r *http.Request) {
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
    }).Methods("GET")

    // Slot management routes
    r.Handle("/api/slots", helpers.AuthMiddleware(http.HandlerFunc(controllers.GetSlots))).Methods("GET")
    r.Handle("/api/book-slot", helpers.AuthMiddleware(http.HandlerFunc(controllers.BookSlot))).Methods("POST")
    r.Handle("/api/release-slot", helpers.AuthMiddleware(http.HandlerFunc(controllers.ReleaseSlot))).Methods("POST")
    r.Handle("/api/slots/{id}/notify", helpers.AuthMiddleware(http.HandlerFunc(controllers.NotifySlot))).Methods("PUT")

    // Register the router
    http.Handle("/", r)
}
