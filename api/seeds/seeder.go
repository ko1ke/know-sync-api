package main

import (
	"know-sync-api/datasources/postgres_db"
	"know-sync-api/seeds/seeds"
	"log"
)

func main() {
	dbClient := postgres_db.Client
	for _, seed := range seeds.All() {
		if err := seed.Run(dbClient); err != nil {
			log.Printf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
}
