package db

import (
	"fmt"
	"log"
	"os"

	"github.com/NotKidding/olympus-server/internal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database instance used across the application
var DB *gorm.DB

// InitDB initializes the connection to PostgreSQL and performs auto-migration
func InitDB() {
	// 1. Load the .env file.
	// In Docker, environment variables are passed directly, so we ignore the error
	// if the file isn't present.
	if err := godotenv.Load(); err != nil {
		log.Println("[!] No .env file found, using system environment variables")
	}

	// 2. Extract database credentials from the environment
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Fallback to default port if not specified
	if port == "" {
		port = "5432"
	}

	// 3. Construct the Data Source Name (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port)

	// 4. Open the connection using GORM
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("[-] Database connection failed: %v. Check your .env file and ensure the Docker container is running.", err)
	}

	// 5. Auto-Migration
	// This creates or updates tables based on your Go structs.
	// We include models.Task{} here to officially kick off Phase 3.
	err = DB.AutoMigrate(&models.Agent{}, &models.Task{})
	if err != nil {
		log.Fatalf("[-] Database migration failed: %v", err)
	}

	fmt.Println("[+] Secure Persistence Bridge established (Agents & Tasks).")
}
