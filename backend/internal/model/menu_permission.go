package model

// MenuPermission 菜单与权限点多对多（导航/授权 UI 用）。
type MenuPermission struct {
	ID           string `gorm:"type:char(36);primaryKey" json:"id"`
	MenuID       string `gorm:"type:char(36);column:menu_id;not null;uniqueIndex:idx_menu_permission_pair" json:"menuId"`
	PermissionID string `gorm:"type:char(36);column:permission_id;not null;uniqueIndex:idx_menu_permission_pair" json:"permissionId"`
}

func (MenuPermission) TableName() string {
	return "menu_permissions"
}
