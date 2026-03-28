package model

type Info struct {
	ID           string  `gorm:"type:char(36);primaryKey" json:"id"`
	TenantScoped
	Title        string  `gorm:"type:varchar(255);not null" json:"title"`
	CategoryID   *string `gorm:"type:char(36);index" json:"categoryId,omitempty"`
	CategoryName string  `gorm:"-" json:"categoryName,omitempty"`
	Author       string  `gorm:"type:varchar(100)" json:"author"`
	Summary      string  `gorm:"type:text" json:"summary"`
	Content      string  `gorm:"type:longtext" json:"content"`
	Status       int8    `gorm:"default:0;index" json:"status"`
	ParentID     *string `gorm:"type:char(36);index" json:"parentId,omitempty"`
	CreatorOptional
	AuditModel
}

func (Info) TableName() string {
	return "infos"
}
