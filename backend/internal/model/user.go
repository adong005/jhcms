package model

import (
	"adcms-backend/internal/pkg/ids"
	"time"
)

type User struct {
	ID                   string     `gorm:"type:char(36);primaryKey" json:"id"`
	Username             string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password             string     `gorm:"type:varchar(255);not null" json:"-"`
	IsAdmin              bool       `gorm:"column:is_admin;not null;default:false" json:"isAdmin"`
	RealName             string     `gorm:"type:varchar(100)" json:"realName"`
	NickName             string     `gorm:"type:varchar(100)" json:"nickName"`
	Email                string     `gorm:"type:varchar(100)" json:"email"`
	Phone                string     `gorm:"type:varchar(20)" json:"phone"`
	SecurityQuestion1    string     `gorm:"type:varchar(255)" json:"-"`
	SecurityAnswer1      string     `gorm:"type:varchar(255)" json:"-"`
	SecurityQuestion2    string     `gorm:"type:varchar(255)" json:"-"`
	SecurityAnswer2      string     `gorm:"type:varchar(255)" json:"-"`
	GoogleAuthSecret     string     `gorm:"type:varchar(255)" json:"-"`
	GoogleAuthBound      bool       `gorm:"not null;default:false" json:"-"`
	NotifyAccountPassword bool      `gorm:"not null;default:true" json:"-"`
	NotifySystemMessage   bool      `gorm:"not null;default:true" json:"-"`
	NotifyTodoTask        bool      `gorm:"not null;default:true" json:"-"`
	Role                 string     `gorm:"type:varchar(50);not null;index" json:"role"`
	ParentID             *string    `gorm:"type:char(36);index" json:"parentId,omitempty"`
	DataScope            string     `gorm:"type:varchar(30);not null;default:'SELF'" json:"dataScope"`
	Status               int8       `gorm:"default:1" json:"status"`
	LastLoginDate        *time.Time `json:"lastLoginDate,omitempty"`
	ExpireDate           *time.Time `json:"expireDate,omitempty"`

	Title       string `gorm:"type:varchar(255)" json:"title,omitempty"`
	Keywords    string `gorm:"type:varchar(500)" json:"keywords,omitempty"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Domain      string `gorm:"type:varchar(255)" json:"domain,omitempty"`
	Logo        string `gorm:"type:varchar(255)" json:"logo,omitempty"`
	ICPCode     string `gorm:"type:varchar(120)" json:"icpCode,omitempty"`
	ContactPhone string `gorm:"type:varchar(64)" json:"contactPhone,omitempty"`
	ContactAddress string `gorm:"type:varchar(255)" json:"contactAddress,omitempty"`
	ContactEmail string `gorm:"type:varchar(120)" json:"contactEmail,omitempty"`

	CreatorOptional
	AuditModel
}

func (User) TableName() string {
	return "users"
}

// EffectiveTenantID 按用户自身关系推导租户ID，不再依赖 users.tenant_id。
func (u User) EffectiveTenantID() string {
	if u.IsAdmin || u.Role == "super_admin" {
		return ids.DefaultTenantUUID
	}
	if u.Role == "admin" {
		return u.ID
	}
	if u.ParentID != nil && *u.ParentID != "" {
		return *u.ParentID
	}
	return u.ID
}
