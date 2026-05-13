package handler

import (
	"encoding/csv"
	"strconv"
	"time"

	"adcms-backend/internal/pkg/common"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type SystemLogHandler struct {
	repo *repository.SystemLogRepository
}

func NewSystemLogHandler(repo *repository.SystemLogRepository) *SystemLogHandler {
	return &SystemLogHandler{repo: repo}
}

func tenantIDFromCtx(c *gin.Context) string {
	v, _ := c.Get("tenant_id")
	s, _ := v.(string)
	return s
}

func buildListParams(c *gin.Context, req *logListReq) repository.SystemLogListParams {
	tenantIDVal, _ := c.Get("tenant_id")
	roleVal, _ := c.Get("role")
	userIDVal, _ := c.Get("user_id")
	userPathVal, _ := c.Get("user_path")
	tenantID, _ := tenantIDVal.(string)
	role, _ := roleVal.(string)
	userID, _ := userIDVal.(string)
	userPath, _ := userPathVal.(string)

	var start, end *time.Time
	if req.StartTime != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, time.Local); err == nil {
			start = &t
		} else if t, err := time.ParseInLocation("2006-01-02", req.StartTime, time.Local); err == nil {
			start = &t
		}
	}
	if req.EndTime != "" {
		if t, err := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, time.Local); err == nil {
			end = &t
		} else if t, err := time.ParseInLocation("2006-01-02", req.EndTime, time.Local); err == nil {
			d := t.Add(24*time.Hour - time.Second)
			end = &d
		}
	}

	return repository.SystemLogListParams{
		TenantID:     tenantID,
		Role:         role,
		UserID:       userID,
		UserPath:     userPath,
		Page:         req.Page,
		PageSize:     req.PageSize,
		UsernameKw:   req.Username,
		Usernames:    req.Usernames,
		Action:       req.Action,
		Status:       req.Status,
		Module:       req.Module,
		IP:           req.IP,
		LogType:      req.LogType,
		StartTime:    start,
		EndTime:      end,
		TenantFilter: req.TenantID,
	}
}

type logListReq struct {
	Page      int      `json:"page"`
	PageSize  int      `json:"pageSize"`
	TenantID  string   `json:"tenantId"`
	Username  string   `json:"username"`
	Usernames []string `json:"usernames"`
	Action    string   `json:"action"`
	Status    string   `json:"status"`
	Module    string   `json:"module"`
	IP        string   `json:"ip"`
	LogType   string   `json:"logType"`
	StartTime string   `json:"startTime"`
	EndTime   string   `json:"endTime"`
}

// GetSystemLogList 获取系统日志列表
func (h *SystemLogHandler) GetSystemLogList(c *gin.Context) {
	var req logListReq
	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		params := buildListParams(c, &req)
		items, total, err := h.repo.List(params)
		if err != nil {
			return nil, 0, err
		}
		result := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			result = append(result, map[string]interface{}{
				"id":          item.ID,
				"requestId":   item.RequestID,
				"tenantId":    item.TenantID,
				"userId":      item.UserID,
				"username":    item.Username,
				"action":      item.Action,
				"module":      item.Module,
				"description": item.Description,
				"targetId":    item.TargetID,
				"ip":          item.IP,
				"method":      item.Method,
				"url":         item.URL,
				"status":      item.Status,
				"logType":     item.LogType,
				"statusCode":  item.StatusCode,
				"duration":    item.Duration,
				"errorMsg":    item.ErrorMsg,
				"requestJson": item.RequestJSON,
				"createTime":  item.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		return result, total, nil
	}, "获取系统日志列表失败")
}

// ExportSystemLog 导出日志为 CSV（最多 100000 条）
func (h *SystemLogHandler) ExportSystemLog(c *gin.Context) {
	var req logListReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	req.Page = 1
	req.PageSize = 100000
	params := buildListParams(c, &req)
	items, _, err := h.repo.List(params)
	if err != nil {
		response.Error(c, "导出日志失败")
		return
	}

	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="system_logs.csv"`)
	w := csv.NewWriter(c.Writer)
	_ = w.Write([]string{"ID", "操作用户", "操作类型", "模块", "描述", "IP", "状态", "耗时(ms)", "操作时间", "错误信息"})
	for _, item := range items {
		_ = w.Write([]string{
			item.ID,
			item.Username,
			item.Action,
			item.Module,
			item.Description,
			item.IP,
			item.Status,
			strconv.Itoa(item.Duration),
			item.CreatedAt.Format("2006-01-02 15:04:05"),
			item.ErrorMsg,
		})
	}
	w.Flush()
}

// PurgeOldLogs 清理 N 天前的日志
func (h *SystemLogHandler) PurgeOldLogs(c *gin.Context) {
	var req struct {
		Days int `json:"days" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误，days 必须为正整数")
		return
	}
	tenantID := tenantIDFromCtx(c)
	before := time.Now().AddDate(0, 0, -req.Days)
	affected, err := h.repo.PurgeBefore(tenantID, before)
	if err != nil {
		response.Error(c, "清理日志失败")
		return
	}
	response.SuccessWithMessage(c, "清理完成", gin.H{"deleted": affected})
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
	tenantID := tenantIDFromCtx(c)
	if err := h.repo.Delete(id, tenantID); err != nil {
		response.Error(c, "删除日志失败")
		return
	}
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
	tenantID := tenantIDFromCtx(c)
	if err := h.repo.BatchDelete(valid, tenantID); err != nil {
		response.Error(c, "批量删除日志失败")
		return
	}
	response.SuccessWithMessage(c, "批量删除日志成功", nil)
}

func (h *SystemLogHandler) ClearSystemLog(c *gin.Context) {
	tenantID := tenantIDFromCtx(c)
	// super_admin 清全平台需显式传 force=true
	var req struct {
		Force bool `json:"force"`
	}
	_ = c.ShouldBindJSON(&req)

	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role == "super_admin" && req.Force {
		tenantID = "__all__"
	}
	if err := h.repo.ClearAll(tenantID); err != nil {
		response.Error(c, "清空日志失败")
		return
	}
	response.SuccessWithMessage(c, "清空日志成功", nil)
}
