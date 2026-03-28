package model

type Form struct {
	ID           string `gorm:"type:char(36);primaryKey" json:"id"`
	TenantScoped
	Contact      string `gorm:"type:varchar(100)" json:"contact"`
	Phone        string `gorm:"type:varchar(20)" json:"phone"`
	Company      string `gorm:"type:varchar(255)" json:"company"`
	IP           string `gorm:"type:varchar(50)" json:"ip"`
	HandleStatus int8   `gorm:"default:0;index" json:"handleStatus"`
	Remark       string `gorm:"type:text" json:"remark"`
	CreatorOptional
	AuditModel
}

func (Form) TableName() string {
	return "forms"
}
