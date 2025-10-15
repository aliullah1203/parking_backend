package routes

import (
	"net/http"
	"parking_management_system_backend/controllers"
)

// SetupRoutes registers HTTP routes without Gin
func SetupRoutes() {
	http.HandleFunc("/api/signup", controllers.Signup)
	http.HandleFunc("/api/login", controllers.Login)
	// add more routes here
}
