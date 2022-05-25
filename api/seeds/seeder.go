package main

import (
	"log"

	"github.com/ko1ke/know-sync-api/datasources/postgres_db"
	"github.com/ko1ke/know-sync-api/seeds/seeds"
)

func main() {
	dbClient := postgres_db.Client
	for _, seed := range seeds.All() {
		if err := seed.Run(dbClient); err != nil {
			log.Printf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
}
