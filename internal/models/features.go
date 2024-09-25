package models

import (
	"time"
)

type Feature struct {
	ID         uint32 `gorm:"primaryKey"`
	Name       string `gorm:"size:200;not null"`
	Permission string `gorm:"size:50;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserFeatureMapping struct {
	ID        uint64 `gorm:"primaryKey"`
	UserId    uint64 `gorm:"unique;not null"`
	FeatureId uint32 `gorm:"unique;not null"`
}
