package repository

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) *PermissionRepository {
	return &PermissionRepository{db: db}
}

func (r *PermissionRepository) List(page, pageSize int, name, code, module string) ([]model.Permission, int64, error) {
	var permissions []model.Permission
	var total int64

	query := r.db.Model(&model.Permission{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if code != "" {
		query = query.Where("code LIKE ?", "%"+code+"%")
	}
	if module != "" {
		query = query.Where("module LIKE ?", "%"+module+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}
	return permissions, total, nil
}

func (r *PermissionRepository) GetByID(id string) (*model.Permission, error) {
	var permission model.Permission
	if err := r.db.Where("id = ?", id).First(&permission).Error; err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *PermissionRepository) Create(permission *model.Permission) error {
	if permission.ID == "" {
		permission.ID = ids.New()
	}
	return r.db.Create(permission).Error
}

func (r *PermissionRepository) Update(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

func (r *PermissionRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Permission{}).Error
}

func (r *PermissionRepository) BatchDelete(idsIn []string) error {
	if len(idsIn) == 0 {
		return nil
	}
	return r.db.Where("id IN ?", idsIn).Delete(&model.Permission{}).Error
}
