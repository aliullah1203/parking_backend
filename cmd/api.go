package main

import (
	"log"
	"net/http"
	"os"
	"parking_management_system_backend/config"
	"parking_management_system_backend/controllers"

	"github.com/joho/godotenv"
)

func Api() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Connect to PostgreSQL
	config.ConnectPostgres()

	// Routes
	http.HandleFunc("/api/signup", controllers.Signup)
	http.HandleFunc("/api/login", controllers.Login) // if you have login

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
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
