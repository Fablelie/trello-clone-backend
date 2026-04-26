package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

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
		getEnv("DB_USER", "myuser"),
		getEnv("DB_PASSWORD", "mypassword"),
		getEnv("DB_NAME", "mydatabase"),
		getEnv("DB_PORT", "5432"),
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

	// Assemble task module
	taskRepo := postgresRepo.NewTaskRepository(db)
	taskUsecase := usecase.NewTaskUsecase(taskRepo, projectRepo)
	taskHandler := handler.NewTaskHandler(taskUsecase)

	// Initialize Fiber and setup Router
	app := fiber.New()
	http.SetupRouter(app, userHandler, projectHandler, taskHandler, jwtSecret)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		port := getEnv("PORT", "8080")
		log.Printf("Server is starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			log.Printf("Listen error: %v", err)
		}
	}()

	<-c
	log.Println("Gracefully shutting down...")

	if err := app.Shutdown(); err != nil {
		log.Printf("Fiber shutdown error: %v", err)
	}

	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	log.Println("Server cleanup completed. Goodbye!")
}
