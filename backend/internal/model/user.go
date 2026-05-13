package model

import (
	"adcms-backend/internal/pkg/ids"
	"strings"
	"time"
)

type User struct {
	ID                    string     `gorm:"type:char(36);primaryKey" json:"id"`
	Username              string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password              string     `gorm:"type:varchar(255);not null" json:"-"`
	IsAdmin               bool       `gorm:"column:is_admin;not null;default:false" json:"isAdmin"`
	RealName              string     `gorm:"type:varchar(100)" json:"realName"`
	NickName              string     `gorm:"type:varchar(100)" json:"nickName"`
	Email                 string     `gorm:"type:varchar(100)" json:"email"`
	Phone                 string     `gorm:"type:varchar(20)" json:"phone"`
	SecurityQuestion1     string     `gorm:"type:varchar(255)" json:"-"`
	SecurityAnswer1       string     `gorm:"type:varchar(255)" json:"-"`
	SecurityQuestion2     string     `gorm:"type:varchar(255)" json:"-"`
	SecurityAnswer2       string     `gorm:"type:varchar(255)" json:"-"`
	GoogleAuthSecret      string     `gorm:"type:varchar(255)" json:"-"`
	GoogleAuthBound       bool       `gorm:"not null;default:false" json:"-"`
	NotifyAccountPassword bool       `gorm:"not null;default:true" json:"-"`
	NotifySystemMessage   bool       `gorm:"not null;default:true" json:"-"`
	NotifyTodoTask        bool       `gorm:"not null;default:true" json:"-"`
	Role                  string     `gorm:"type:varchar(50);not null;index" json:"role"`
	ParentID              *string    `gorm:"type:char(36);index" json:"parentId,omitempty"`
	DataScope             string     `gorm:"type:varchar(30);not null;default:'SELF'" json:"dataScope"`
	Status                int8       `gorm:"default:1" json:"status"`
	LastLoginDate         *time.Time `json:"lastLoginDate,omitempty"`
	ExpireDate            *time.Time `json:"expireDate,omitempty"`
	Path                  string     `gorm:"type:varchar(512);not null;default:'';index" json:"path,omitempty"`

	CreatorOptional
	AuditModel
}

func (User) TableName() string {
	return "users"
}

// TenantID 从 path 提取顶层租户 ID（path 第一段），super_admin 返回 DefaultTenantUUID。
// path 格式：'/' (super_admin) | '/adminId/' | '/adminId/agentId/userId/'
func (u User) TenantID() string {
	if u.IsAdmin || u.Role == "super_admin" {
		return ids.DefaultTenantUUID
	}
	if u.Path == "" {
		return u.EffectiveTenantID()
	}
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) > 0 && parts[0] != "" {
		return parts[0]
	}
	return ids.DefaultTenantUUID
}

// EffectiveTenantID 按用户自身关系推导租户ID。优先使用 path，回退到 parent_id。
func (u User) EffectiveTenantID() string {
	if u.IsAdmin || u.Role == "super_admin" {
		return ids.DefaultTenantUUID
	}
	if u.Path != "" {
		parts := strings.Split(strings.Trim(u.Path, "/"), "/")
		if len(parts) > 0 && parts[0] != "" {
			return parts[0]
		}
	}
	if u.Role == "admin" {
		return u.ID
	}
	if u.ParentID != nil && *u.ParentID != "" {
		return *u.ParentID
	}
	return u.ID
}
