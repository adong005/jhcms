package model

type SiteGroup struct {
	ID          string  `gorm:"type:char(36);primaryKey" json:"id"`
	TenantScoped
	Keyword     string  `gorm:"type:varchar(100)" json:"keyword"`
	Subdomain   string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"subdomain"`
	Title       string  `gorm:"type:varchar(255)" json:"title"`
	Keywords    string  `gorm:"type:varchar(500)" json:"keywords"`
	Description string  `gorm:"type:text" json:"description"`
	ParentID    *string `gorm:"type:char(36);index" json:"parentId,omitempty"`
	CreatorOptional
	AuditModel
}

func (SiteGroup) TableName() string {
	return "site_groups"
}
