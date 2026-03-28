package handler

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/common"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/repository"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleRepo *repository.RoleRepository
}

func NewRoleHandler(roleRepo *repository.RoleRepository) *RoleHandler {
	return &RoleHandler{roleRepo: roleRepo}
}

// GetRoleList 获取角色列表
func (h *RoleHandler) GetRoleList(c *gin.Context) {
	var req struct {
		Page     int                       `json:"page"`
		PageSize int                       `json:"pageSize"`
		Name     string                    `json:"name"`
		Code     string                    `json:"code"`
		Status   common.OptionalListStatus `json:"status"`
	}

	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		tenantIDVal, _ := c.Get("tenant_id")
		operatorRoleVal, _ := c.Get("role")
		currentUserIDVal, _ := c.Get("user_id")
		tenantID, _ := tenantIDVal.(string)
		operatorRole, _ := operatorRoleVal.(string)
		currentUserID, _ := currentUserIDVal.(string)
		roles, total, err := h.roleRepo.List(tenantID, operatorRole, currentUserID, req.Page, req.PageSize, req.Name, req.Code, req.Status.Ptr())
		if err != nil {
			return nil, 0, err
		}
		items := make([]map[string]interface{}, 0, len(roles))
		for _, role := range roles {
			items = append(items, map[string]interface{}{
				"id":          role.ID,
				"name":        role.Name,
				"code":        role.Code,
				"description": role.Description,
				"status":      role.Status,
				"createTime":  role.CreatedAt.Format("2006-01-02 15:04:05"),
				"updateTime":  role.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		return items, total, nil
	}, "获取角色列表失败")
}

// CreateRole 创建角色（兼容前端 CRUD 路由）
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	code := strings.TrimSpace(req.Code)
	if code == "" {
		code = buildRoleCode(req.Name)
	}
	if operatorRoleVal, ok := c.Get("role"); ok {
		if operatorRole, ok2 := operatorRoleVal.(string); ok2 && operatorRole == "admin" {
			if code == "super_admin" || code == "admin" || code == "user" {
				response.Error(c, "租户管理员不可创建系统保留角色")
				return
			}
		}
	}
	role := &model.Role{
		TenantScoped: model.TenantScoped{TenantID: ids.DefaultTenantUUID},
		Name:         req.Name,
		Code:         code,
		Description:  req.Description,
		Status:       1,
	}
	if tenantIDVal, ok := c.Get("tenant_id"); ok {
		if tenantID, ok2 := tenantIDVal.(string); ok2 && tenantID != "" {
			role.TenantID = tenantID
		}
	}
	if req.Status != nil {
		role.Status = int8(*req.Status)
	}
	if userID, ok := c.Get("user_id"); ok {
		if uid, ok2 := userID.(string); ok2 && uid != "" {
			role.CreatedBy = &uid
		}
	}
	if err := h.roleRepo.Create(role); err != nil {
		response.Error(c, "创建角色失败")
		return
	}
	response.SuccessWithMessage(c, "创建角色成功", nil)
}

var nonRoleCodeRegexp = regexp.MustCompile(`[^a-z0-9_]+`)

func buildRoleCode(name string) string {
	base := strings.TrimSpace(strings.ToLower(name))
	base = strings.ReplaceAll(base, "-", "_")
	base = nonRoleCodeRegexp.ReplaceAllString(base, "_")
	base = strings.Trim(base, "_")
	if base == "" {
		base = "role"
	}
	if len(base) > 36 {
		base = base[:36]
	}
	return base + "_" + strings.ReplaceAll(ids.New(), "-", "")[:8]
}

// UpdateRole 更新角色（兼容前端 POST /role/update）
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	var req struct {
		ID          string `json:"id" binding:"required"`
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的角色ID")
		return
	}
	role, err := h.roleRepo.GetByID(req.ID)
	if err != nil {
		response.Error(c, "角色不存在")
		return
	}
	if req.Name != "" {
		role.Name = req.Name
	}
	if req.Code != "" {
		role.Code = req.Code
	}
	if req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != nil {
		role.Status = int8(*req.Status)
	}
	if err = h.roleRepo.Update(role); err != nil {
		response.Error(c, "更新角色失败")
		return
	}
	response.SuccessWithMessage(c, "更新角色成功", nil)
}

