package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "ADHD-Skills-Tracker/models"
    "github.com/google/uuid"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
   var user models.User
   json.NewDecoder(r.Body).Decode(&user) 
   user.ID = uuid.New()
   user.CreatedAt = time.Now()
   user.UpdatedAt = time.Now()

   err := models.CreateUser(&user)
   if err != nil {
        http.Error(w, err.Error(),http.StatusInternalServerError)
        return
   }

   json.NewEncoder(w).Encode(user)
}
