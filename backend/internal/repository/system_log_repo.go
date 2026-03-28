package repository

import (
	"adcms-backend/internal/model"

	"gorm.io/gorm"
)

type SystemLogRepository struct {
	db *gorm.DB
}

func NewSystemLogRepository(db *gorm.DB) *SystemLogRepository {
	return &SystemLogRepository{db: db}
}

func (r *SystemLogRepository) List(
	tenantID, role, currentUserID, currentUsername string,
	page, pageSize int,
	username, action, status, tenantFilter string,
	usernames []string,
) ([]model.SystemLog, int64, error) {
	var items []model.SystemLog
	var total int64
	query := r.db.Model(&model.SystemLog{})

	switch role {
	case "super_admin":
		if tenantFilter != "" {
			query = query.Where("tenant_id = ?", tenantFilter)
		}
	case "admin":
		if tenantID != "" {
			query = query.Where("tenant_id = ?", tenantID)
		}
		// 管理员仅看自己和下级日志（下级日志 parent_id 指向管理员 user_id）。
		query = query.Where("username = ? OR parent_id = ?", currentUsername, currentUserID)
	default:
		// 普通用户仅看自己的日志。
		if tenantID != "" {
			query = query.Where("tenant_id = ?", tenantID)
		}
		query = query.Where("username = ?", currentUsername)
	}

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if len(usernames) > 0 {
		query = query.Where("username IN ?", usernames)
	}
	if action != "" {
		query = query.Where("action LIKE ?", "%"+action+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *SystemLogRepository) Delete(id string, tenantID string) error {
	query := r.db.Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Delete(&model.SystemLog{}).Error
}
func (r *SystemLogRepository) BatchDelete(ids []string, tenantID string) error {
	if len(ids) == 0 {
		return nil
	}
	query := r.db.Where("id IN ?", ids)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Delete(&model.SystemLog{}).Error
}
func (r *SystemLogRepository) ClearAll(tenantID string) error {
	query := r.db.Model(&model.SystemLog{})
	if tenantID != "" {
		return query.Where("tenant_id = ?", tenantID).Delete(&model.SystemLog{}).Error
	}
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.SystemLog{}).Error
}

func (r *SystemLogRepository) Create(item *model.SystemLog) error {
	return r.db.Create(item).Error
}
