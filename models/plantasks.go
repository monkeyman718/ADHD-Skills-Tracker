package models

import (
    "time"
    "github.com/google/uuid"
)

type PlanTask struct {
    ID              uuid.UUID   `json:"id" db:"id"`
    PlanID          uuid.UUID   `json:"plan_id" db:"plan_id"`
    TaskID          uuid.UUID   `json:"task_id" db:"task_id"`
    Status          string      `json:"status" db:"status"`
    CreatedAt       time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time   `json:"updated_at" db:"updated_at"`
}
