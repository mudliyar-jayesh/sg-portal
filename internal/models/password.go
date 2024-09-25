package models

import (
	"time"
)

// UserPassword represents the password entity for a user.
type UserPassword struct {
    ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`       // Auto-incrementing primary key
    UserID     uint64    `gorm:"not null" json:"user_id"`                  // Foreign key for User, required
    Password   string    `gorm:"not null" json:"password"`                 // User's password, required
    Salt       string    `gorm:"not null" json:"salt"`                     // Password salt, required
    CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`         // Automatically set when the record is first created
    UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`         // Automatically updated when the record is modified
}

// NewUserPassword is a helper function to initialize a new UserPassword with the current date.
func NewUserPassword(userID uint64, password, salt string) *UserPassword {
    return &UserPassword{
        UserID:   userID,
        Password: password,
        Salt:     salt,
    }
}
