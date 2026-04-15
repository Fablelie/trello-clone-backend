package main

import (
	"log"
	"os"

	"github.com/fablelie/trello-clone-backend/internal/infrastructure/database"
	"github.com/joho/godotenv"
)

func main() {
	// load Config from .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	_ = database.NewPostgresDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	log.Println("Migration completed successfully!")
}
