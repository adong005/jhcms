package repository

import (
	"errors"
	"time"

	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/datascope"

	"gorm.io/gorm"
)

type SystemLogRepository struct {
	db *gorm.DB
}

func NewSystemLogRepository(db *gorm.DB) *SystemLogRepository {
	return &SystemLogRepository{db: db}
}

// ListParams 日志列表查询参数。
type SystemLogListParams struct {
	TenantID     string
	Role         string
	UserID       string
	UserPath     string
	Username     string
	Page         int
	PageSize     int
	UsernameKw   string
	Usernames    []string
	Action       string
	Status       string
	Module       string
	IP           string
	LogType      string
	StartTime    *time.Time
	EndTime      *time.Time
	TenantFilter string
}

func (r *SystemLogRepository) List(p SystemLogListParams) ([]model.SystemLog, int64, error) {
	var items []model.SystemLog
	var total int64
	query := r.db.Model(&model.SystemLog{})

	switch p.Role {
	case "super_admin":
		if p.TenantFilter != "" {
			query = query.Where("tenant_id = ?", p.TenantFilter)
		}
	case "admin":
		if p.TenantID != "" {
			query = query.Where("tenant_id = ?", p.TenantID)
		}
		// 用 path 子树隔离：admin 及其所有层级子用户的日志。
		if p.UserPath != "" && p.UserPath != "/" {
			query = datascope.ApplyUserPathSubtree(query, "path", p.UserPath)
		} else if p.UserID != "" {
			query = query.Where("user_id = ? OR parent_id = ?", p.UserID, p.UserID)
		}
	default:
		// 普通用户仅看自己的日志。
		if p.TenantID != "" {
			query = query.Where("tenant_id = ?", p.TenantID)
		}
		query = query.Where("user_id = ?", p.UserID)
	}

	if p.UsernameKw != "" {
		query = query.Where("username LIKE ?", "%"+p.UsernameKw+"%")
	}
	if len(p.Usernames) > 0 {
		query = query.Where("username IN ?", p.Usernames)
	}
	if p.Action != "" {
		query = query.Where("action = ?", p.Action)
	}
	if p.Status != "" {
		query = query.Where("status = ?", p.Status)
	}
	if p.Module != "" {
		query = query.Where("module LIKE ?", "%"+p.Module+"%")
	}
	if p.IP != "" {
		query = query.Where("ip LIKE ?", p.IP+"%")
	}
	if p.LogType != "" {
		query = query.Where("log_type = ?", p.LogType)
	}
	if p.StartTime != nil {
		query = query.Where("created_at >= ?", *p.StartTime)
	}
	if p.EndTime != nil {
		query = query.Where("created_at <= ?", *p.EndTime)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (p.Page - 1) * p.PageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(p.PageSize).Find(&items).Error; err != nil {
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

// ClearAll 清空日志。tenantID 为空时拒绝操作，防止误全表删除。
// super_admin 如需清全平台请传 "__all__" 作为 tenantID。
func (r *SystemLogRepository) ClearAll(tenantID string) error {
	if tenantID == "" {
		return errors.New("tenant_id is required for clear operation")
	}
	if tenantID == "__all__" {
		return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.SystemLog{}).Error
	}
	return r.db.Where("tenant_id = ?", tenantID).Delete(&model.SystemLog{}).Error
}

// PurgeBefore 删除 before 时刻之前的日志（归档/清理用）。
func (r *SystemLogRepository) PurgeBefore(tenantID string, before time.Time) (int64, error) {
	query := r.db.Where("created_at < ?", before)
	if tenantID != "" && tenantID != "__all__" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	result := query.Delete(&model.SystemLog{})
	return result.RowsAffected, result.Error
}

func (r *SystemLogRepository) Create(item *model.SystemLog) error {
	return r.db.Create(item).Error
}
