package db

import (
	"fmt"
	"log"
	"os" // Required to read Environment Variables

	"github.com/NotKidding/olympus-server/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Pull secrets from the environment
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable",
		host, user, pass, name)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("[-] Database connection failed. Check your environment variables.")
	}

	DB.AutoMigrate(&models.Agent{})
	fmt.Println("[+] Secure Persistence Bridge established.")
}
