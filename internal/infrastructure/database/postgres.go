package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	my_postgres "github.com/fablelie/trello-clone-backend/internal/repository/postgres"
)

func NewPostgresDB(host, user, password, dbname, port string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	err = db.AutoMigrate(
		&my_postgres.UserSchema{},
		&my_postgres.ProjectSchema{},
		&my_postgres.ColumnSchema{},
		&my_postgres.ProjectMemberSchema{},
		&my_postgres.TaskSchema{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	fmt.Println("Database connection and migration successful!")
	return db
}
