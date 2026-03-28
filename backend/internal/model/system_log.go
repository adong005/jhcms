package model

import "time"

type SystemLog struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	TenantID    string    `gorm:"type:char(36);index;not null" json:"tenantId"`
	Username    string    `gorm:"type:varchar(50);index" json:"username"`
	Action      string    `gorm:"type:varchar(100);index" json:"action"`
	Module      string    `gorm:"type:varchar(100)" json:"module"`
	Description string    `gorm:"type:text" json:"description"`
	IP          string    `gorm:"type:varchar(50)" json:"ip"`
	Status      string    `gorm:"type:varchar(20)" json:"status"`
	Duration    int       `json:"duration"`
	ErrorMsg    string    `gorm:"type:text" json:"errorMsg"`
	RequestJSON string    `gorm:"type:longtext" json:"requestJson"`
	ParentID    *string   `gorm:"type:char(36);index" json:"parentId,omitempty"`
	CreatedAt   time.Time `json:"createTime"`
}

func (SystemLog) TableName() string {
	return "system_logs"
}
