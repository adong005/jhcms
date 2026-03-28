package repository

import (
	"gorm.io/gorm"
)

// BaseRepository 基础仓库结构
type BaseRepository struct {
	db        *gorm.DB
	modelType interface{}
}

// NewBaseRepository 创建基础仓库
func NewBaseRepository(db *gorm.DB, modelType interface{}) *BaseRepository {
	return &BaseRepository{
		db:        db,
		modelType: modelType,
	}
}

// GetDB 获取数据库连接
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}

// BuildQuery 构建查询（可被子类覆盖）
func (r *BaseRepository) BuildQuery(filters map[string]interface{}) *gorm.DB {
	query := r.db.Model(r.modelType)

	// 通用过滤逻辑
	for key, value := range filters {
		if value != nil && value != "" {
			switch key {
			case "name", "username", "title":
				// 模糊搜索
				query = query.Where(key+" LIKE ?", "%"+value.(string)+"%")
			case "status":
				// 精确匹配
				query = query.Where(key+" = ?", value)
			default:
				// 其他字段精确匹配
				query = query.Where(key+" = ?", value)
			}
		}
	}

	return query
}

// ApplyPagination 应用分页
func (r *BaseRepository) ApplyPagination(query *gorm.DB, page, pageSize int) *gorm.DB {
	offset := (page - 1) * pageSize
	return query.Offset(offset).Limit(pageSize)
}

// ApplyOrder 应用排序
func (r *BaseRepository) ApplyOrder(query *gorm.DB, orderBy string) *gorm.DB {
	if orderBy == "" {
		orderBy = "created_at DESC"
	}
	return query.Order(orderBy)
}
