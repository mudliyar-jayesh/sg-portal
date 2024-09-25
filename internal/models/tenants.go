package models

import (
	"time"
)

type UserTenantMapping struct {
	ID       uint64 `gorm:"primaryKey"`
	UserId   uint64 `gorm:"unique;not null"`
	TenantId uint64 `gorm:"unique;not null"`
}

type Tenant struct {
	ID            uint64 `gorm:"primaryKey"`
	CompanyGuid   string `gorm:"size:50;unique;not null"`
	CompanyName   string `gorm:"size:250"`
	Host          string `gorm:"size:250"`
	BmrmPort      uint32
	SgBizPort     uint32
	TallySyncPort uint32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
