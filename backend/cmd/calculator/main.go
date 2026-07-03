package main

import (
	"Calculator/internal/api"
	"Calculator/internal/service"
	"Calculator/internal/storage"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Cant load .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Db config error: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Cant connect to database: %v", err)
	}

	postgresStore := storage.NewPostgresStorage(db)
	calculator := service.NewCalculator(postgresStore)
	server := api.NewEchoServer(calculator)

	err = server.StartListening()
	if err != nil {
		panic(err)
	}
}
