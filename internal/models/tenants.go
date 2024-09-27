package models

import (
	"time"
)

type UserTenantMapping struct {
	ID       uint64 `gorm:"primaryKey"`
	UserId   uint64 `gorm:"uniqueIndex;not null"`
	TenantId uint64 `gorm:"uniqueIndex;not null"`
}

type Tenant struct {
	ID            uint64 `gorm:"primaryKey"`
	CompanyGuid   string `gorm:"size:50;uniqueIndex;not null"`
	CompanyName   string `gorm:"size:250;uniqueIndex"`
	Host          string `gorm:"size:250;uniqueIndex"`
	BmrmPort      uint32
	SgBizPort     uint32
	TallySyncPort uint32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
