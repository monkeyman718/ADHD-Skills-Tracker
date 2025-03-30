package database

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/jackc/pgx/v5/pgxpool"
    _ "github.com/joho/godotenv/autoload"
)

var DB *pgxpool.Pool

func ConnectDB() {
    var err error

    databaseUrl := os.Getenv("DATABASE_URL")
    if databaseUrl == "" {
        log.Fatal("DATABASE_URL environment variable not found.")
    }

    DB, err = pgxpool.New(context.Background(), databaseUrl)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }
    
    fmt.Println("Connected to the database!")
}

func CloseDB() {
    if DB != nil {
        DB.Close()
        fmt.Println("Database connection closed.")
    }
}
