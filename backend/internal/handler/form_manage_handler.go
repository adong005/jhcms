package handler

import (
	"adcms-backend/internal/pkg/common"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/repository"
	"strings"

	"github.com/gin-gonic/gin"
)

type FormManageHandler struct {
	repo *repository.FormRepository
}

func NewFormManageHandler(repo *repository.FormRepository) *FormManageHandler {
	return &FormManageHandler{repo: repo}
}

// GetFormList 获取表单列表
func (h *FormManageHandler) GetFormList(c *gin.Context) {
	var req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"pageSize"`
		Contact  string `json:"contact"`
		Phone    string `json:"phone"`
		Company  string `json:"company"`
	}

	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		tenantIDVal, _ := c.Get("tenant_id")
		roleVal, _ := c.Get("role")
		currentUserIDVal, _ := c.Get("user_id")
		tenantID, _ := tenantIDVal.(string)
		role, _ := roleVal.(string)
		currentUserID, _ := currentUserIDVal.(string)
		items, total, err := h.repo.List(tenantID, role, currentUserID, req.Page, req.PageSize, req.Contact, req.Phone, req.Company)
		if err != nil {
			if strings.Contains(err.Error(), "doesn't exist") {
				return []map[string]interface{}{}, 0, nil
			}
			return nil, 0, err
		}

		result := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			row := map[string]interface{}{
				"id":           item.ID,
				"contact":      item.Contact,
				"phone":        item.Phone,
				"company":      item.Company,
				"ip":           item.IP,
				"handleStatus": item.HandleStatus,
				"remark":       item.Remark,
				"createTime":   item.CreatedAt.Format("2006-01-02 15:04:05"),
				"updateTime":   item.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
			if item.CreatedBy != nil {
				row["createdBy"] = *item.CreatedBy
			}
			result = append(result, row)
		}

		return result, total, nil
	}, "获取表单列表失败")
}

func (h *FormManageHandler) DeleteForm(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		var req struct {
			ID string `json:"id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, "请求参数错误")
			return
		}
		if !ids.Valid(req.ID) {
			response.Error(c, "无效的表单ID")
			return
		}
		id = req.ID
	} else if !ids.Valid(id) {
		response.Error(c, "无效的表单ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	currentUserIDVal, _ := c.Get("user_id")
	currentUserID, _ := currentUserIDVal.(string)
	if err := h.repo.Delete(id, tenantID, role, currentUserID); err != nil && !strings.Contains(err.Error(), "doesn't exist") {
		response.Error(c, "删除表单失败")
		return
	}
	response.SuccessWithMessage(c, "删除表单成功", nil)
}

func (h *FormManageHandler) BatchDeleteForm(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if len(req.IDs) > 0 {
		valid := make([]string, 0, len(req.IDs))
		for _, raw := range req.IDs {
			if ids.Valid(raw) {
				valid = append(valid, raw)
			}
		}
		if len(valid) == 0 {
			response.Error(c, "缺少有效的表单ID")
			return
		}
		tenantIDVal, _ := c.Get("tenant_id")
		tenantID, _ := tenantIDVal.(string)
		roleVal, _ := c.Get("role")
		role, _ := roleVal.(string)
		currentUserIDVal, _ := c.Get("user_id")
		currentUserID, _ := currentUserIDVal.(string)
		if err := h.repo.BatchDelete(valid, tenantID, role, currentUserID); err != nil && !strings.Contains(err.Error(), "doesn't exist") {
			response.Error(c, "批量删除表单失败")
			return
		}
	}
	response.SuccessWithMessage(c, "批量删除表单成功", nil)
}

func (h *FormManageHandler) ExportForm(c *gin.Context) {
	response.SuccessWithMessage(c, "导出成功", map[string]interface{}{})
}
