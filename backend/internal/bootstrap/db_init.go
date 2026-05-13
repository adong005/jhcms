package bootstrap

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func InitDatabase(db *gorm.DB, mode string) error {
	if strings.EqualFold(mode, "reset") {
		if err := resetAllTables(db); err != nil {
			return err
		}
	}

	if err := MigrateSchema(db); err != nil {
		return err
	}

	if err := normalizeRolePermissionSchema(db); err != nil {
		return err
	}
	return seedDefaultData(db)
}

func MigrateSchema(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.Permission{},
		&model.User{},
		&model.Role{},
		&model.RolePermission{},
		&model.Menu{},
		&model.SystemLog{},
		&model.TenantSiteConfig{},
	); err != nil {
		return err
	}
	// users 不再保留 tenant_id，由 users.id / users.parent_id 推导租户关系。
	if db.Migrator().HasColumn(&model.User{}, "tenant_id") {
		if err := db.Migrator().DropColumn(&model.User{}, "tenant_id"); err != nil {
			return err
		}
	}
	// 超管标识统一为 is_admin，兼容历史列 is_platform_super_admin。
	if db.Migrator().HasColumn(&model.User{}, "is_platform_super_admin") {
		if db.Migrator().HasColumn(&model.User{}, "is_admin") {
			if err := db.Exec("UPDATE users SET is_admin = 1 WHERE is_platform_super_admin = 1").Error; err != nil {
				return err
			}
		}
		if err := db.Migrator().DropColumn(&model.User{}, "is_platform_super_admin"); err != nil {
			return err
		}
	}
	// tenants 表已废弃，迁移时自动清理。
	if db.Migrator().HasTable("tenants") {
		if err := db.Migrator().DropTable("tenants"); err != nil {
			return err
		}
	}
	// 历史数据纠正：管理员创建但仍落在默认租户的角色，tenant_id 归正为创建人用户ID。
	if err := normalizeRoleTenantID(db); err != nil {
		return err
	}
	// 历史数据兼容：菜单显示字段默认按显示处理，避免升级后菜单意外不可见。
	if db.Migrator().HasColumn(&model.Menu{}, "is_show") {
		if err := db.Exec("UPDATE menus SET is_show = 1 WHERE is_show IS NULL").Error; err != nil {
			return err
		}
	}
	return nil
}

func normalizeRoleTenantID(db *gorm.DB) error {
	return db.Exec(
		`UPDATE roles r
		INNER JOIN users u ON u.id = r.created_by
		SET r.tenant_id = r.created_by
		WHERE r.tenant_id = ? AND u.role = 'admin' AND r.created_by IS NOT NULL AND r.created_by <> ''`,
		ids.DefaultTenantUUID,
	).Error
}

func resetAllTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&model.RolePermission{},
		&model.SystemLog{},
		&model.Menu{},
		&model.TenantSiteConfig{},
		&model.User{},
		&model.Role{},
		&model.Permission{},
	)
}

func normalizeRolePermissionSchema(db *gorm.DB) error {
	if db.Migrator().HasIndex(&model.RolePermission{}, "idx_role_menu_unique") {
		if err := db.Migrator().DropIndex(&model.RolePermission{}, "idx_role_menu_unique"); err != nil {
			return err
		}
	}
	// 清理 permission_id 为空导致的唯一索引冲突数据。
	if err := db.Where("permission_id = '' OR permission_id IS NULL").Delete(&model.RolePermission{}).Error; err != nil {
		return err
	}
	// 去重，保留同 role_id + permission_id 的一条记录。
	if err := db.Exec(
		`DELETE rp1 FROM role_permissions rp1
		INNER JOIN role_permissions rp2
		ON rp1.role_id = rp2.role_id
		AND rp1.permission_id = rp2.permission_id
		AND rp1.id > rp2.id`,
	).Error; err != nil {
		return err
	}
	if db.Migrator().HasColumn(&model.RolePermission{}, "menu_id") {
		if err := db.Migrator().DropColumn(&model.RolePermission{}, "menu_id"); err != nil {
			return err
		}
	}
	if !db.Migrator().HasIndex(&model.RolePermission{}, "idx_role_permission_unique") {
		if err := db.Migrator().CreateIndex(&model.RolePermission{}, "idx_role_permission_unique"); err != nil {
			return err
		}
	}
	return nil
}

func seedDefaultData(db *gorm.DB) error {
	if err := seedDefaultRoles(db); err != nil {
		return err
	}
	permissionCodeToID, err := seedDefaultPermissions(db)
	if err != nil {
		return err
	}
	if err := seedDefaultMenus(db); err != nil {
		return err
	}
	if err := seedDefaultAdmin(db); err != nil {
		return err
	}
	return seedRolePermissions(db, permissionCodeToID)
}

