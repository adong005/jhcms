package model

type InfoCategory struct {
	ID          string `gorm:"type:char(36);primaryKey" json:"id"`
	TenantScoped
	Name        string `gorm:"type:varchar(120);not null;index" json:"name"`
	Code        string `gorm:"type:varchar(80);uniqueIndex;not null" json:"code"`
	IsHome      int8   `gorm:"default:1;index" json:"isHome"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Description string `gorm:"type:text" json:"description"`
	Status      int8   `gorm:"default:1;index" json:"status"`
	CreatorOptional
	AuditModel
}

func (InfoCategory) TableName() string {
	return "info_categories"
}
