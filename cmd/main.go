package main

import (
	"github.com/joho/godotenv"
	"github.com/ko1ke/know-sync-api/cmd/app"
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
	app.StartApp()
}
