package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

func main() {
	// Load .env if present
	_ = godotenv.Load()

	dbURL := os.Getenv("POSTGRES_DSN")
	if dbURL == "" {
		log.Fatal("POSTGRES_DSN environment variable is required")
	}

	migrationsDir := flag.String("dir", "internal/infrastructure/db/migrations", "migrations directory")
	flag.Parse()

	db, err := goose.OpenDBWithDriver("pgx", dbURL)
	if err != nil {
		log.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	if err := goose.Up(db, *migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("Migrations applied successfully.")
}
