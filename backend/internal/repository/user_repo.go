package repository

import (
	"fmt"
	"strings"

	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func dedupeUUIDs(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	var out []string
	for _, s := range in {
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// FindByUsername 根据用户名查找用户
func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID 根据 ID 查找用户
func (r *UserRepository) FindByID(id string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	if user.ID == "" {
		user.ID = ids.New()
	}
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// UpdateLastLogin 更新最后登录时间
func (r *UserRepository) UpdateLastLogin(userID string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).
		Update("last_login_date", gorm.Expr("NOW()")).Error
}

// GetSiteConfigByUserID 根据用户ID获取网站配置（支持 super_admin 和 admin）
func (r *UserRepository) GetSiteConfigByUserID(userID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetSiteConfigByTenantAdmin 租户内普通用户读取：取该租户下 role=admin 的站点配置
func (r *UserRepository) GetSiteConfigByTenantAdmin(tenantID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("id = ? AND role = ?", tenantID, "admin").First(&user).Error
	if err != nil {
		// 兼容历史数据：若 tenantID 指向普通用户，则回溯其 parent(admin)。
		var member model.User
		if e := r.db.Select("parent_id").Where("id = ?", tenantID).First(&member).Error; e == nil &&
			member.ParentID != nil && *member.ParentID != "" {
			return r.GetSiteConfigByUserID(*member.ParentID)
		}
		return nil, err
	}
	return r.GetSiteConfigByUserID(user.ID)
}

// UpdateSiteConfig 更新网站配置
func (r *UserRepository) UpdateSiteConfig(
	userID string,
	title, keywords, description, domain, logo, icpCode, contactPhone, contactAddress, contactEmail string,
) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Select("title", "keywords", "description", "domain", "logo", "icp_code", "contact_phone", "contact_address", "contact_email").
		Updates(map[string]interface{}{
			"title":           title,
			"keywords":        keywords,
			"description":     description,
			"domain":          domain,
			"logo":            logo,
			"icp_code":        icpCode,
			"contact_phone":   contactPhone,
			"contact_address": contactAddress,
			"contact_email":   contactEmail,
		}).Error
}

// UpdateSiteLogo 更新网站 Logo 路径
func (r *UserRepository) UpdateSiteLogo(userID, logo string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("logo", logo).Error
}

// List 获取用户列表（带租户隔离）
func (r *UserRepository) List(tenantID, role, currentUserID string, page, pageSize int, username string, status *int) ([]model.User, int64, error) {
	query := r.db.Model(&model.User{})

	if role == "super_admin" {
		// 平台超级管理员可查看全部用户，不按租户过滤。
	} else if currentUserID != "" {
		query = query.Where("id = ? OR created_by = ?", currentUserID, currentUserID)
	}

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	if status != nil {
		query = query.Where("status = ?", int8(*status))
	}

	var total int64
	query.Count(&total)

	var users []model.User
	err := query.Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

func (r *UserRepository) UpdateStatus(id string, status int8) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *UserRepository) UpdatePasswordHash(id string, passwordHash string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", passwordHash).Error
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.User{}).Error
}

func (r *UserRepository) BatchDelete(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Where("id IN ?", ids).Delete(&model.User{}).Error
}

func (r *UserRepository) GetAccessCodesByUser(userID string, role string) ([]string, error) {
	if role == "super_admin" {
		var allCodes []string
		if err := r.db.Model(&model.Permission{}).
			Where("code <> ''").
			Pluck("code", &allCodes).Error; err != nil {
			return nil, err
		}
		return allCodes, nil
	}

	var user model.User
	if err := r.db.Select("id, role, parent_id, is_admin").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	tenantID := user.EffectiveTenantID()

	var roleRecord model.Role
	if err := r.db.Where("tenant_id = ? AND code = ?", tenantID, user.Role).First(&roleRecord).Error; err != nil {
		// 租户未初始化独立角色时，回退默认租户模板角色
		if err2 := r.db.Where("tenant_id = ? AND code = ?", ids.DefaultTenantUUID, user.Role).First(&roleRecord).Error; err2 != nil {
			return []string{}, nil
		}
	}

	var codes []string
	err := r.db.Table("role_permissions rp").
		Joins("JOIN permissions p ON p.id = rp.permission_id").
		Where("rp.role_id = ? AND p.code <> ''", roleRecord.ID).
		Distinct("p.code").
		Pluck("p.code", &codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}

// NickNameByUserIDs 批量解析用户展示名
func (r *UserRepository) NickNameByUserIDs(ids []string) map[string]string {
	dedup := dedupeUUIDs(ids)
	out := make(map[string]string, len(dedup))
	if len(dedup) == 0 {
		return out
	}
	var users []model.User
	if err := r.db.Model(&model.User{}).
		Select("id", "nick_name", "real_name", "username").
		Where("id IN ?", dedup).
		Find(&users).Error; err != nil {
		return out
	}
	for _, u := range users {
		label := strings.TrimSpace(u.NickName)
		if label == "" {
			label = strings.TrimSpace(u.RealName)
		}
		if label == "" {
			label = strings.TrimSpace(u.Username)
		}
		out[u.ID] = label
	}
	return out
}

// AscriptionByTenantIDs 列表「归属」：租户 admin 昵称等
func (r *UserRepository) AscriptionByTenantIDs(tenantIDs []string) map[string]string {
	out := make(map[string]string)
	for _, tid := range dedupeUUIDs(tenantIDs) {
		var u model.User
		err := r.db.Where("id = ? AND role = ?", tid, "admin").First(&u).Error
		if err != nil {
			_ = r.db.Where("parent_id = ?", tid).Order("created_at ASC").First(&u).Error
		}
		label := strings.TrimSpace(u.NickName)
		if label == "" {
			label = strings.TrimSpace(u.RealName)
		}
		if label == "" && u.Username != "" {
			label = u.Username
		}
		if label == "" {
			suffix := tid
			if len(suffix) > 8 {
				suffix = suffix[:8]
			}
			label = fmt.Sprintf("租户#%s", suffix)
		}
		out[tid] = label
	}
	return out
}

// RoleExistsByCode 判断角色编码是否存在且启用
func (r *UserRepository) RoleExistsByCode(code string) (bool, error) {
	var cnt int64
	if err := r.db.Model(&model.Role{}).
		Where("code = ? AND status = 1", code).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// RoleExistsByCodeAndCreator 判断角色编码是否由指定创建人创建且启用
func (r *UserRepository) RoleExistsByCodeAndCreator(code, creatorID string) (bool, error) {
	var cnt int64
	if err := r.db.Model(&model.Role{}).
		Where("code = ? AND created_by = ? AND status = 1", code, creatorID).
		Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}
