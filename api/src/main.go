package main

import (
	"know-sync-api/app"

	"github.com/sirupsen/logrus"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Printf("failed to load .env: %v", err)
	}
}

func main() {
	logrus.SetReportCaller(true)
	app.StartApp()
}
