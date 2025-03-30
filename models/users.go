package users

import (
    "time"
    "github.com/google/uuid"
)

// User represents a user in the database
type User struct {
    ID          uuid.UUID   `json:"id"  db:"id"`
    Email       string      `json:"email" db:"email"`
    Password    string      `json:"-" db:"password"`
    CreatedAt   time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}
