package db

import (
	"avito_spring_staj_2025/internal/service/dsn"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func DbConnect() *sql.DB {
	dsnString := dsn.FromEnv()
	if dsnString == "" {
		log.Fatal("DSN not provided. Please check your environment variables.")
	}

	db, err := sql.Open("postgres", dsnString)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	const maxOpenConns = 1000
	const maxIdleConns = 500
	const connMaxLifetime = 30 * time.Minute
	const connMaxIdleTime = 5 * time.Minute

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetConnMaxIdleTime(connMaxIdleTime)

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("Connected to database")
	return db
}
