package postgres_db

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Client *gorm.DB
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env: %v", err)
	}
	dsnString := os.Getenv("DATABASE_URL")
	client, err := gorm.Open(postgres.Open(dsnString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	log.Println("database successfully configure")
	Client = client
}
