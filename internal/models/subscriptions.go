package models

import (
	"time"
)

// Subscription model, no changes needed here
type Subscription struct {
	ID        uint32 `gorm:"primaryKey"`
	Name      string `gorm:"size:200;not null"`
	Code      string `gorm:"size:50;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserSubscriptionMapping - combination of UserId and SubscriptionId must be unique
type UserSubscriptionMapping struct {
	ID             uint64 `gorm:"primaryKey"`
	UserId         uint64 `gorm:"not null;uniqueIndex:idx_user_subscription"` // Part of composite unique index
	SubscriptionId uint32 `gorm:"not null;uniqueIndex:idx_user_subscription"` // Part of composite unique index
}

// FeatureSubscriptionMapping - combination of FeatureId and SubscriptionId must be unique
type FeatureSubscriptionMapping struct {
	ID             uint64 `gorm:"primaryKey"`
	FeatureId      uint32 `gorm:"not null;uniqueIndex:idx_feature_subscription"` // Part of composite unique index
	SubscriptionId uint32 `gorm:"not null;uniqueIndex:idx_feature_subscription"` // Part of composite unique index
}

// UserSubscriptionHistory - combination of UserId and SubscriptionId must be unique
type UserSubscriptionHistory struct {
	ID               uint64 `gorm:"primaryKey"`
	UserId           uint64 `gorm:"not null;uniqueIndex:idx_user_subscription_history"` // Part of composite unique index
	SubscriptionId   uint32 `gorm:"not null;uniqueIndex:idx_user_subscription_history"` // Part of composite unique index
	StartDate        time.Time
	RenewalDate      time.Time
	ExpiryDate       time.Time
	NumberOfRenewals uint16
}

