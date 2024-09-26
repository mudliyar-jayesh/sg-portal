package models

import (
	"time"

	"github.com/google/uuid"
)

// Token represents the token entity for user authentication.
type Token struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"` // Auto-incrementing primary key
	UserID    uint64    `gorm:"not null" json:"user_id"`            // Foreign key for User, required
	Value     uuid.UUID `gorm:"type:uuid;not null" json:"value"`    // UUID for the token value
	Expiry    time.Time `gorm:"not null" json:"expiry"`             // Token expiration time, required
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`   // Automatically set when the record is first created
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`   // Automatically updated when the record is modified
}

// NewToken is a helper function to initialize a new Token with the current time.
func NewToken(userID uint64, expiry time.Time) *Token {
	return &Token{
		UserID: userID,
		Value:  uuid.New(), // Generate a new UUID for the token
		Expiry: expiry,     // Set the provided expiration time
	}
}
