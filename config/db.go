package config

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectPostgres() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Database environment variables are missing")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
}
