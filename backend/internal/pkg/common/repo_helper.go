package common

import (
	"gorm.io/gorm"
)

// BuildListQuery 构建列表查询（通用方法）
func BuildListQuery(db *gorm.DB, model interface{}, filters map[string]interface{}) *gorm.DB {
	query := db.Model(model)

	// 通用过滤逻辑
	for key, value := range filters {
		if value == nil || value == "" {
			continue
		}

		switch key {
		case "name", "username", "title", "realName", "nickName":
			// 模糊搜索
			query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
		case "status", "role", "type":
			// 精确匹配
			query = query.Where(key+" = ?", value)
		default:
			// 其他字段精确匹配
			query = query.Where(key+" = ?", value)
		}
	}

	return query
}

// ApplyPagination 应用分页
func ApplyPagination(query *gorm.DB, page, pageSize int) *gorm.DB {
	offset := (page - 1) * pageSize
	return query.Offset(offset).Limit(pageSize)
}

// ApplyOrder 应用排序
func ApplyOrder(query *gorm.DB, orderBy string) *gorm.DB {
	if orderBy == "" {
		orderBy = "created_at DESC"
	}
	return query.Order(orderBy)
}

// GetListWithPagination 获取分页列表（通用方法）
func GetListWithPagination(db *gorm.DB, model interface{}, page, pageSize int, filters map[string]interface{}, orderBy string) (interface{}, int64, error) {
	query := BuildListQuery(db, model, filters)

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 应用排序和分页
	query = ApplyOrder(query, orderBy)
	query = ApplyPagination(query, page, pageSize)

	// 查询数据
	var results interface{}
	if err := query.Find(&results).Error; err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

// SoftDelete 软删除
func SoftDelete(db *gorm.DB, model interface{}, id int64) error {
	return db.Delete(model, id).Error
}

// BatchDelete 批量删除
func BatchDelete(db *gorm.DB, model interface{}, ids []int64) error {
	return db.Delete(model, ids).Error
}

// UpdateFields 更新指定字段
func UpdateFields(db *gorm.DB, model interface{}, id int64, fields map[string]interface{}) error {
	return db.Model(model).Where("id = ?", id).Updates(fields).Error
}

// ExistsByID 检查记录是否存在
func ExistsByID(db *gorm.DB, model interface{}, id int64) bool {
	var count int64
	db.Model(model).Where("id = ?", id).Count(&count)
	return count > 0
}

// ExistsByField 检查字段值是否存在
func ExistsByField(db *gorm.DB, model interface{}, field string, value interface{}) bool {
	var count int64
	db.Model(model).Where(field+" = ?", value).Count(&count)
	return count > 0
}
