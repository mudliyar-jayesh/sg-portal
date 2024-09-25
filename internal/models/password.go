package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserPassword represents the password entity for a user.
type UserPassword struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"` // Auto-incrementing primary key
	UserID     uint64    `gorm:"not null" json:"user_id"`            // Foreign key for User, required
	Password   string    `gorm:"not null" json:"password"`           // User's hashed password, required
	Salt       string    `gorm:"not null" json:"salt"`               // Password salt, required
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`   // Automatically set when the record is first created
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`   // Automatically updated when the record is modified
}

// NewUserPassword is a helper function to initialize a new UserPassword with the current date.
func NewUserPassword(userID uint64, password, salt string) *UserPassword {
	return &UserPassword{
		UserID:   userID,
		Password: password,
		Salt:     salt,
	}
}

// HashPassword hashes a password with the given salt using bcrypt.
func HashPassword(password, salt string) (string, error) {
	combined := password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(combined), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ValidatePassword compares the provided password and salt against the stored hashed password.
func ValidatePassword(providedPassword, salt, storedHashedPassword string) error {
	combined := providedPassword + salt
	return bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(combined))
}

// GenerateSalt creates a cryptographically secure random salt.
func GenerateSalt() (string, error) {
	// Define the salt size in bytes (e.g., 16 bytes = 128 bits)
	saltSize := 16
	salt := make([]byte, saltSize)

	// Generate a random salt using crypto/rand
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Convert the byte slice to a hex string
	return hex.EncodeToString(salt), nil
}
