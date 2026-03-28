package handler

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/common"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/repository"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type SystemLogHandler struct {
	repo *repository.SystemLogRepository
}

func NewSystemLogHandler(repo *repository.SystemLogRepository) *SystemLogHandler {
	return &SystemLogHandler{repo: repo}
}

func (h *SystemLogHandler) appendSelfOperationLog(c *gin.Context, action, description string, payload interface{}) {
	tenantID, _ := c.Get("tenant_id")
	username, _ := c.Get("username")
	userID, _ := c.Get("user_id")
	tenantIDStr, _ := tenantID.(string)
	usernameStr, _ := username.(string)
	userIDStr, _ := userID.(string)
	if tenantIDStr == "" || usernameStr == "" || userIDStr == "" {
		return
	}
	var requestJSON string
	if payload != nil {
		if raw, err := json.Marshal(payload); err == nil {
			requestJSON = string(raw)
		}
	}
	var parentID *string
	if p, ok := c.Get("parent_id"); ok {
		if pid, ok2 := p.(string); ok2 && pid != "" {
			parentID = &pid
		}
	}
	_ = h.repo.Create(&model.SystemLog{
		ID:          ids.New(),
		TenantID:    tenantIDStr,
		Username:    usernameStr,
		Action:      action,
		Module:      "system-logs",
		Description: description,
		IP:          c.ClientIP(),
		Status:      "success",
		Duration:    0,
		RequestJSON: requestJSON,
		ParentID:    parentID,
		CreatedAt:   time.Now(),
	})
}

// GetSystemLogList 获取系统日志列表
func (h *SystemLogHandler) GetSystemLogList(c *gin.Context) {
	var req struct {
		Page      int    `json:"page"`
		PageSize  int    `json:"pageSize"`
		TenantID  string `json:"tenantId"`
		Username  string `json:"username"`
		Usernames []string `json:"usernames"`
		Action    string `json:"action"`
		Status    string `json:"status"`
	}

	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		tenantIDVal, _ := c.Get("tenant_id")
		roleVal, _ := c.Get("role")
		userIDVal, _ := c.Get("user_id")
		usernameVal, _ := c.Get("username")
		tenantID, _ := tenantIDVal.(string)
		role, _ := roleVal.(string)
		userID, _ := userIDVal.(string)
		currentUsername, _ := usernameVal.(string)
		items, total, err := h.repo.List(
			tenantID, role, userID, currentUsername,
			req.Page, req.PageSize, req.Username, req.Action, req.Status, req.TenantID, req.Usernames,
		)
		if err != nil {
			return nil, 0, err
		}

		result := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			result = append(result, map[string]interface{}{
				"id":          item.ID,
				"tenantId":    item.TenantID,
				"username":    item.Username,
				"action":      item.Action,
				"module":      item.Module,
				"description": item.Description,
				"ip":          item.IP,
				"status":      item.Status,
				"duration":    item.Duration,
				"errorMsg":    item.ErrorMsg,
				"requestJson": item.RequestJSON,
				"createTime":  item.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		return result, total, nil
	}, "获取系统日志列表失败")
}

func (h *SystemLogHandler) DeleteSystemLog(c *gin.Context) {
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
			response.Error(c, "无效日志ID")
			return
		}
		id = req.ID
	} else if !ids.Valid(id) {
		response.Error(c, "无效日志ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.repo.Delete(id, tenantID); err != nil {
		response.Error(c, "删除日志失败")
		return
	}
	h.appendSelfOperationLog(c, "delete", "DELETE /api/system-logs/:id", gin.H{"id": id})
	response.SuccessWithMessage(c, "删除日志成功", nil)
}

func (h *SystemLogHandler) BatchDeleteSystemLog(c *gin.Context) {
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
			response.Error(c, "缺少有效日志ID")
			return
		}
		tenantIDVal, _ := c.Get("tenant_id")
		tenantID, _ := tenantIDVal.(string)
		if err := h.repo.BatchDelete(valid, tenantID); err != nil {
			response.Error(c, "批量删除日志失败")
			return
		}
		h.appendSelfOperationLog(c, "delete", "POST /api/system-logs/batch-delete", gin.H{"ids": valid})
	}
	response.SuccessWithMessage(c, "批量删除日志成功", nil)
}

func (h *SystemLogHandler) ClearSystemLog(c *gin.Context) {
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.repo.ClearAll(tenantID); err != nil {
		response.Error(c, "清空日志失败")
		return
	}
	h.appendSelfOperationLog(c, "delete", "POST /api/system-logs/clear", gin.H{"scope": "tenant"})
	response.SuccessWithMessage(c, "清空日志成功", nil)
}
