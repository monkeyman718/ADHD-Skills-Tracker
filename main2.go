package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/jackc/pgx/v5"
    "github.com/golang-jwt/jwt/v5"
    _ "github.com/joho/godotenv/autoload"
)

func main(){
   conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
   if err != nil {
        log.Printf("Cannot connect to database: %v", err)
   }

   log.Printf("Connected to database!...")

   defer conn.Close(context.Background())

   var queryResponse string

   err = conn.QueryRow(context.Background(), "SELECT version()").Scan(&queryResponse)
   if err != nil {
        log.Printf("Query did not return any results: %v", err)
   }

   fmt.Printf("Query response: %s\n", queryResponse)
}
