package main

import (
	"log"
	"os"

	"github.com/fablelie/trello-clone-backend/internal/delivery/http"
	"github.com/fablelie/trello-clone-backend/internal/delivery/http/handler"
	"github.com/fablelie/trello-clone-backend/internal/infrastructure/database"
	postgresRepo "github.com/fablelie/trello-clone-backend/internal/repository/postgres"
	"github.com/fablelie/trello-clone-backend/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func main() {
	// load Config from .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	db := database.NewPostgresDB(
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "5432"),
		getEnv("DB_PASSWORD", "myuser"),
		getEnv("DB_NAME", "mypassword"),
		getEnv("DB_PORT", "mydatabase"),
	)

	jwtSecret := getEnv("JWT_SECRET", "secret_key")

	// Assemble user module
	userRepo := postgresRepo.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, jwtSecret)
	userHandler := handler.NewUserHandler(userUsecase)

	// Assemble project module
	projectRepo := postgresRepo.NewProjectRepository(db)
	projectUsecase := usecase.NewProjectUsecase(projectRepo)
	projectHandler := handler.NewProjectHandler(projectUsecase)

	// Initialize Fiber and setup Router
	app := fiber.New()
	http.SetupRouter(app, userHandler, projectHandler, jwtSecret)

	port := getEnv("PORT", "8080")

	log.Printf("Server is running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
