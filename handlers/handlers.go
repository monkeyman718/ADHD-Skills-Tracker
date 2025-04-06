package handlers

import (
    "encoding/json"
    "fmt"
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

   fmt.Printf("id: %v, email: %v, passwd: %v, created: %v, updated: %v\n",user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)

   err := models.CreateUser(&user)
   if err != nil {
        http.Error(w, err.Error(),http.StatusInternalServerError)
        return
   }

   json.NewEncoder(w).Encode(user)
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
    var creds struct {
        Email       string `json:"email"`
        Password    string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
        http.Error(w, "Invalid login payload", http.StatusBadRequest)
        return
    }

    // authenticate user with models/AuthenticateUser
    // check if user authentication was valid
    user, err := models.AuthenticateUser(creds.Email, creds.Password)
    if err != nil {
        http.Error(w, "Login failed: " + err.Error(), http.StatusUnauthorized)
       // fmt.Fprintf(w, "password_hash: %v, entered_password: %v\n", user.Password, creds.Password)
        return
    }

    // for now return user ID on successful login
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Login successful", 
        "user_id": user.ID.String(),
    })
}
