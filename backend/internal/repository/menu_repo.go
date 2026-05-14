package repository

import (
	"strings"

	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"

	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) *MenuRepository {
	return &MenuRepository{db: db}
}

// GetAllMenus 获取所有菜单（用于构建树形结构）
func (r *MenuRepository) GetAllMenus(tenantID string) ([]model.Menu, error) {
	var menus []model.Menu
	loadByTenant := func(tid string) ([]model.Menu, error) {
		var rows []model.Menu
		query := r.db.Where("status = ? AND is_show = ?", 1, 1)
		if tid != "" {
			query = query.Where("tenant_id = ?", tid)
		}
		if err := query.Order("sort_order ASC, id ASC").Find(&rows).Error; err != nil {
			return nil, err
		}
		return rows, nil
	}

	rows, err := loadByTenant(tenantID)
	if err != nil {
		return nil, err
	}
	// 租户尚未初始化菜单时，回退默认租户模板菜单，避免登录后左侧菜单为空。
	if len(rows) == 0 && tenantID != "" && tenantID != ids.DefaultTenantUUID {
		rows, err = loadByTenant(ids.DefaultTenantUUID)
		if err != nil {
			return nil, err
		}
	}
	menus = rows

	return menus, nil
}

// GetMenusByParentID 根据父ID获取菜单
func (r *MenuRepository) GetMenusByParentID(parentID *string) ([]model.Menu, error) {
	var menus []model.Menu

	query := r.db.Where("status = ?", 1)

	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}

	err := query.Order("sort_order ASC, id ASC").Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return menus, nil
}

// GetMenuByID 根据ID获取菜单
func (r *MenuRepository) GetMenuByID(id string, tenantID string) (*model.Menu, error) {
	var menu model.Menu
	query := r.db.Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// List 获取菜单列表（带分页和过滤）
func (r *MenuRepository) List(tenantID string, page, pageSize int, name, menuType string, status *int) ([]model.Menu, int64, error) {
	run := func(tid string) ([]model.Menu, int64, error) {
		var rows []model.Menu
		var total int64
		query := r.db.Model(&model.Menu{})
		if tid != "" {
			query = query.Where("tenant_id = ?", tid)
		}
		if name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}
		switch menuType {
		case "catalog":
			query = query.Where("(component = '' OR component IS NULL) AND (permission_code = '' OR permission_code IS NULL)")
		case "menu":
			query = query.Where("component != '' AND component IS NOT NULL")
		case "button":
			query = query.Where("permission_code != '' AND permission_code IS NOT NULL")
		}
		if status != nil {
			query = query.Where("status = ?", *status)
		}
		if err := query.Count(&total).Error; err != nil {
			return nil, 0, err
		}
		offset := (page - 1) * pageSize
		if err := query.Order("sort_order ASC, id ASC").
			Offset(offset).
			Limit(pageSize).
			Find(&rows).Error; err != nil {
			return nil, 0, err
		}
		return rows, total, nil
	}

	menus, total, err := run(tenantID)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 && tenantID != "" && tenantID != ids.DefaultTenantUUID {
		return run(ids.DefaultTenantUUID)
	}
	return menus, total, nil
}

// buildPathChain 根据父节点构建 path_chain。
func buildPathChain(tx *gorm.DB, menuID string, parentID *string) (string, error) {
	if parentID == nil || *parentID == "" {
		return "/" + menuID + "/", nil
	}
	var parent model.Menu
	if err := tx.Select("id, path_chain").Where("id = ?", *parentID).First(&parent).Error; err != nil {
		return "", err
	}
	if parent.PathChain == "" {
		parent.PathChain = "/" + parent.ID + "/"
	}
	return parent.PathChain + menuID + "/", nil
}

func (r *MenuRepository) Create(menu *model.Menu) error {
	if menu.ID == "" {
		menu.ID = ids.New()
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if menu.PathChain == "" {
			chain, err := buildPathChain(tx, menu.ID, menu.ParentID)
			if err != nil {
				return err
			}
			menu.PathChain = chain
		}
		return tx.Create(menu).Error
	})
}

