package main

import (
	"log"
	"net/http"
	"os"
	"parking_management_system_backend/config"
	"parking_management_system_backend/helpers"
	"parking_management_system_backend/routes"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	helpers.InitJWTSecret()

	log.Println("DB_HOST from .env:", os.Getenv("DB_HOST"))

	// Connect to PostgreSQL
	config.ConnectPostgres()

	// Register routes
	routes.SetupRoutes()

	// CORS middleware
	handler := corsMiddleware(http.DefaultServeMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server started at port", port)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}

// Simple CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
