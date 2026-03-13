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

var DB *gorm.DB

func InitDB() {
	// 1. Load the .env file into the process environment
	if err := godotenv.Load(); err != nil {
		log.Println("[!] No .env file found, using system environment variables")
	}

	// 2. Now os.Getenv will actually find your values
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		host, user, pass, name)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("[-] Database connection failed. Check your .env file and ensure the Docker container is running.")
	}

	DB.AutoMigrate(&models.Agent{})
	fmt.Println("[+] Secure Persistence Bridge established.")
}
