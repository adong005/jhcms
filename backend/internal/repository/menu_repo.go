package repository

import (
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
		query = query.Where("parent_id_menu IS NULL")
	} else {
		query = query.Where("parent_id_menu = ?", *parentID)
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
		if menuType != "" {
			if menuType == "catalog" {
				query = query.Where("component = '' OR component IS NULL")
			} else if menuType == "menu" {
				query = query.Where("component != '' AND component IS NOT NULL")
			}
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

func (r *MenuRepository) Create(menu *model.Menu) error {
	if menu.ID == "" {
		menu.ID = ids.New()
	}
	return r.db.Create(menu).Error
}

func (r *MenuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
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
	query := r.db.Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Delete(&model.Menu{}).Error
}

func (r *MenuRepository) BatchDelete(ids []string, tenantID string) error {
	if len(ids) == 0 {
		return nil
	}
	query := r.db.Where("id IN ?", ids)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Delete(&model.Menu{}).Error
}
