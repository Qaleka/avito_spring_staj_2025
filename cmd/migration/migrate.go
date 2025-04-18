package main

import (
	"avito_spring_staj_2025/internal/db"
	"database/sql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"log"
	"path/filepath"
)

func main() {
	_ = godotenv.Load()
	databaseURL := db.FromEnv()
	if databaseURL == "" {
		log.Fatal("DSN not provided. Please set DB_HOST and other DB_* env variables.")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	migrationsDir := filepath.Join("cmd", "migration", "migrations")

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
