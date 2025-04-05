package main

import (
    "log"
    "net/http"

    "ADHD-Skills-Tracker/database"
    "ADHD-Skills-Tracker/routes"
)

func main(){
    database.ConnectDB()
    defer database.CloseDB()

    r := routes.CreateRoutes() 

    log.Println("Server running on port 8080...")
    http.ListenAndServe(":8080", r)
}
