package database

import (
    "database/sql"
    "log"
    "os"
	"github.com/joho/godotenv"
    _ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func ConnectDB() {
	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

	dsn := os.Getenv("DATABASE_URL")
    
	var err error
    DB, err = sql.Open("pgx", dsn) // pgx driver
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }

    if err := DB.Ping(); err != nil {
        log.Fatalf("Database not reachable: %v", err)
    }

    log.Println("Connected to PostgreSQL successfully!")
}

func CloseDB() {
    if DB != nil {
        DB.Close()
        log.Println("Database connection closed.")
    }
}
