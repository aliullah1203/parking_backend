package cmd

import (
	"log"
	"net/http"
	"os"

	"parking_management_system_backend/config"
	"parking_management_system_backend/routes"

	"github.com/joho/godotenv"
)

// Api initializes environment, DB, and starts the HTTP server
func Api() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Initialize PostgreSQL connection
	config.ConnectPostgres()

	// Register routes
	routes.RegisterRoutes()

	// Get port from env or default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
