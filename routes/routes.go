package routes

import (
	"net/http"
	"strings"

	"parking_management_system_backend/controllers"
	"parking_management_system_backend/helpers"
)

func RegisterRoutes() {
	// Public routes
	http.HandleFunc("/api/signup", controllers.Signup)
	http.HandleFunc("/api/login", controllers.Login)
	http.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message":"pong"}`))
	})

	// Protected routes
	http.Handle("/api/users", helpers.AuthMiddleware(
		helpers.AuthorizeRole(http.HandlerFunc(controllers.GetUsers), "ADMIN", "SUPER_ADMIN"),
	))

	http.Handle("/api/users/", helpers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/api/users/")
		id = strings.TrimPrefix(id, "/") // remove leading slash if exists

		if id == "" {
			http.Error(w, "User ID required", http.StatusBadRequest)
			return
		}

		r.Header.Set("userID", id)
		controllers.GetUser(w, r)
	})))

}
