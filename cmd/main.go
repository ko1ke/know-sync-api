package main

import (
	"github.com/joho/godotenv"
	"github.com/ko1ke/know-sync-api/cmd/app"
	"github.com/ko1ke/know-sync-api/cmd/datasources/postgres_db"
	"github.com/ko1ke/know-sync-api/cmd/datasources/redis_db"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Printf("failed to load .env: %v", err)
	}
}

func main() {
	logrus.SetReportCaller(true)
	postgres_db.Start()
	redis_db.Start()
	app.StartApp()
}