func seedDefaultRoles(db *gorm.DB) error {
	defaultRoles := []model.Role{
		{TenantScoped: model.TenantScoped{TenantID: ids.DefaultTenantUUID}, Name: "超级管理员", Code: "super_admin", DataScope: "TENANT_ALL", Description: "系统超级管理员", Status: 1},
		{TenantScoped: model.TenantScoped{TenantID: ids.DefaultTenantUUID}, Name: "管理员", Code: "admin", DataScope: "TENANT_ALL", Description: "租户管理员", Status: 1},
		{TenantScoped: model.TenantScoped{TenantID: ids.DefaultTenantUUID}, Name: "普通用户", Code: "user", DataScope: "SELF", Description: "普通用户", Status: 1},
	}
	for i := range defaultRoles {
		role := defaultRoles[i]
		var existing model.Role
		err := db.Where("tenant_id = ? AND code = ?", role.TenantID, role.Code).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		role.ID = ids.New()
		if err := db.Create(&role).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedDefaultPermissions(db *gorm.DB) (map[string]string, error) {
	// 统一维护权限点，覆盖已接入 withPermission 的接口，并补齐核心业务模块权限。
	permissions := []model.Permission{
		// 系统管理（路由 withPermission 已使用）
		{Code: "system:user:list", Name: "用户列表", Module: "system"},
		{Code: "system:user:create", Name: "创建用户", Module: "system"},
		{Code: "system:user:update", Name: "编辑用户", Module: "system"},
		{Code: "system:user:delete", Name: "删除用户", Module: "system"},
		{Code: "system:role:list", Name: "角色列表", Module: "system"},
		{Code: "system:role:create", Name: "创建角色", Module: "system"},
		{Code: "system:role:update", Name: "编辑角色", Module: "system"},
		{Code: "system:role:permission", Name: "角色授权", Module: "system", IsDelegable: false},
		{Code: "system:menu:list", Name: "菜单列表", Module: "system"},
		{Code: "system:menu:update", Name: "编辑菜单", Module: "system"},
		{Code: "system:permission:list", Name: "权限管理", Module: "system"},
		// 日志
		{Code: "log:list", Name: "日志列表", Module: "log"},
		{Code: "log:delete", Name: "删除日志", Module: "log"},
		{Code: "log:clear", Name: "清空日志", Module: "log"},
	}
	codeToID := make(map[string]string, len(permissions))
	for i := range permissions {
		p := permissions[i]
		var existing model.Permission
		err := db.Where("code = ?", p.Code).First(&existing).Error
		if err == nil {
			codeToID[p.Code] = existing.ID
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		p.ID = ids.New()
		if err := db.Create(&p).Error; err != nil {
			return nil, err
		}
		codeToID[p.Code] = p.ID
	}
	return codeToID, nil
}

func seedDefaultMenus(db *gorm.DB) error {
	tenantID := ids.DefaultTenantUUID
	_ = db.Where("tenant_id = ?", tenantID).Delete(&model.Menu{}).Error

	type menuSeed struct {
		Name      string
		Path      string
		Component string
		Icon      string
		SortOrder int
		ParentID  *string
	}

	systemID := ids.New()
	dashboardID := ids.New()

	seeds := []struct {
		ID string
		menuSeed
	}{
		{dashboardID, menuSeed{Name: "工作台", Path: "/dashboard", Icon: "mdi:view-dashboard", SortOrder: 1}},
		{systemID, menuSeed{Name: "系统管理", Path: "/system", Icon: "mdi:cog", SortOrder: 2}},
		{ids.New(), menuSeed{Name: "个人中心", Path: "/profile", Component: "/_core/profile/index", Icon: "mdi:account-circle", SortOrder: 4}},
		{ids.New(), menuSeed{Name: "日志管理", Path: "/system-logs/list", Component: "/system-logs/list", Icon: "lucide:file-text", SortOrder: 5}},
		{ids.New(), menuSeed{Name: "概览", Path: "/analytics", Component: "/dashboard/analytics/index", Icon: "mdi:view-dashboard", SortOrder: 1, ParentID: &dashboardID}},
		{ids.New(), menuSeed{Name: "用户管理", Path: "/users/list", Component: "/users/list", Icon: "lucide:users", SortOrder: 1, ParentID: &systemID}},
		{ids.New(), menuSeed{Name: "角色管理", Path: "/roles/list", Component: "/roles/list", Icon: "lucide:shield", SortOrder: 2, ParentID: &systemID}},
		{ids.New(), menuSeed{Name: "菜单管理", Path: "/menus/list", Component: "/menus/list", Icon: "lucide:menu", SortOrder: 3, ParentID: &systemID}},
		{ids.New(), menuSeed{Name: "权限管理", Path: "/permissions/list", Component: "/permissions/list", Icon: "lucide:key-round", SortOrder: 4, ParentID: &systemID}},
	}

	for _, s := range seeds {
		m := model.Menu{
			ID:           s.ID,
			TenantScoped: model.TenantScoped{TenantID: tenantID},
			Name:         s.Name,
			Path:         s.Path,
			Component:    s.Component,
			Icon:         s.Icon,
			SortOrder:    s.SortOrder,
			IsShow:       1,
			Status:       1,
			ParentID:     s.ParentID,
		}
		if err := db.Create(&m).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedDefaultAdmin(db *gorm.DB) error {
	type userSeed struct {
		Username      string
		PasswordPlain string
		IsAdmin       bool
		RealName      string
		NickName      string
		Email         string
		Phone         string
		Role          string
		DataScope     string
		Status        int8
	}

	seeds := []userSeed{
		{Username: "admin", PasswordPlain: "admin123", IsAdmin: true, RealName: "超级管理员", NickName: "Admin", Email: "admin@adcms.com", Phone: "13800000001", Role: "super_admin", DataScope: "TENANT_ALL", Status: 1},
		{Username: "tenant_admin", PasswordPlain: "admin123", IsAdmin: false, RealName: "租户管理员", NickName: "TenantAdmin", Email: "tenant-admin@adcms.com", Phone: "13800000002", Role: "admin", DataScope: "TENANT_ALL", Status: 1},
		{Username: "demo_user", PasswordPlain: "admin123", IsAdmin: false, RealName: "演示用户", NickName: "DemoUser", Email: "demo-user@adcms.com", Phone: "13800000003", Role: "user", DataScope: "SELF", Status: 1},
	}

	for _, s := range seeds {
		var existing model.User
		err := db.Where("username = ?", s.Username).First(&existing).Error
		if err == nil {
			// 已存在用户时补齐可展示的模拟字段，不覆盖历史密码。
			updateData := map[string]interface{}{
				"is_admin":   s.IsAdmin,
				"real_name":  s.RealName,
				"nick_name":  s.NickName,
				"email":      s.Email,
				"phone":      s.Phone,
				"role":       s.Role,
				"data_scope": s.DataScope,
				"status":     s.Status,
			}
			if err := db.Model(&model.User{}).Where("id = ?", existing.ID).Updates(updateData).Error; err != nil {
				return err
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		pwd, hashErr := bcrypt.GenerateFromPassword([]byte(s.PasswordPlain), bcrypt.DefaultCost)
		if hashErr != nil {
			return hashErr
		}

		u := model.User{
			ID:        ids.New(),
			Username:  s.Username,
			Password:  string(pwd),
			IsAdmin:   s.IsAdmin,
			RealName:  s.RealName,
			NickName:  s.NickName,
			Email:     s.Email,
			Phone:     s.Phone,
			Role:      s.Role,
			DataScope: s.DataScope,
			Status:    s.Status,
		}
		if s.Role == "user" {
			// 默认普通用户归属 tenant_admin
			adminUsername := "tenant_admin"
			var admin model.User
			if err := db.Select("id").Where("username = ?", adminUsername).First(&admin).Error; err == nil {
				u.ParentID = &admin.ID
			}
		}
		if err := db.Create(&u).Error; err != nil {
			return err
		}
	}

	return nil
}

func seedRolePermissions(db *gorm.DB, codeToID map[string]string) error {
	rolePerms := map[string][]string{
		"super_admin": {
			"system:user:list", "system:user:create", "system:user:update", "system:user:delete",
			"system:role:list", "system:role:create", "system:role:update", "system:role:permission",
			"system:menu:list", "system:menu:update", "system:permission:list",
			"log:list", "log:delete", "log:clear",
		},
		"admin": {
			"system:user:list", "system:user:create", "system:user:update", "system:user:delete",
			"system:role:list", "system:role:create", "system:role:update", "system:role:permission",
			"system:menu:list", "system:menu:update", "system:permission:list",
			"log:list", "log:delete", "log:clear",
		},
		"user": {},
	}
	for roleCode, codes := range rolePerms {
		var role model.Role
		if err := db.Where("tenant_id = ? AND code = ?", ids.DefaultTenantUUID, roleCode).First(&role).Error; err != nil {
			continue
		}
		if err := db.Where("role_id = ?", role.ID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}
		rows := make([]model.RolePermission, 0, len(codes))
		for _, code := range codes {
			pid, ok := codeToID[code]
			if !ok || pid == "" {
				continue
			}
			rows = append(rows, model.RolePermission{
				ID:           ids.New(),
				RoleID:       role.ID,
				PermissionID: pid,
			})
		}
		if len(rows) > 0 {
			if err := db.Create(&rows).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
