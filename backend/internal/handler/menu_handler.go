package handler

import (
	"strings"

	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/common"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/permissionalias"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuRepo *repository.MenuRepository
	userRepo *repository.UserRepository
}

func NewMenuHandler(menuRepo *repository.MenuRepository, userRepo *repository.UserRepository) *MenuHandler {
	return &MenuHandler{
		menuRepo: menuRepo,
		userRepo: userRepo,
	}
}

// GetAllMenus 获取所有菜单（用于系统导航，树形结构）
func (h *MenuHandler) GetAllMenus(c *gin.Context) {
	common.HandleAllRequest(c, func() (interface{}, error) {
		tenantIDVal, _ := c.Get("tenant_id")
		tenantID, _ := tenantIDVal.(string)
		menus, err := h.menuRepo.GetAllMenus(tenantID)
		if err != nil {
			return nil, err
		}
		ids := make([]string, len(menus))
		for i := range menus {
			ids[i] = menus[i].ID
		}
		authMap, err := h.menuRepo.PermissionCodesByMenuIDs(ids)
		if err != nil {
			return nil, err
		}
		visible, err := h.filterNavMenusForRequest(c, menus, authMap)
		if err != nil {
			return nil, err
		}
		return h.buildMenuTree(visible, authMap), nil
	}, "获取菜单失败")
}

// GetMenuList 获取菜单列表（用于菜单管理，树形结构，带分页）
func (h *MenuHandler) GetMenuList(c *gin.Context) {
	var req struct {
		Page     int                       `json:"page"`
		PageSize int                       `json:"pageSize"`
		Name     string                    `json:"name"`
		Type     string                    `json:"type"`
		Status   common.OptionalListStatus `json:"status"`
	}

	common.HandleTreeListRequest(c, &req, 100, func() (interface{}, int64, error) {
		tenantIDVal, _ := c.Get("tenant_id")
		tenantID, _ := tenantIDVal.(string)
		menus, total, err := h.menuRepo.List(tenantID, req.Page, req.PageSize, req.Name, req.Type, req.Status.Ptr())
		if err != nil {
			return nil, 0, err
		}
		ids := make([]string, 0, len(menus))
		for i := range menus {
			ids = append(ids, menus[i].ID)
		}
		authMap, err := h.menuRepo.PermissionCodesByMenuIDs(ids)
		if err != nil {
			return nil, 0, err
		}
		items := h.convertMenusToListFormat(menus, authMap)
		return items, total, nil
	}, "获取菜单列表失败")
}

// buildMenuTree 构建菜单树形结构（用于系统导航）
func (h *MenuHandler) buildMenuTree(menus []model.Menu, authMap map[string][]string) []map[string]interface{} {
	menuMap := make(map[string]*model.Menu)
	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	var tree []map[string]interface{}

	for _, menu := range menus {
		if menu.ParentID == nil {
			menuItem := h.convertMenuToNavigationFormat(&menu, menuMap, authMap)
			tree = append(tree, menuItem)
		}
	}

	return tree
}

// convertMenuToNavigationFormat 转换菜单为导航格式
func (h *MenuHandler) convertMenuToNavigationFormat(menu *model.Menu, menuMap map[string]*model.Menu, authMap map[string][]string) map[string]interface{} {
	item := map[string]interface{}{
		"name": menu.Name,
		"path": menu.Path,
		"meta": map[string]interface{}{
			"title": menu.Name,
			"icon":  menu.Icon,
			"order": menu.SortOrder,
		},
	}

	if menu.Component != "" {
		item["component"] = menu.Component
	}
	// 兼容历史数据：个人中心若未配置组件，自动映射到内置页面
	if menu.Path == "/profile" && menu.Component == "" {
		item["component"] = "/_core/profile/index"
	}

	if auth := navAuthorityCodes(menu, authMap); len(auth) > 0 {
		item["meta"].(map[string]interface{})["authority"] = auth
	}

	var children []map[string]interface{}
	for _, m := range menuMap {
		if m.ParentID != nil && *m.ParentID == menu.ID {
			childItem := map[string]interface{}{
				"name":      m.Name,
				"path":      m.Path,
				"component": m.Component,
				"meta": map[string]interface{}{
					"title": m.Name,
					"icon":  m.Icon,
				},
			}
			if ca := navAuthorityCodes(m, authMap); len(ca) > 0 {
				childItem["meta"].(map[string]interface{})["authority"] = ca
			}
			children = append(children, childItem)
		}
	}

	if len(children) > 0 {
		item["children"] = children
	}

	return item
}

