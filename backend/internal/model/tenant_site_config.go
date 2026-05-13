package model

type TenantSiteConfig struct {
	ID             string `gorm:"type:char(36);primaryKey" json:"id"`
	TenantID       string `gorm:"type:char(36);uniqueIndex;not null" json:"tenantId"`
	Title          string `gorm:"type:varchar(255);default:''" json:"title"`
	Keywords       string `gorm:"type:varchar(500);default:''" json:"keywords"`
	Description    string `gorm:"type:text" json:"description"`
	Domain         string `gorm:"type:varchar(255);default:'';index" json:"domain"`
	Logo           string `gorm:"type:varchar(255);default:''" json:"logo"`
	ICPCode        string `gorm:"type:varchar(120);default:''" json:"icpCode"`
	ContactPhone   string `gorm:"type:varchar(64);default:''" json:"contactPhone"`
	ContactAddress string `gorm:"type:varchar(255);default:''" json:"contactAddress"`
	ContactEmail   string `gorm:"type:varchar(120);default:''" json:"contactEmail"`
	AuditModel
}

func (TenantSiteConfig) TableName() string {
	return "tenant_site_configs"
}
