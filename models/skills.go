package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "ADHD-Skills-Tracker/database"
)

// Skill struct for the "skills" table
type Skill struct {
    ID        uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    UserID    uuid.UUID `json:"user_id"`
    Name      string    `json:"name"`
    Priority  string    `json:"priority"`
    Goal      string    `json:"goal"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

// CreateSkill creates a new skill in the database
func CreateSkill(skill *Skill) error {
    result := database.DB.Create(skill)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

// GetSkill retrieves a skill by its ID
func GetSkill(id uuid.UUID) (Skill, error) {
    var skill Skill
    result := database.DB.First(&skill, "id = ?", id)
    if result.Error != nil {
        return skill, result.Error
    }
    return skill, nil
}

