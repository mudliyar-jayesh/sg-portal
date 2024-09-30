package models

import (
	"time"
)

type UserTenantMapping struct {
	ID       uint64 `gorm:"primaryKey"`
	UserId   uint64 `gorm:"uniqueIndex:idx_tnt_mapping;not null"`
	TenantId uint64 `gorm:"uniqueIndex:idx_tnt_mapping;not null"`
}

type Tenant struct {
	ID            uint64 `gorm:"primaryKey"`
	CompanyGuid   string `gorm:"size:50;uniqueIndex:idx_tnt:;not null"`
	CompanyName   string `gorm:"size:250;uniqueIndex:idx_tnt"`
	Host          string `gorm:"size:250;uniqueIndex:idx_tnt"`
	BmrmPort      uint32
	SgBizPort     uint32
	TallySyncPort uint32
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
