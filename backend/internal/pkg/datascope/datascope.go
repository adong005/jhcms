package datascope

import (
	"gorm.io/gorm"
)

// ApplyUserPathSubtree 在 column 上应用用户 path 前缀子树过滤（常用于管理员数据范围）。
func ApplyUserPathSubtree(db *gorm.DB, column, userPath string) *gorm.DB {
	if userPath != "" && userPath != "/" {
		return db.Where(column+" LIKE ?", userPath+"%")
	}
	return db
}
