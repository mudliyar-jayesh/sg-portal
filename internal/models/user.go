package models

import (
	"time"
)

// User represents a user entity in the system.
type User struct {
    ID           uint64      `gorm:"primaryKey;autoIncrement" json:"id"`       // Auto-incrementing primary key
    Email        string      `gorm:"not null;unique" json:"email"`             // Required and unique
    Password     string      `gorm:"-" json:"password"`                        // Exclude from the database
    Name         string      `gorm:"not null" json:"name"`                     // Required, non-null
    CountryID    *int        `gorm:"default:0" json:"country_id"`              // Optional, defaults to 0
    MobileNumber string      `gorm:"not null" json:"mobile_number"`            // Required field
    IsActive     bool        `gorm:"default:true" json:"is_active"`            // Defaults to true
    LastLogin    *time.Time  `json:"last_login_time"`                          // Nullable datetime
    CreatedAt    time.Time   `json:"created_at"`                               // GORM will automatically handle this
    UpdatedAt    time.Time   `json:"updated_at"`                               // GORM will automatically handle this
}

