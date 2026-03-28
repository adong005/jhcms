package model

type Menu struct {
	ID             string  `gorm:"type:char(36);primaryKey" json:"id"`
	TenantScoped
	ParentIDMenu   *string `gorm:"type:char(36);index" json:"parentIdMenu,omitempty"`
	Name           string  `gorm:"type:varchar(100);not null" json:"name"`
	Path           string  `gorm:"type:varchar(255)" json:"path"`
	Component      string  `gorm:"type:varchar(255)" json:"component"`
	Icon           string  `gorm:"type:varchar(100)" json:"icon"`
	PermissionCode string  `gorm:"type:varchar(120);index" json:"permissionCode"`
	IsDelegable    bool    `gorm:"not null;default:true" json:"isDelegable"`
	SortOrder      int     `gorm:"default:0" json:"sortOrder"`
	IsShow         int8    `gorm:"column:is_show;default:1;index" json:"isShow"`
	Status         int8    `gorm:"default:1" json:"status"`
	AuditModel
}

func (Menu) TableName() string {
	return "menus"
}
