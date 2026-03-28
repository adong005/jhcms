package model

type Role struct {
	ID          string  `gorm:"type:char(36);primaryKey" json:"id"`
	TenantScoped
	Name        string  `gorm:"type:varchar(100);not null" json:"name"`
	Code        string  `gorm:"type:varchar(50);not null;index:idx_roles_tenant_code,unique" json:"code"`
	DataScope   string  `gorm:"type:varchar(30);not null;default:'SELF'" json:"dataScope"`
	Description string  `gorm:"type:text" json:"description"`
	ParentID    *string `gorm:"type:char(36);index" json:"parentId,omitempty"`
	Status      int8    `gorm:"default:1" json:"status"`
	CreatorOptional
	AuditModel
}

func (Role) TableName() string {
	return "roles"
}
