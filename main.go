package main

import (
    "context"
    "log"
    "net/http"

    "ADHD-Skills-Tracker/database"
    "github.com/gorilla/mux"
)

func main(){
    database.ConnectDB()
    defer database.CloseDB()

    // Example query to test connection
    var version string
    err := database.DB.QueryRow(context.Background(), "SELECT version()").Scan(&version)
    if err != nil {
        log.Fatalf("Query failed: %v", err)
    }

    log.Println("Connected to PostgreSQL:", version)

    // Set up router
    r := mux.NewRouter()
    // Add routes here...
    log.Println("Server running on port 8080...")
    http.ListenAndServe(":8080", r)
}
