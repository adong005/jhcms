package repository

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"

	"gorm.io/gorm"
)

type InfoRepository struct {
	db *gorm.DB
}

func NewInfoRepository(db *gorm.DB) *InfoRepository {
	return &InfoRepository{db: db}
}

// List 获取信息列表（带租户隔离）
func (r *InfoRepository) List(tenantID string, role string, userID string, page, pageSize int, title string, status *int) ([]model.Info, int64, error) {
	var infos []model.Info
	var total int64

	query := r.db.Model(&model.Info{})

	if role == "super_admin" {
		// 超管全量
	} else if userID != "" {
		query = query.Where("created_by = ? OR created_by IN (?)",
			userID,
			r.db.Model(&model.User{}).Select("id").Where("created_by = ?", userID),
		)
	} else if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&infos).Error; err != nil {
		return nil, 0, err
	}

	return infos, total, nil
}

// GetByID 根据ID获取信息
func (r *InfoRepository) GetByID(id string, tenantID string) (*model.Info, error) {
	var info model.Info
	query := r.db.Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.First(&info).Error; err != nil {
		return nil, err
	}
	return &info, nil
}

// Create 创建信息
func (r *InfoRepository) Create(tenantID, userID string, title, content string, categoryID *string, status int8, publishDate string) error {
	cb := userID
	info := &model.Info{
		ID:              ids.New(),
		TenantScoped:    model.TenantScoped{TenantID: tenantID},
		Title:           title,
		Content:         content,
		CategoryID:      categoryID,
		Status:          status,
		CreatorOptional: model.CreatorOptional{CreatedBy: &cb},
	}

	return r.db.Create(info).Error
}

// Update 更新信息
func (r *InfoRepository) Update(id string, tenantID string, title, content string, categoryID *string, status int8, publishDate, author, summary string) error {
	updates := map[string]interface{}{
		"title":       title,
		"content":     content,
		"category_id": categoryID,
		"status":      status,
		"author":      author,
		"summary":     summary,
	}

	return r.db.Model(&model.Info{}).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		Updates(updates).Error
}

// Delete 删除信息（软删除）
func (r *InfoRepository) Delete(id string, tenantID string) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Info{}).Error
}

func (r *InfoRepository) UpdateStatus(id string, tenantID string, status int8) error {
	return r.db.Model(&model.Info{}).Where("id = ? AND tenant_id = ?", id, tenantID).Update("status", status).Error
}

func (r *InfoRepository) BatchDelete(ids []string, tenantID string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Where("id IN ? AND tenant_id = ?", ids, tenantID).Delete(&model.Info{}).Error
}
