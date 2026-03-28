package model

type Permission struct {
	ID          string `gorm:"type:char(36);primaryKey" json:"id"`
	Code        string `gorm:"type:varchar(120);uniqueIndex;not null" json:"code"`
	Name        string `gorm:"type:varchar(120);not null" json:"name"`
	Module      string `gorm:"type:varchar(80);index" json:"module"`
	IsDelegable bool   `gorm:"not null;default:true" json:"isDelegable"`
	AuditModel
}

func (Permission) TableName() string {
	return "permissions"
}
