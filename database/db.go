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

/*func ConnectDB() {
    databaseUrl := os.Getenv("DATABASE_URL")
    if databaseUrl == "" {
        log.Fatal("DATABASE_URL environment variable not found.")
    }

    config, err := pgxpool.ParseConfig(databaseUrl)
    if err != nil {
        log.Fatalf("Unable to parse database URL: %v", err)
    }

    // 🔥 This disables the statement cache to avoid the Supabase error
    config.ConnConfig.PreferSimpleProtocol = true

    DB, err = pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }

    fmt.Println("Connected to the database!")
}*/


func CloseDB() {
    if DB != nil {
        DB.Close()
        fmt.Println("Database connection closed.")
    }
}
