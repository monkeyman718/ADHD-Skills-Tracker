package models   

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "log"
    "time"

    "ADHD-Skills-Tracker/database"
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
)

// User represents a user in the database
type User struct {
    ID          uuid.UUID   `json:"id"  db:"id"`
    Email       string      `json:"email" db:"email"`
    Password    string      `json:"password" db:"password_hash"`
    CreatedAt   time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}

func AuthenticateUser(email, password string) (*User, error) {
    var user User

    queryStr := "SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1"
   row, err := database.DB.Query(context.Background(), queryStr, email)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, errors.New("Invalid email or password")
        }
        log.Printf("Error querying user: %v", err)
        return nil, fmt.Errorf("error querying user: %v", err)
    }

    for row.Next() {
        err := row.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, fmt.Errorf("unable to scan row: %w", err)
        }
    }

    // Log the password hash in the database and the provided password (for debugging)
    log.Printf("Stored password hash: %v", user.Password)
    log.Printf("Provided password: %v", password)

    // Compare the provided password with the hashed password stored in the database
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        log.Printf("Password comparison failed: %v", err)
        return nil, errors.New("Invalid email or password")
    }

    return &user, nil
}


func CreateUser(user *User) error {
    // Hash the password before storing it
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        return err
    }
    user.Password = string(hashedPassword)  // Store the hashed password

    queryStr := "INSERT INTO users (id, email, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5);"

    _, err = database.DB.Exec(context.Background(), queryStr, user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
    if err != nil {
      log.Printf("Error inserting user into DB: %v", err)
      return fmt.Errorf("error inserting new user: %w", err)
   }


    return nil
}

