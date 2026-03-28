package repository

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type SiteGroupRepository struct {
	db *gorm.DB
}

func NewSiteGroupRepository(db *gorm.DB) *SiteGroupRepository {
	return &SiteGroupRepository{db: db}
}

func (r *SiteGroupRepository) List(tenantID, role, currentUserID, adminID string, page, pageSize int, keyword, subdomain string) ([]model.SiteGroup, int64, error) {
	var items []model.SiteGroup
	var total int64

	query := r.db.Model(&model.SiteGroup{})
	if role == "super_admin" {
		if adminID != "" {
			query = query.Where("tenant_id = ?", adminID)
		}
	} else if currentUserID != "" {
		query = query.Where("created_by = ? OR created_by IN (?)",
			currentUserID,
			r.db.Model(&model.User{}).Select("id").Where("created_by = ?", currentUserID),
		)
	} else if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if keyword != "" {
		query = query.Where("keyword LIKE ?", "%"+keyword+"%")
	}
	if subdomain != "" {
		query = query.Where("subdomain LIKE ?", "%"+subdomain+"%")
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

type AdminOption struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	NickName string `json:"nickName"`
}

func (r *SiteGroupRepository) ListAdminOptions(role, currentUserID string) ([]AdminOption, error) {
	var rows []AdminOption
	q := r.db.Table("users").
		Select("id AS user_id, username, nick_name").
		Where("role = ?", "admin").
		Order("created_at ASC")
	if role == "admin" && currentUserID != "" {
		q = q.Where("id = ?", currentUserID)
	}
	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *SiteGroupRepository) AdminLabelByTenantIDs(tenantIDs []string) map[string]string {
	out := make(map[string]string)
	if len(tenantIDs) == 0 {
		return out
	}
	var rows []struct {
		ID       string `gorm:"column:id"`
		Username string `gorm:"column:username"`
		NickName string `gorm:"column:nick_name"`
	}
	if err := r.db.Table("users").Select("id, username, nick_name").Where("id IN ?", tenantIDs).Scan(&rows).Error; err != nil {
		return out
	}
	for _, r := range rows {
		label := strings.TrimSpace(r.NickName)
		if label == "" {
			label = strings.TrimSpace(r.Username)
		}
		out[r.ID] = label
	}
	return out
}

func (r *SiteGroupRepository) Create(item *model.SiteGroup) error {
	if item.ID == "" {
		item.ID = ids.New()
	}
	return r.db.Create(item).Error
}
func (r *SiteGroupRepository) Update(item *model.SiteGroup) error { return r.db.Save(item).Error }
func (r *SiteGroupRepository) GetByID(id string, tenantID string) (*model.SiteGroup, error) {
	var item model.SiteGroup
	query := r.db.Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if err := query.First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
func (r *SiteGroupRepository) Delete(id string, tenantID string) error {
	query := r.db.Where("id = ?", id)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Delete(&model.SiteGroup{}).Error
}
func (r *SiteGroupRepository) BatchDelete(ids []string, tenantID string) error {
	if len(ids) == 0 {
		return nil
	}
	query := r.db.Where("id IN ?", ids)
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return query.Delete(&model.SiteGroup{}).Error
}

type CityItem struct {
	CityCode int    `json:"cityCode"`
	Name     string `json:"name"`
	Pinyin   string `json:"pinyin"`
}

func normalizeDomain(domain string) string {
	domain = strings.TrimSpace(strings.ToLower(domain))
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimSuffix(domain, "/")
	return domain
}

func (r *SiteGroupRepository) GetAdminSiteDomain(tenantID string) (string, error) {
	if tenantID == "" {
		return "", nil
	}
	var row struct {
		Domain string `gorm:"column:domain"`
	}
	if err := r.db.Table("users").Select("domain").Where("id = ?", tenantID).First(&row).Error; err != nil {
		return "", err
	}
	return normalizeDomain(row.Domain), nil
}

func (r *SiteGroupRepository) ListAreaCities() ([]CityItem, error) {
	var rows []model.City
	if err := r.db.Model(&model.City{}).
		Where("status = ?", 1).
		Order("city_code ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]CityItem, 0, len(rows))
	for _, row := range rows {
		name := strings.TrimSpace(row.Name)
		py := strings.TrimSpace(row.Pinyin)
		if row.CityCode <= 0 || name == "" || py == "" {
			continue
		}
		items = append(items, CityItem{
			CityCode: row.CityCode,
			Name:     name,
			Pinyin:   py,
		})
	}
	return items, nil
}

// BuildVirtualCitySiteGroups 运行时根据城市列表构造默认站群（不落库）。
func (r *SiteGroupRepository) BuildVirtualCitySiteGroups(tenantID, createdBy, domain, keyword, subdomain string, page, pageSize int) ([]model.SiteGroup, int64, error) {
	cities, err := r.ListAreaCities()
	if err != nil {
		return nil, 0, err
	}
	domain = normalizeDomain(domain)
	kw := strings.ToLower(strings.TrimSpace(keyword))
	sub := strings.ToLower(strings.TrimSpace(subdomain))
	var createdByPtr *string
	if strings.TrimSpace(createdBy) != "" {
		createdByPtr = &createdBy
	}

	filtered := make([]model.SiteGroup, 0, len(cities))
	for _, c := range cities {
		fullSubdomain := c.Pinyin
		if domain != "" {
			fullSubdomain = fmt.Sprintf("%s.%s", c.Pinyin, domain)
		}
		if kw != "" && !strings.Contains(strings.ToLower(c.Name), kw) {
			continue
		}
		if sub != "" && !strings.Contains(strings.ToLower(fullSubdomain), sub) {
			continue
		}
		filtered = append(filtered, model.SiteGroup{
			ID:           "virtual-" + c.Pinyin,
			TenantScoped: model.TenantScoped{TenantID: tenantID},
			Keyword:      c.Name,
			Subdomain:    fullSubdomain,
			Title:        "",
			Keywords:     "",
			Description:  "",
			CreatorOptional: model.CreatorOptional{
				CreatedBy: createdByPtr,
			},
		})
	}

	total := int64(len(filtered))
	if pageSize <= 0 {
		return filtered, total, nil
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * pageSize
	if offset >= len(filtered) {
		return []model.SiteGroup{}, total, nil
	}
	end := offset + pageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	return filtered[offset:end], total, nil
}

// EnsureAdminDefaultCityGroups 管理员默认补齐城市站群数据。
func (r *SiteGroupRepository) EnsureAdminDefaultCityGroups(tenantID, createdBy, domain string) error {
	domain = normalizeDomain(domain)
	if tenantID == "" || domain == "" {
		return nil
	}

	cities, err := r.ListAreaCities()
	if err != nil {
		return err
	}
	if len(cities) == 0 {
		return nil
	}

	var existed []string
	if err := r.db.Model(&model.SiteGroup{}).
		Where("tenant_id = ?", tenantID).
		Pluck("subdomain", &existed).Error; err != nil {
		return err
	}
	existSet := make(map[string]struct{}, len(existed))
	for _, s := range existed {
		existSet[strings.ToLower(strings.TrimSpace(s))] = struct{}{}
	}

	var createdByPtr *string
	if strings.TrimSpace(createdBy) != "" {
		createdByPtr = &createdBy
	}

	toCreate := make([]model.SiteGroup, 0, len(cities))
	for _, c := range cities {
		keyword := strings.TrimSpace(c.Name)
		py := strings.TrimSpace(c.Pinyin)
		if keyword == "" || py == "" {
			continue
		}
		subdomain := fmt.Sprintf("%s.%s", py, domain)
		key := strings.ToLower(subdomain)
		if _, ok := existSet[key]; ok {
			continue
		}
		existSet[key] = struct{}{}
		toCreate = append(toCreate, model.SiteGroup{
			ID:           ids.New(),
			TenantScoped: model.TenantScoped{TenantID: tenantID},
			Keyword:      keyword,
			Subdomain:    subdomain,
			Title:        "",
			Keywords:     "",
			Description:  "",
			CreatorOptional: model.CreatorOptional{
				CreatedBy: createdByPtr,
			},
		})
	}
	if len(toCreate) == 0 {
		return nil
	}
	return r.db.Create(&toCreate).Error
}
