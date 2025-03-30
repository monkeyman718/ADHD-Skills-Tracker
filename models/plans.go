package models

import (
    "time"
    "github.com/google/uuid"
)

type Plan struct {
    ID              uuid.UUID   `json:"id" db:"id"`
    UserID          uuid.UUID   `json:"user_id" db:"user_id"`
    Date            time.Time   `json:"date" db:"date"`
    Type            string      `json:"type" db:"type"`
    CreatedAt       time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}
