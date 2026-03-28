package handler

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/common"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/repository"
	"strings"

	"github.com/gin-gonic/gin"
)

type SiteGroupHandler struct {
	repo *repository.SiteGroupRepository
}

func NewSiteGroupHandler(repo *repository.SiteGroupRepository) *SiteGroupHandler {
	return &SiteGroupHandler{repo: repo}
}

// GetSiteGroupList 获取站群列表
func (h *SiteGroupHandler) GetSiteGroupList(c *gin.Context) {
	var req struct {
		Page      int    `json:"page"`
		PageSize  int    `json:"pageSize"`
		Keyword   string `json:"keyword"`
		Subdomain string `json:"subdomain"`
		AdminID   string `json:"adminId"`
	}

	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		tenantIDVal, _ := c.Get("tenant_id")
		userIDVal, _ := c.Get("user_id")
		roleVal, _ := c.Get("role")
		tenantID, _ := tenantIDVal.(string)
		userID, _ := userIDVal.(string)
		role, _ := roleVal.(string)

		items, total, err := h.repo.List(tenantID, role, userID, req.AdminID, req.Page, req.PageSize, req.Keyword, req.Subdomain)
		if err != nil {
			return nil, 0, err
		}
		// 管理员首次进入站群管理时，如果还没有落库数据，则按 area_code 动态生成默认城市站群（不落库）。
		if role == "admin" && total == 0 {
			domain, _ := h.repo.GetAdminSiteDomain(tenantID)
			virtualItems, virtualTotal, vErr := h.repo.BuildVirtualCitySiteGroups(
				tenantID, userID, domain, req.Keyword, req.Subdomain, req.Page, req.PageSize,
			)
			if vErr == nil {
				items = virtualItems
				total = virtualTotal
			}
		}
		tenantIDs := make([]string, 0, len(items))
		for _, it := range items {
			tenantIDs = append(tenantIDs, it.TenantID)
		}
		adminLabels := h.repo.AdminLabelByTenantIDs(tenantIDs)

		result := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			createTime := ""
			updateTime := ""
			if !item.CreatedAt.IsZero() {
				createTime = item.CreatedAt.Format("2006-01-02 15:04:05")
			}
			if !item.UpdatedAt.IsZero() {
				updateTime = item.UpdatedAt.Format("2006-01-02 15:04:05")
			}
			result = append(result, map[string]interface{}{
				"id":          item.ID,
				"adminId":     item.TenantID,
				"adminName":   adminLabels[item.TenantID],
				"keyword":     item.Keyword,
				"subdomain":   item.Subdomain,
				"title":       item.Title,
				"keywords":    item.Keywords,
				"description": item.Description,
				"createTime":  createTime,
				"updateTime":  updateTime,
			})
		}

		return result, total, nil
	}, "获取站群列表失败")
}

// GetCityList 获取城市列表（仅超管使用）
func (h *SiteGroupHandler) GetCityList(c *gin.Context) {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "super_admin" {
		response.Forbidden(c, "无权访问城市列表")
		return
	}

	var req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Name     string `json:"name"`
	}
	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		items, err := h.repo.ListAreaCities()
		if err != nil {
			return nil, 0, err
		}
		keyword := req.Name
		if keyword != "" {
			filtered := make([]repository.CityItem, 0, len(items))
			for _, item := range items {
				if strings.Contains(item.Name, keyword) || strings.Contains(item.Pinyin, strings.ToLower(keyword)) {
					filtered = append(filtered, item)
				}
			}
			items = filtered
		}
		total := int64(len(items))
		if req.Page <= 0 {
			req.Page = 1
		}
		if req.PageSize <= 0 {
			req.PageSize = 10
		}
		offset := (req.Page - 1) * req.PageSize
		if offset >= len(items) {
			return []repository.CityItem{}, total, nil
		}
		end := offset + req.PageSize
		if end > len(items) {
			end = len(items)
		}
		return items[offset:end], total, nil
	}, "获取城市列表失败")
}

