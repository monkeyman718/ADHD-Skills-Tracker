package models   

import (
    "context"
    "time"

    "ADHD-Skills-Tracker/database"
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
)

// User represents a user in the database
type User struct {
    ID          uuid.UUID   `json:"id"  db:"id"`
    Email       string      `json:"email" db:"email"`
    Password    string      `json:"-" db:"password_hash"`
    CreatedAt   time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

func CreateUser(user *User) error {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }
    user.Password = string(hashedPassword)

    queryStr := "INSERT INTO users (id, email, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5);"
    _, err = database.DB.Exec(context.Background(), queryStr, user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt) 
    return err 
}
