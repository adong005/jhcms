package model

import (
	"time"

	"gorm.io/gorm"
)

// AuditModel 创建时间、更改时间、GORM 软删除。
type AuditModel struct {
	CreatedAt time.Time      `json:"createTime"`
	UpdatedAt time.Time      `json:"updateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TenantScoped 多租户 tenant_id（UUID）。
type TenantScoped struct {
	TenantID string `gorm:"type:char(36);index;not null" json:"tenantId"`
}

// CreatorOptional 创建人用户 ID（UUID）。
type CreatorOptional struct {
	CreatedBy *string `gorm:"type:char(36);index" json:"createdBy,omitempty"`
}
