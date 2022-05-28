package postgres_db

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	Client *gorm.DB
)

func Start() {
	dsnString := os.Getenv("DATABASE_URL")
	client, err := gorm.Open(postgres.Open(dsnString), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	log.Println("database successfully configure")
	Client = client
}
