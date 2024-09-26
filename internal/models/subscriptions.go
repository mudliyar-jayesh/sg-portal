package models

import (
	"time"
)

type Subscription struct {
	ID        uint32 `gorm:"primaryKey"`
	Name      string `gorm:"size:200;not null"`
	Code      string `gorm:"size:50;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserSubscriptionMapping struct {
	ID             uint64 `gorm:"primaryKey"`
	UserId         uint64 `gorm:"unique;not null"`
	SubscriptionId uint32 `gorm:"unique;not null"`
}

type FeatureSubscriptionMapping struct {
	ID             uint64 `gorm:"primaryKey"`
	FeatureId      uint32 `gorm:"unique;not null"`
	SubscriptionId uint32 `gorm:"unique;not null"`
}

type UserSubscriptionHistory struct {
	ID               uint64 `gorm:"primaryKey"`
	UserId           uint64 `gorm:"unique;not null"`
	SubscriptionId   uint32 `gorm:"unique;not null"`
	StartDate        time.Time
	RenewalDate      time.Time
	ExpiryDate       time.Time
	NumberOfRenewals uint16
}
