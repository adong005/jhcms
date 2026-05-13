package model

import "time"

type SystemLog struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	RequestID   string    `gorm:"type:char(36);not null;default:''" json:"requestId"`
	TenantID    string    `gorm:"type:char(36);index;not null" json:"tenantId"`
	UserID      string    `gorm:"type:char(36);not null;default:'';index:idx_logs_user_id" json:"userId"`
	Path        string    `gorm:"type:varchar(512);not null;default:'';index:idx_logs_path" json:"path"`
	Username    string    `gorm:"type:varchar(50);index" json:"username"`
	Action      string    `gorm:"type:varchar(100);index:idx_logs_action_st" json:"action"`
	Module      string    `gorm:"type:varchar(100)" json:"module"`
	Description string    `gorm:"type:text" json:"description"`
	TargetID    string    `gorm:"type:varchar(64);not null;default:''" json:"targetId"`
	IP          string    `gorm:"type:varchar(50)" json:"ip"`
	Method      string    `gorm:"type:varchar(10);not null;default:''" json:"method"`
	URL         string    `gorm:"type:varchar(500);not null;default:''" json:"url"`
	UserAgent   string    `gorm:"type:varchar(255);not null;default:''" json:"userAgent"`
	Status      string    `gorm:"type:varchar(20);index:idx_logs_action_st" json:"status"`
	LogType     string    `gorm:"type:varchar(20);not null;default:'api'" json:"logType"`
	StatusCode  int       `gorm:"type:smallint;not null;default:0" json:"statusCode"`
	Duration    int       `json:"duration"`
	ErrorMsg    string    `gorm:"type:text" json:"errorMsg"`
	RequestJSON string    `gorm:"type:longtext" json:"requestJson"`
	ParentID    *string   `gorm:"type:char(36);index" json:"parentId,omitempty"`
	CreatedAt   time.Time `gorm:"index:idx_logs_tenant_time" json:"createTime"`
}

func (SystemLog) TableName() string {
	return "system_logs"
}
