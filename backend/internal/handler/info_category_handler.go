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

type InfoCategoryHandler struct {
	repo     *repository.InfoCategoryRepository
	userRepo *repository.UserRepository
}

func NewInfoCategoryHandler(repo *repository.InfoCategoryRepository, userRepo *repository.UserRepository) *InfoCategoryHandler {
	return &InfoCategoryHandler{repo: repo, userRepo: userRepo}
}

// GetCategoryList 获取信息分类列表
func (h *InfoCategoryHandler) GetCategoryList(c *gin.Context) {
	var req struct {
		Page     int                   `json:"page"`
		PageSize int                   `json:"pageSize"`
		Name     string                `json:"name"`
		Status   common.OptionalListStatus `json:"status"`
	}

	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		tenantIDVal, _ := c.Get("tenant_id")
		roleVal, _ := c.Get("role")
		currentUserIDVal, _ := c.Get("user_id")
		tenantID, _ := tenantIDVal.(string)
		role, _ := roleVal.(string)
		currentUserID, _ := currentUserIDVal.(string)
		items, total, err := h.repo.List(tenantID, role, currentUserID, req.Page, req.PageSize, req.Name, req.Status.Ptr())
		if err != nil {
			return nil, 0, err
		}
		tenantIDs := make([]string, 0, len(items))
		for _, item := range items {
			tenantIDs = append(tenantIDs, item.TenantID)
		}
		ascription := h.userRepo.AscriptionByTenantIDs(tenantIDs)
		result := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			result = append(result, map[string]interface{}{
				"id":            item.ID,
				"name":          item.Name,
				"code":          item.Code,
				"isHome":        item.IsHome,
				"sort":          item.Sort,
				"description":   item.Description,
				"status":        item.Status,
				"tenantId":      item.TenantID,
				"ownerNickName": ascription[item.TenantID],
				"createTime":    item.CreatedAt.Format("2006-01-02 15:04:05"),
				"updateTime":    item.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		return result, total, nil
	}, "获取信息分类列表失败")
}

func (h *InfoCategoryHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code"`
		IsHome      *int   `json:"isHome"`
		Sort        int    `json:"sort"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	code := strings.TrimSpace(req.Code)
	if code == "" {
		code = "cat_" + strings.ReplaceAll(ids.New(), "-", "")
	}
	item := &model.InfoCategory{
		TenantScoped: model.TenantScoped{TenantID: ids.DefaultTenantUUID},
		Name:         req.Name,
		Code:         code,
		IsHome:       1,
		Sort:         req.Sort,
		Description:  req.Description,
		Status:       1,
	}
	if tenantIDVal, ok := c.Get("tenant_id"); ok {
		if tenantID, ok2 := tenantIDVal.(string); ok2 && tenantID != "" {
			item.TenantID = tenantID
		}
	}
	if req.Status != nil {
		item.Status = int8(*req.Status)
	}
	if req.IsHome != nil {
		item.IsHome = int8(*req.IsHome)
	}
	if userID, ok := c.Get("user_id"); ok {
		if uid, ok2 := userID.(string); ok2 && uid != "" {
			item.CreatedBy = &uid
		}
	}
	if err := h.repo.Create(item); err != nil {
		response.Error(c, "创建分类失败")
		return
	}
	response.SuccessWithMessage(c, "创建分类成功", nil)
}

func (h *InfoCategoryHandler) UpdateCategory(c *gin.Context) {
	var req struct {
		ID          string `json:"id" binding:"required"`
		Name        string `json:"name"`
		Code        string `json:"code"`
		IsHome      *int   `json:"isHome"`
		Sort        interface{} `json:"sort"`
		Description string `json:"description"`
		Status      *int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效分类ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	roleVal, _ := c.Get("role")
	tenantID, _ := tenantIDVal.(string)
	role, _ := roleVal.(string)
	if role == "super_admin" {
		tenantID = ""
	}
	item, err := h.repo.GetByID(req.ID, tenantID)
	if err != nil {
		response.Error(c, "分类不存在")
		return
	}
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Code != "" {
		item.Code = req.Code
	}
	if req.IsHome != nil {
		item.IsHome = int8(*req.IsHome)
	}
	sort, parseErr := common.ParseOptionalInt(req.Sort)
	if parseErr != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if sort != nil {
		item.Sort = *sort
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if req.Status != nil {
		item.Status = int8(*req.Status)
	}
	if err = h.repo.Update(item); err != nil {
		response.Error(c, "更新分类失败")
		return
	}
	response.SuccessWithMessage(c, "更新分类成功", nil)
}

func (h *InfoCategoryHandler) UpdateCategoryStatus(c *gin.Context) {
	var req struct {
		ID     string `json:"id" binding:"required"`
		Status int    `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效分类ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	roleVal, _ := c.Get("role")
	tenantID, _ := tenantIDVal.(string)
	role, _ := roleVal.(string)
	if role == "super_admin" {
		tenantID = ""
	}
	if err := h.repo.UpdateStatus(req.ID, tenantID, int8(req.Status)); err != nil {
		response.Error(c, "更新分类状态失败")
		return
	}
	response.SuccessWithMessage(c, "更新分类状态成功", nil)
}

func (h *InfoCategoryHandler) DeleteCategory(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效分类ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	roleVal, _ := c.Get("role")
	tenantID, _ := tenantIDVal.(string)
	role, _ := roleVal.(string)
	if role == "super_admin" {
		tenantID = ""
	}
	if err := h.repo.Delete(req.ID, tenantID); err != nil {
		response.Error(c, "删除分类失败")
		return
	}
	response.SuccessWithMessage(c, "删除分类成功", nil)
}

func (h *InfoCategoryHandler) BatchDeleteCategory(c *gin.Context) {
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
		response.Error(c, "缺少有效分类ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	roleVal, _ := c.Get("role")
	tenantID, _ := tenantIDVal.(string)
	role, _ := roleVal.(string)
	if role == "super_admin" {
		tenantID = ""
	}
	if err := h.repo.BatchDelete(validIDs, tenantID); err != nil {
		response.Error(c, "批量删除分类失败")
		return
	}
	response.SuccessWithMessage(c, "批量删除分类成功", nil)
}
