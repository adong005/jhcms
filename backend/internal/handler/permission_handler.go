package handler

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/common"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionRepo *repository.PermissionRepository
}

func NewPermissionHandler(permissionRepo *repository.PermissionRepository) *PermissionHandler {
	return &PermissionHandler{permissionRepo: permissionRepo}
}

func (h *PermissionHandler) GetPermissionList(c *gin.Context) {
	var req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Name     string `json:"name"`
		Code     string `json:"code"`
		Module   string `json:"module"`
	}

	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		permissions, total, err := h.permissionRepo.List(req.Page, req.PageSize, req.Name, req.Code, req.Module)
		if err != nil {
			return nil, 0, err
		}

		items := make([]map[string]interface{}, 0, len(permissions))
		for _, p := range permissions {
			items = append(items, map[string]interface{}{
				"id":          p.ID,
				"name":        p.Name,
				"code":        p.Code,
				"module":      p.Module,
				"isDelegable": p.IsDelegable,
				"createTime":  p.CreatedAt.Format("2006-01-02 15:04:05"),
				"updateTime":  p.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		return items, total, nil
	}, "获取权限列表失败")
}

func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Module      string `json:"module"`
		IsDelegable *bool  `json:"isDelegable"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	permission := &model.Permission{
		Name:        req.Name,
		Code:        req.Code,
		Module:      req.Module,
		IsDelegable: true,
	}
	if req.IsDelegable != nil {
		permission.IsDelegable = *req.IsDelegable
	}

	if err := h.permissionRepo.Create(permission); err != nil {
		response.Error(c, "创建权限失败")
		return
	}
	response.SuccessWithMessage(c, "创建权限成功", nil)
}

func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	var req struct {
		ID          string `json:"id" binding:"required"`
		Name        string `json:"name"`
		Code        string `json:"code"`
		Module      string `json:"module"`
		IsDelegable *bool  `json:"isDelegable"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的权限ID")
		return
	}

	permission, err := h.permissionRepo.GetByID(req.ID)
	if err != nil {
		response.Error(c, "权限不存在")
		return
	}
	if req.Name != "" {
		permission.Name = req.Name
	}
	if req.Code != "" {
		permission.Code = req.Code
	}
	if req.Module != "" {
		permission.Module = req.Module
	}
	if req.IsDelegable != nil {
		permission.IsDelegable = *req.IsDelegable
	}

	if err = h.permissionRepo.Update(permission); err != nil {
		response.Error(c, "更新权限失败")
		return
	}
	response.SuccessWithMessage(c, "更新权限成功", nil)
}

func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的权限ID")
		return
	}

	if err := h.permissionRepo.Delete(req.ID); err != nil {
		response.Error(c, "删除权限失败")
		return
	}
	response.SuccessWithMessage(c, "删除权限成功", nil)
}

func (h *PermissionHandler) BatchDeletePermission(c *gin.Context) {
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
		response.Error(c, "缺少有效权限ID")
		return
	}

	if err := h.permissionRepo.BatchDelete(validIDs); err != nil {
		response.Error(c, "批量删除权限失败")
		return
	}
	response.SuccessWithMessage(c, "批量删除权限成功", nil)
}