// convertMenusToListFormat 转换菜单列表为前端格式（用于菜单管理）
func (h *MenuHandler) convertMenusToListFormat(menus []model.Menu, authMap map[string][]string) []map[string]interface{} {
	var items []map[string]interface{}

	for _, menu := range menus {
		pc := strings.TrimSpace(menu.PermissionCode)
		comp := strings.TrimSpace(menu.Component)
		menuType := "menu"
		switch {
		case pc != "" && comp == "":
			menuType = "button"
		case comp == "":
			menuType = "catalog"
		default:
			menuType = "menu"
		}

		pcs := append([]string(nil), authMap[menu.ID]...)
		if pc != "" && !stringSliceContains(pcs, pc) {
			pcs = append(pcs, pc)
		}

		item := map[string]interface{}{
			"id":               menu.ID,
			"name":             menu.Name,
			"path":             menu.Path,
			"type":             menuType,
			"icon":             menu.Icon,
			"permissionCode":   menu.PermissionCode,
			"permissionCodes":  pcs,
			"component":        menu.Component,
			"parentId":         menu.ParentID,
			"order":            menu.SortOrder,
			"isShow":           menu.IsShow,
			"status":           menu.Status,
			"createTime":       menu.CreatedAt.Format("2006-01-02 15:04:05"),
			"updateTime":       menu.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		items = append(items, item)
	}

	return items
}

func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var req struct {
		Name            string      `json:"name" binding:"required"`
		Path            string      `json:"path"`
		Type            string      `json:"type"`
		Icon            string      `json:"icon"`
		PermissionCode  string      `json:"permissionCode"`
		PermissionCodes []string    `json:"permissionCodes"`
		Component       string      `json:"component"`
		ParentID        *string     `json:"parentId"`
		Order           interface{} `json:"order"`
		IsShow          *int8       `json:"isShow"`
		Status          *int8       `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	order, err := common.ParseOptionalInt(req.Order)
	if err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	linkCodes := mergeMenuPermissionCodesFromReq(req.PermissionCodes, req.PermissionCode)
	menu := &model.Menu{
		TenantScoped:   model.TenantScoped{TenantID: ids.DefaultTenantUUID},
		Name:           req.Name,
		Path:           req.Path,
		Icon:           req.Icon,
		PermissionCode: req.PermissionCode,
		Component:      req.Component,
		ParentID:       req.ParentID,
		SortOrder:      0,
		IsShow:         1,
		Status:         1,
	}
	if len(linkCodes) > 0 {
		menu.PermissionCode = linkCodes[0]
	}
	if order != nil {
		menu.SortOrder = *order
	}
	if tenantIDVal, ok := c.Get("tenant_id"); ok {
		if tenantID, ok2 := tenantIDVal.(string); ok2 && tenantID != "" {
			menu.TenantID = tenantID
		}
	}
	if req.Type == "catalog" {
		menu.Component = ""
	}
	if req.IsShow != nil {
		menu.IsShow = *req.IsShow
	}
	if req.Status != nil {
		menu.Status = *req.Status
	}
	if err := h.menuRepo.Create(menu); err != nil {
		response.Error(c, "创建菜单失败")
		return
	}
	if err := h.menuRepo.ReplaceMenuPermissionLinks(menu.ID, linkCodes); err != nil {
		response.Error(c, "创建菜单失败")
		return
	}
	response.SuccessWithMessage(c, "创建菜单成功", nil)
}

func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	var req struct {
		ID              string      `json:"id" binding:"required"`
		Name            string      `json:"name"`
		Path            string      `json:"path"`
		Type            string      `json:"type"`
		Icon            string      `json:"icon"`
		PermissionCode  string      `json:"permissionCode"`
		PermissionCodes []string    `json:"permissionCodes"`
		Component       string      `json:"component"`
		ParentID        *string     `json:"parentId"`
		Order           interface{} `json:"order"`
		IsShow          *int8       `json:"isShow"`
		Status          *int8       `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	order, err := common.ParseOptionalInt(req.Order)
	if err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的菜单ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	menu, err := h.menuRepo.GetMenuByID(req.ID, tenantID)
	if err != nil {
		response.Error(c, "菜单不存在")
		return
	}
	savedPermissionCode := menu.PermissionCode
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Path != "" {
		menu.Path = req.Path
	}
	if req.Icon != "" {
		menu.Icon = req.Icon
	}
	if req.Component != "" {
		menu.Component = req.Component
	}
	if req.Type == "catalog" {
		menu.Component = ""
	}
	if req.ParentID != nil {
		menu.ParentID = req.ParentID
	}
	if order != nil {
		menu.SortOrder = *order
	}
	if req.IsShow != nil {
		menu.IsShow = *req.IsShow
	}
	if req.Status != nil {
		menu.Status = *req.Status
	}
	linkCodes := mergeMenuPermissionCodesFromReq(req.PermissionCodes, req.PermissionCode)
	if len(linkCodes) == 0 {
		linkCodes = mergeMenuPermissionCodesFromReq(nil, savedPermissionCode)
	}
	if len(linkCodes) > 0 {
		menu.PermissionCode = linkCodes[0]
	}
	if err = h.menuRepo.Update(menu); err != nil {
		response.Error(c, "更新菜单失败")
		return
	}
	if err := h.menuRepo.ReplaceMenuPermissionLinks(menu.ID, linkCodes); err != nil {
		response.Error(c, "更新菜单失败")
		return
	}
	response.SuccessWithMessage(c, "更新菜单成功", nil)
}

func (h *MenuHandler) UpdateMenuStatus(c *gin.Context) {
	var req struct {
		ID     string `json:"id" binding:"required"`
		Status *int   `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的菜单ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.menuRepo.UpdateStatus(req.ID, tenantID, int8(*req.Status)); err != nil {
		response.Error(c, "更新菜单状态失败")
		return
	}
	response.SuccessWithMessage(c, "更新菜单状态成功", nil)
}

func (h *MenuHandler) UpdateMenuShow(c *gin.Context) {
	var req struct {
		ID     string `json:"id" binding:"required"`
		IsShow *int   `json:"isShow" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的菜单ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.menuRepo.UpdateShow(req.ID, tenantID, int8(*req.IsShow)); err != nil {
		response.Error(c, "更新菜单显示状态失败")
		return
	}
	response.SuccessWithMessage(c, "更新菜单显示状态成功", nil)
}

func (h *MenuHandler) DeleteMenu(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的菜单ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.menuRepo.Delete(req.ID, tenantID); err != nil {
		response.Error(c, "删除菜单失败")
		return
	}
	response.SuccessWithMessage(c, "删除菜单成功", nil)
}

func (h *MenuHandler) BatchDeleteMenu(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	valid := make([]string, 0, len(req.IDs))
	for _, idRaw := range req.IDs {
		if ids.Valid(idRaw) {
			valid = append(valid, idRaw)
		}
	}
	if len(valid) == 0 {
		response.Error(c, "缺少有效的菜单ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.menuRepo.BatchDelete(valid, tenantID); err != nil {
		response.Error(c, "批量删除菜单失败")
		return
	}
	response.SuccessWithMessage(c, "批量删除菜单成功", nil)
}

func (h *MenuHandler) filterNavMenusForRequest(c *gin.Context, menus []model.Menu, authMap map[string][]string) ([]model.Menu, error) {
	if v, ok := c.Get("is_platform_super_admin"); ok {
		if b, ok2 := v.(bool); ok2 && b {
			return menus, nil
		}
	}
	uidVal, _ := c.Get("user_id")
	roleVal, _ := c.Get("role")
	uid, _ := uidVal.(string)
	role, _ := roleVal.(string)
	if uid == "" {
		return menus, nil
	}
	rawCodes, err := h.userRepo.GetAccessCodesByUser(uid, role)
	if err != nil {
		return nil, err
	}
	eff := permissionalias.EffectiveCodes(rawCodes)

	idToMenu := make(map[string]model.Menu, len(menus))
	for i := range menus {
		idToMenu[menus[i].ID] = menus[i]
	}
	childrenByParent := make(map[string][]string)
	for i := range menus {
		m := menus[i]
		pk := ""
		if m.ParentID != nil && *m.ParentID != "" {
			pk = *m.ParentID
		}
		childrenByParent[pk] = append(childrenByParent[pk], m.ID)
	}

	memo := make(map[string]bool)
	var vis func(string) bool
	vis = func(menuID string) bool {
		if v, ok := memo[menuID]; ok {
			return v
		}
		m := idToMenu[menuID]
		req := navMenuRequiredCodes(m, authMap)
		var ok bool
		if len(req) > 0 {
			ok = navAnyCodeAllowed(eff, req)
		} else {
			kids := childrenByParent[menuID]
			if len(kids) == 0 {
				ok = true
			} else {
				for _, kid := range kids {
					if vis(kid) {
						ok = true
						break
					}
				}
			}
		}
		memo[menuID] = ok
		return ok
	}

	for i := range menus {
		vis(menus[i].ID)
	}

	out := make([]model.Menu, 0, len(menus))
	for _, m := range menus {
		if memo[m.ID] {
			out = append(out, m)
		}
	}
	return out, nil
}

func navMenuRequiredCodes(m model.Menu, authMap map[string][]string) []string {
	var req []string
	seen := map[string]struct{}{}
	for _, c := range authMap[m.ID] {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		req = append(req, c)
	}
	if pc := strings.TrimSpace(m.PermissionCode); pc != "" {
		if _, ok := seen[pc]; !ok {
			req = append(req, pc)
		}
	}
	return req
}

func navAnyCodeAllowed(eff map[string]struct{}, codes []string) bool {
	for _, c := range codes {
		if _, ok := eff[c]; ok {
			return true
		}
	}
	return false
}

func navAuthorityCodes(m *model.Menu, authMap map[string][]string) []string {
	return navMenuRequiredCodes(*m, authMap)
}

func stringSliceContains(list []string, v string) bool {
	for _, x := range list {
		if x == v {
			return true
		}
	}
	return false
}

func mergeMenuPermissionCodesFromReq(codes []string, legacy string) []string {
	var out []string
	seen := map[string]struct{}{}
	for _, c := range codes {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		if _, ok := seen[c]; ok {
			continue
		}
		seen[c] = struct{}{}
		out = append(out, c)
	}
	if x := strings.TrimSpace(legacy); x != "" {
		if _, ok := seen[x]; !ok {
			out = append(out, x)
		}
	}
	return out
}
