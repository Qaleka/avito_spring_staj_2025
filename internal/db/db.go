package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func DbConnect() *sql.DB {
	dsnString := FromEnv()
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

func FromEnv() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		return ""
	}
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, dbname)
}
