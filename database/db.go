package database

import (
    "fmt"
    "log"
    "os"

    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    _ "github.com/joho/godotenv/autoload"
)

var DB *gorm.DB

// ConnectDB connects to the database using GORM
func ConnectDB() {
    var err error

    dsn := os.Getenv("DATABASE_URL") // Database connection string

    // Open connection to the database
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Unable to connect to database: %v", err)
    }

    fmt.Println("Connected to the database!")
}

// CloseDB closes the database connection
func CloseDB() {
    db, err := DB.DB()
    if err != nil {
        log.Fatalf("Unable to get DB instance: %v", err)
    }
    db.Close()
    fmt.Println("Database connection closed.")
}

