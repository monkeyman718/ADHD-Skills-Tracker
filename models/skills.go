package models

import (
    "time"
    "github.com/google/uuid"
)

type Skill struct {
    ID          uuid.UUID   `json:"id" db:"id"`
    UserID      uuid.UUID   `json:"user_id" db:"user_id"`
    Name        string      `json:"name" db:"name"`
    Priority    int         `json:"priority db:priority"`
    Goal        string      `json:"goal" db:"goal"`
    Status      string      `json:"status" db:"status"`
    CreatedAt   time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt   time.Time   `json:"updated_at" db:"updated_at"`
}