func (r *MenuRepository) Update(menu *model.Menu) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var old model.Menu
		if err := tx.Select("id, parent_id, path_chain").Where("id = ?", menu.ID).First(&old).Error; err != nil {
			return err
		}
		oldParentID := ""
		if old.ParentID != nil {
			oldParentID = *old.ParentID
		}
		newParentID := ""
		if menu.ParentID != nil {
			newParentID = *menu.ParentID
		}
		if oldParentID != newParentID {
			chain, err := buildPathChain(tx, menu.ID, menu.ParentID)
			if err != nil {
				return err
			}
			newChain := chain
			oldChain := old.PathChain
			menu.PathChain = newChain
			// 级联更新所有子菜单的 path_chain 前缀。
			if oldChain != "" {
				var descendants []model.Menu
				if err := tx.Where("path_chain LIKE ? AND id != ?", oldChain+"%", menu.ID).Find(&descendants).Error; err != nil {
					return err
				}
				for i := range descendants {
					updated := strings.Replace(descendants[i].PathChain, oldChain, newChain, 1)
					if err := tx.Model(&model.Menu{}).Where("id = ?", descendants[i].ID).Update("path_chain", updated).Error; err != nil {
						return err
					}
				}
			}
		}
		return tx.Save(menu).Error
	})
}

func (r *MenuRepository) UpdateStatus(id string, tenantID string, status int8) error {
	query := r.db.Model(&model.Menu{}).Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Update("status", status).Error
}

func (r *MenuRepository) UpdateShow(id string, tenantID string, isShow int8) error {
	query := r.db.Model(&model.Menu{}).Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Update("is_show", isShow).Error
}

func (r *MenuRepository) Delete(id string, tenantID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id = ?", id).Delete(&model.MenuPermission{}).Error; err != nil {
			return err
		}
		q := tx.Where("id = ?", id)
		if tenantID != "" {
			q = q.Where("tenant_id = ?", tenantID)
		}
		return q.Delete(&model.Menu{}).Error
	})
}

func (r *MenuRepository) BatchDelete(ids []string, tenantID string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id IN ?", ids).Delete(&model.MenuPermission{}).Error; err != nil {
			return err
		}
		q := tx.Where("id IN ?", ids)
		if tenantID != "" {
			q = q.Where("tenant_id = ?", tenantID)
		}
		return q.Delete(&model.Menu{}).Error
	})
}

// PermissionCodesByMenuIDs 返回每个菜单关联的权限码（含 menu_permissions 与 menus.permission_code）。
func (r *MenuRepository) PermissionCodesByMenuIDs(menuIDs []string) (map[string][]string, error) {
	out := make(map[string][]string)
	if len(menuIDs) == 0 {
		return out, nil
	}
	type joinRow struct {
		MenuID string `gorm:"column:menu_id"`
		Code   string `gorm:"column:code"`
	}
	var jrows []joinRow
	if err := r.db.Table("menu_permissions AS mp").
		Select("mp.menu_id, p.code").
		Joins("JOIN permissions AS p ON p.id = mp.permission_id").
		Where("mp.menu_id IN ?", menuIDs).
		Scan(&jrows).Error; err != nil {
		return nil, err
	}
	for _, row := range jrows {
		if row.Code == "" {
			continue
		}
		out[row.MenuID] = appendUniqueCode(out[row.MenuID], row.Code)
	}
	var menus []model.Menu
	if err := r.db.Select("id", "permission_code").Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
		return nil, err
	}
	for _, m := range menus {
		pc := strings.TrimSpace(m.PermissionCode)
		if pc == "" {
			continue
		}
		out[m.ID] = appendUniqueCode(out[m.ID], pc)
	}
	return out, nil
}

func appendUniqueCode(list []string, code string) []string {
	for _, x := range list {
		if x == code {
			return list
		}
	}
	return append(list, code)
}

// ReplaceMenuPermissionLinks 用权限码列表重写某菜单的 menu_permissions 行。
func (r *MenuRepository) ReplaceMenuPermissionLinks(menuID string, permissionCodes []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("menu_id = ?", menuID).Delete(&model.MenuPermission{}).Error; err != nil {
			return err
		}
		if len(permissionCodes) == 0 {
			return nil
		}
		var perms []model.Permission
		if err := tx.Where("code IN ?", permissionCodes).Find(&perms).Error; err != nil {
			return err
		}
		for i := range perms {
			row := model.MenuPermission{
				ID:           ids.New(),
				MenuID:       menuID,
				PermissionID: perms[i].ID,
			}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
