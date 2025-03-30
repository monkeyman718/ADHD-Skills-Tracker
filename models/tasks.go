package models

import (
    "time"
    "github.com/google/uuid"
)

type Task struct {
    ID              uuid.UUID   `json:"id" db:"id"`
    SkillID         uuid.UUID   `json:"skill_id" db:"skill_id"`
    Name            string      `json:"name" db:"name"`
    Status          string      `json:"status" db:"status"`
    DueDate         time.Time   `json:"due_date" db:"due_date"`
    CreatedAt       time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}
