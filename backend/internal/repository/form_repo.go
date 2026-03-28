package repository

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type FormRepository struct {
	db *gorm.DB
}

func NewFormRepository(db *gorm.DB) *FormRepository {
	return &FormRepository{db: db}
}

var formTableNamePattern = regexp.MustCompile(`^[a-z0-9_]{1,64}$`)

func buildFormTableName(adminID string) (string, error) {
	if !ids.Valid(adminID) {
		return "", fmt.Errorf("invalid admin id")
	}
	table := "forms_u_" + strings.ToLower(strings.ReplaceAll(adminID, "-", ""))
	if !formTableNamePattern.MatchString(table) {
		return "", fmt.Errorf("invalid table name")
	}
	return table, nil
}

func (r *FormRepository) ensureFormTable(tableName string) error {
	if !formTableNamePattern.MatchString(tableName) {
		return fmt.Errorf("invalid table name")
	}
	hasBase := r.db.Migrator().HasTable("forms")
	if !hasBase {
		if err := r.db.AutoMigrate(&model.Form{}); err != nil {
			return err
		}
	}
	if r.db.Migrator().HasTable(tableName) {
		return nil
	}
	return r.db.Exec(fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` LIKE `forms`", tableName)).Error
}

func (r *FormRepository) listByTable(tableName string, page, pageSize int, contact, phone, company string) ([]model.Form, int64, error) {
	var items []model.Form
	var total int64
	query := r.db.Table(tableName)
	if contact != "" {
		query = query.Where("contact LIKE ?", "%"+contact+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if company != "" {
		query = query.Where("company LIKE ?", "%"+company+"%")
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

func (r *FormRepository) List(tenantID, role, currentUserID string, page, pageSize int, contact, phone, company string) ([]model.Form, int64, error) {
	if role != "super_admin" {
		adminID := tenantID
		if role == "admin" {
			adminID = currentUserID
		}
		tableName, err := buildFormTableName(adminID)
		if err != nil {
			return []model.Form{}, 0, nil
		}
		if err := r.ensureFormTable(tableName); err != nil {
			return nil, 0, err
		}
		return r.listByTable(tableName, page, pageSize, contact, phone, company)
	}

	// 超管聚合所有管理员分表数据，再做内存分页。
	var adminIDs []string
	if err := r.db.Model(&model.User{}).Where("role = ?", "admin").Pluck("id", &adminIDs).Error; err != nil {
		return nil, 0, err
	}
	all := make([]model.Form, 0)
	for _, adminID := range adminIDs {
		tableName, err := buildFormTableName(adminID)
		if err != nil {
			continue
		}
		if !r.db.Migrator().HasTable(tableName) {
			continue
		}
		rows, _, err := r.listByTable(tableName, 1, 1_000_000, contact, phone, company)
		if err != nil {
			return nil, 0, err
		}
		all = append(all, rows...)
	}

	sort.SliceStable(all, func(i, j int) bool {
		ai := all[i].CreatedAt
		aj := all[j].CreatedAt
		if ai.Equal(aj) {
			return all[i].ID > all[j].ID
		}
		if ai.IsZero() {
			return false
		}
		if aj.IsZero() {
			return true
		}
		return ai.After(aj)
	})

	total := int64(len(all))
	if total == 0 {
		return []model.Form{}, 0, nil
	}
	start := (page - 1) * pageSize
	if start >= len(all) {
		return []model.Form{}, total, nil
	}
	end := start + pageSize
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], total, nil
}

func (r *FormRepository) deleteByTable(tableName string, id string) error {
	return r.db.Table(tableName).Where("id = ?", id).Delete(&model.Form{}).Error
}

func (r *FormRepository) Delete(id string, tenantID, role, currentUserID string) error {
	if role != "super_admin" {
		adminID := tenantID
		if role == "admin" {
			adminID = currentUserID
		}
		tableName, err := buildFormTableName(adminID)
		if err != nil {
			return nil
		}
		if !r.db.Migrator().HasTable(tableName) {
			return nil
		}
		return r.deleteByTable(tableName, id)
	}
	// 超管优先按 tenantID 指定分表删除；未指定时扫描全部管理员分表。
	if ids.Valid(tenantID) {
		tableName, err := buildFormTableName(tenantID)
		if err == nil && r.db.Migrator().HasTable(tableName) {
			return r.deleteByTable(tableName, id)
		}
	}
	var adminIDs []string
	if err := r.db.Model(&model.User{}).Where("role = ?", "admin").Pluck("id", &adminIDs).Error; err != nil {
		return err
	}
	for _, adminID := range adminIDs {
		tableName, err := buildFormTableName(adminID)
		if err != nil || !r.db.Migrator().HasTable(tableName) {
			continue
		}
		if err := r.deleteByTable(tableName, id); err != nil {
			return err
		}
	}
	return nil
}

func (r *FormRepository) BatchDelete(formIDs []string, tenantID, role, currentUserID string) error {
	if len(formIDs) == 0 {
		return nil
	}
	if role != "super_admin" {
		adminID := tenantID
		if role == "admin" {
			adminID = currentUserID
		}
		tableName, err := buildFormTableName(adminID)
		if err != nil || !r.db.Migrator().HasTable(tableName) {
			return nil
		}
		return r.db.Table(tableName).Where("id IN ?", formIDs).Delete(&model.Form{}).Error
	}
	// 超管按指定租户优先；否则对全管理员分表执行批量删除。
	if ids.Valid(tenantID) {
		tableName, err := buildFormTableName(tenantID)
		if err == nil && r.db.Migrator().HasTable(tableName) {
			return r.db.Table(tableName).Where("id IN ?", formIDs).Delete(&model.Form{}).Error
		}
	}
	var adminIDs []string
	if err := r.db.Model(&model.User{}).Where("role = ?", "admin").Pluck("id", &adminIDs).Error; err != nil {
		return err
	}
	for _, adminID := range adminIDs {
		tableName, err := buildFormTableName(adminID)
		if err != nil || !r.db.Migrator().HasTable(tableName) {
			continue
		}
		if err := r.db.Table(tableName).Where("id IN ?", formIDs).Delete(&model.Form{}).Error; err != nil {
			return err
		}
	}
	return nil
}
