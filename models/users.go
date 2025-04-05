package models   

import (
    "context"
    "errors"
    "fmt"
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
    fmt.Printf("created_hashed_password: %v\n", user.Password)

    queryStr := "INSERT INTO users (id, email, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5);"
    _, err = database.DB.Exec(context.Background(), queryStr, user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt) 
    return err 
}

func AuthenticateUser(email, password string) (*User, error) {
    var user User

    queryStr := "SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1"

    err := database.DB.QueryRow(context.Background(), queryStr, email).Scan(&user.ID,&user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
    if err != nil {
        return nil, errors.New("Invalid email or password")
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return nil, errors.New("Provided password is not correct")
    }

    return &user, nil
}