// GetAdminOptions 获取管理员下拉选项（超管可见全部，管理员仅自己）
func (h *SiteGroupHandler) GetAdminOptions(c *gin.Context) {
	roleVal, _ := c.Get("role")
	userIDVal, _ := c.Get("user_id")
	role, _ := roleVal.(string)
	userID, _ := userIDVal.(string)
	items, err := h.repo.ListAdminOptions(role, userID)
	if err != nil {
		response.Error(c, "获取管理员列表失败")
		return
	}
	response.Success(c, items)
}

func (h *SiteGroupHandler) CreateSiteGroup(c *gin.Context) {
	var req struct {
		Keyword     string `json:"keyword"`
		Subdomain   string `json:"subdomain" binding:"required"`
		Title       string `json:"title"`
		Keywords    string `json:"keywords"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	userID, _ := c.Get("user_id")
	var createdBy *string
	if uid, ok := userID.(string); ok && uid != "" {
		createdBy = &uid
	}
	data := &model.SiteGroup{
		TenantScoped:    model.TenantScoped{TenantID: ids.DefaultTenantUUID},
		Keyword:         req.Keyword,
		Subdomain:       req.Subdomain,
		Title:           req.Title,
		Keywords:        req.Keywords,
		Description:     req.Description,
		CreatorOptional: model.CreatorOptional{CreatedBy: createdBy},
	}
	if tenantIDVal, ok := c.Get("tenant_id"); ok {
		if tenantID, ok2 := tenantIDVal.(string); ok2 && tenantID != "" {
			data.TenantID = tenantID
		}
	}
	if err := h.repo.Create(data); err != nil {
		response.Error(c, "创建站群失败")
		return
	}
	response.SuccessWithMessage(c, "创建站群成功", nil)
}

func (h *SiteGroupHandler) UpdateSiteGroup(c *gin.Context) {
	var req struct {
		ID          string `json:"id"`
		Keyword     string `json:"keyword"`
		Subdomain   string `json:"subdomain"`
		Title       string `json:"title"`
		Keywords    string `json:"keywords"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if req.ID == "" {
		req.ID = c.Param("id")
	}
	if req.ID == "" || !ids.Valid(req.ID) {
		response.Error(c, "缺少或无效的站群ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	item, err := h.repo.GetByID(req.ID, tenantID)
	if err != nil {
		response.Error(c, "站群不存在")
		return
	}
	item.Keyword = req.Keyword
	item.Subdomain = req.Subdomain
	item.Title = req.Title
	item.Keywords = req.Keywords
	item.Description = req.Description
	if err := h.repo.Update(item); err != nil {
		response.Error(c, "更新站群失败")
		return
	}
	response.SuccessWithMessage(c, "更新站群成功", nil)
}

func (h *SiteGroupHandler) DeleteSiteGroup(c *gin.Context) {
	id := c.Param("id")
	if id != "" {
		if !ids.Valid(id) {
			response.Error(c, "无效的站群ID")
			return
		}
		tenantIDVal, _ := c.Get("tenant_id")
		tenantID, _ := tenantIDVal.(string)
		if err := h.repo.Delete(id, tenantID); err != nil {
			response.Error(c, "删除站群失败")
			return
		}
	}
	response.SuccessWithMessage(c, "删除站群成功", nil)
}

func (h *SiteGroupHandler) BatchDeleteSiteGroup(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	valid := make([]string, 0, len(req.IDs))
	for _, raw := range req.IDs {
		if ids.Valid(raw) {
			valid = append(valid, raw)
		}
	}
	if len(valid) > 0 {
		tenantIDVal, _ := c.Get("tenant_id")
		tenantID, _ := tenantIDVal.(string)
		if err := h.repo.BatchDelete(valid, tenantID); err != nil {
			response.Error(c, "批量删除站群失败")
			return
		}
	}
	response.SuccessWithMessage(c, "批量删除站群成功", nil)
}