// UpdateRoleStatus 更新角色状态
func (h *RoleHandler) UpdateRoleStatus(c *gin.Context) {
	var req struct {
		ID     string `json:"id" binding:"required"`
		Status *int   `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	if !ids.Valid(req.ID) {
		response.Error(c, "无效的角色ID")
		return
	}
	role, err := h.roleRepo.GetByID(req.ID)
	if err != nil {
		response.Error(c, "角色不存在")
		return
	}
	if role.Code == "super_admin" && *req.Status == 0 {
		response.Error(c, "超级管理员角色不允许禁用")
		return
	}
	if err = h.roleRepo.UpdateStatus(req.ID, int8(*req.Status)); err != nil {
		response.Error(c, "更新角色状态失败")
		return
	}
	response.SuccessWithMessage(c, "更新角色状态成功", nil)
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	if !ids.Valid(req.ID) {
		response.Error(c, "无效的角色ID")
		return
	}
	if err := h.roleRepo.Delete(req.ID); err != nil {
		response.Error(c, "删除角色失败")
		return
	}
	response.SuccessWithMessage(c, "删除角色成功", nil)
}

// BatchDeleteRole 批量删除角色
func (h *RoleHandler) BatchDeleteRole(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	validIDs := make([]string, 0, len(req.IDs))
	for _, v := range req.IDs {
		if ids.Valid(v) {
			validIDs = append(validIDs, v)
		}
	}
	if len(validIDs) == 0 {
		response.Error(c, "缺少有效角色ID")
		return
	}
	if err := h.roleRepo.BatchDelete(validIDs); err != nil {
		response.Error(c, "批量删除角色失败")
		return
	}
	response.SuccessWithMessage(c, "批量删除角色成功", nil)
}

func (h *RoleHandler) GetRolePermission(c *gin.Context) {
	roleIDStr := c.Param("id")
	if !ids.Valid(roleIDStr) {
		response.Error(c, "无效的角色ID")
		return
	}
	permissionIDs, err := h.roleRepo.GetPermissionIDs(roleIDStr)
	if err != nil {
		response.Error(c, "获取角色权限失败")
		return
	}
	response.Success(c, gin.H{
		"permissionIds": permissionIDs,
		"menuIds":       permissionIDs, // 兼容旧前端字段
	})
}

func (h *RoleHandler) UpdateRolePermission(c *gin.Context) {
	var req struct {
		RoleID        string   `json:"roleId" binding:"required"`
		PermissionIDs []string `json:"permissionIds"`
		MenuIDs       []string `json:"menuIds"` // 兼容旧字段
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.RoleID) {
		response.Error(c, "无效的角色ID")
		return
	}
	rawPermissionIDs := req.PermissionIDs
	if len(rawPermissionIDs) == 0 && len(req.MenuIDs) > 0 {
		rawPermissionIDs = req.MenuIDs
	}
	permissionIDs := make([]string, 0, len(rawPermissionIDs))
	for _, v := range rawPermissionIDs {
		if ids.Valid(v) {
			permissionIDs = append(permissionIDs, v)
		}
	}

	if isSuper, _ := c.Get("is_platform_super_admin"); isSuper != true {
		roleCode, _ := c.Get("role")
		tenantIDVal, _ := c.Get("tenant_id")
		roleCodeStr, _ := roleCode.(string)
		tenantID, _ := tenantIDVal.(string)
		operatorRole, err := h.roleRepo.GetByCode(tenantID, roleCodeStr)
		if err != nil {
			response.Error(c, "获取操作者角色失败")
			return
		}
		assignable, err := h.roleRepo.GetDelegablePermissionIDsByRoleID(operatorRole.ID)
		if err != nil {
			response.Error(c, "获取可委派权限失败")
			return
		}
		allowMap := make(map[string]struct{}, len(assignable))
		for _, mid := range assignable {
			allowMap[mid] = struct{}{}
		}
		for _, targetID := range permissionIDs {
			if _, ok := allowMap[targetID]; !ok {
				response.Error(c, "存在超出当前角色可委派范围的权限")
				return
			}
		}
	}

	if err := h.roleRepo.SetPermissions(req.RoleID, permissionIDs); err != nil {
		response.Error(c, "更新角色权限失败")
		return
	}
	response.SuccessWithMessage(c, "更新角色权限成功", nil)
}
