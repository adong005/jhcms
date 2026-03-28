package repository

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"

	"gorm.io/gorm"
)

type InfoCategoryRepository struct {
	db *gorm.DB
}

func NewInfoCategoryRepository(db *gorm.DB) *InfoCategoryRepository {
	return &InfoCategoryRepository{db: db}
}

func (r *InfoCategoryRepository) List(tenantID, role, currentUserID string, page, pageSize int, name string, status *int) ([]model.InfoCategory, int64, error) {
	var items []model.InfoCategory
	var total int64
	query := r.db.Model(&model.InfoCategory{})
	if role == "super_admin" {
		// 超管全量
	} else if currentUserID != "" {
		// 管理员仅看自己创建的分类；普通用户同样仅看自己创建的分类。
		query = query.Where("created_by = ?", currentUserID)
	} else if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	if err := query.Order("sort ASC, id ASC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// MapDisplayLabelByIDs 批量查询分类，返回 id -> 展示文案
func (r *InfoCategoryRepository) MapDisplayLabelByIDs(tenantID string, ids []string) (map[string]string, error) {
	dedup := dedupeUUIDs(ids)
	if len(dedup) == 0 {
		return map[string]string{}, nil
	}
	var rows []model.InfoCategory
	q := r.db.Model(&model.InfoCategory{}).Select("id", "name", "code").Where("id IN ?", dedup)
	if tenantID != "" {
		q = q.Where("tenant_id = ?", tenantID)
	}
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make(map[string]string, len(rows))
	for _, row := range rows {
		if row.Code != "" {
			out[row.ID] = row.Name + "（" + row.Code + "）"
		} else {
			out[row.ID] = row.Name
		}
	}
	return out, nil
}

func (r *InfoCategoryRepository) GetByID(id string, tenantID string) (*model.InfoCategory, error) {
	var item model.InfoCategory
	query := r.db.Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *InfoCategoryRepository) Create(item *model.InfoCategory) error {
	if item.ID == "" {
		item.ID = ids.New()
	}
	return r.db.Create(item).Error
}
func (r *InfoCategoryRepository) Update(item *model.InfoCategory) error { return r.db.Save(item).Error }
func (r *InfoCategoryRepository) UpdateStatus(id string, tenantID string, status int8) error {
	query := r.db.Model(&model.InfoCategory{}).Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Update("status", status).Error
}
func (r *InfoCategoryRepository) Delete(id string, tenantID string) error {
	query := r.db.Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Delete(&model.InfoCategory{}).Error
}
func (r *InfoCategoryRepository) BatchDelete(ids []string, tenantID string) error {
	if len(ids) == 0 {
		return nil
	}
	query := r.db.Where("id IN ?", ids)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Delete(&model.InfoCategory{}).Error
}
