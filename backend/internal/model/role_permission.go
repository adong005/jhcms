package model

type RolePermission struct {
	ID           string `gorm:"type:char(36);primaryKey" json:"id"`
	RoleID       string `gorm:"type:char(36);index:idx_role_permission_unique,unique;not null" json:"roleId"`
	PermissionID string `gorm:"type:char(36);index:idx_role_permission_unique,unique;not null" json:"permissionId"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
