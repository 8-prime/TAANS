package main

import (
	"log"
	"taans/internal/app"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app, err := app.NewApplication()
	if err != nil {
		log.Fatal(err)
	}

	app.RegisterRoutes()

	if err := app.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
