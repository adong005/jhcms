package repository

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) List(tenantID, operatorRole, currentUserID string, page, pageSize int, name, code string, status *int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	run := func(tid string) ([]model.Role, int64, error) {
		var rows []model.Role
		var cnt int64
		query := r.db.Model(&model.Role{})
		if tid != "" {
			query = query.Where("tenant_id = ?", tid)
		}
		if name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}
		if code != "" {
			query = query.Where("code LIKE ?", "%"+code+"%")
		}
		if operatorRole == "admin" {
			query = query.Where("created_by = ?", currentUserID)
		}
		if status != nil {
			query = query.Where("status = ?", *status)
		}
		if err := query.Count(&cnt).Error; err != nil {
			return nil, 0, err
		}
		offset := (page - 1) * pageSize
		if err := query.Order("created_at ASC").Offset(offset).Limit(pageSize).Find(&rows).Error; err != nil {
			return nil, 0, err
		}
		return rows, cnt, nil
	}

	roles, total, err := run(tenantID)
	if err != nil {
		return nil, 0, err
	}
	// 租户尚未初始化角色时，回退默认租户模板角色，避免角色列表为空。
	if total == 0 && operatorRole != "admin" && tenantID != "" && tenantID != ids.DefaultTenantUUID {
		return run(ids.DefaultTenantUUID)
	}
	return roles, total, nil
}

func (r *RoleRepository) GetByID(id string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) Create(role *model.Role) error {
	if role.ID == "" {
		role.ID = ids.New()
	}
	return r.db.Create(role).Error
}

func (r *RoleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) UpdateStatus(id string, status int8) error {
	return r.db.Model(&model.Role{}).Where("id = ?", id).Update("status", status).Error
}

func (r *RoleRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Role{}).Error
}

func (r *RoleRepository) BatchDelete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Where("id IN ?", ids).Delete(&model.Role{}).Error
}

func (r *RoleRepository) GetPermissionIDs(roleID string) ([]string, error) {
	var rows []model.RolePermission
	if err := r.db.Where("role_id = ?", roleID).Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]string, 0, len(rows))
	for _, row := range rows {
		out = append(out, row.PermissionID)
	}
	return out, nil
}

func (r *RoleRepository) SetPermissions(roleID string, permissionIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}
		if len(permissionIDs) == 0 {
			return nil
		}
		rows := make([]model.RolePermission, 0, len(permissionIDs))
		for _, permissionID := range permissionIDs {
			rows = append(rows, model.RolePermission{
				ID:           ids.New(),
				RoleID:       roleID,
				PermissionID: permissionID,
			})
		}
		return tx.Create(&rows).Error
	})
}

func (r *RoleRepository) GetByCode(tenantID string, code string) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("tenant_id = ? AND code = ?", tenantID, code).First(&role).Error; err != nil {
		if tenantID != "" && tenantID != ids.DefaultTenantUUID {
			if err2 := r.db.Where("tenant_id = ? AND code = ?", ids.DefaultTenantUUID, code).First(&role).Error; err2 == nil {
				return &role, nil
			}
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetDelegablePermissionIDsByRoleID(roleID string) ([]string, error) {
	var idsOut []string
	err := r.db.Table("role_permissions rp").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("rp.role_id = ? AND p.is_delegable = ?", roleID, true).
		Pluck("rp.permission_id", &idsOut).Error
	if err != nil {
		return nil, err
	}
	return idsOut, nil
}
