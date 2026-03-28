package handler

import (
	"adcms-backend/internal/pkg/common"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/repository"
	"strings"

	"github.com/gin-gonic/gin"
)

type InfoHandler struct {
	infoRepo         *repository.InfoRepository
	userRepo         *repository.UserRepository
	infoCategoryRepo *repository.InfoCategoryRepository
}

func NewInfoHandler(
	infoRepo *repository.InfoRepository,
	userRepo *repository.UserRepository,
	infoCategoryRepo *repository.InfoCategoryRepository,
) *InfoHandler {
	return &InfoHandler{
		infoRepo:         infoRepo,
		userRepo:         userRepo,
		infoCategoryRepo: infoCategoryRepo,
	}
}

// GetInfoList 获取信息列表
func (h *InfoHandler) GetInfoList(c *gin.Context) {
	var req struct {
		Page     int                       `json:"page"`
		PageSize int                       `json:"pageSize"`
		Title    string                    `json:"title"`
		Status   common.OptionalListStatus `json:"status"`
	}

	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		tenantID, _ := c.Get("tenant_id")
		role, _ := c.Get("role")
		userID, _ := c.Get("user_id")

		tenantIDStr, _ := tenantID.(string)
		roleStr := role.(string)
		userIDStr, _ := userID.(string)

		infos, total, err := h.infoRepo.List(tenantIDStr, roleStr, userIDStr, req.Page, req.PageSize, req.Title, req.Status.Ptr())
		if err != nil {
			return nil, 0, err
		}

		catIDs := make([]string, 0, len(infos))
		for _, item := range infos {
			if item.CategoryID != nil && *item.CategoryID != "" {
				catIDs = append(catIDs, *item.CategoryID)
			}
		}
		catLabels, _ := h.infoCategoryRepo.MapDisplayLabelByIDs(tenantIDStr, catIDs)

		publisherIDs := make([]string, 0, len(infos))
		for _, item := range infos {
			if item.CreatedBy != nil && *item.CreatedBy != "" {
				publisherIDs = append(publisherIDs, *item.CreatedBy)
			}
		}
		publisherLabels := h.userRepo.NickNameByUserIDs(publisherIDs)

		tenantIDs := make([]string, 0, len(infos))
		for _, item := range infos {
			tenantIDs = append(tenantIDs, item.TenantID)
		}
		ascription := h.userRepo.AscriptionByTenantIDs(tenantIDs)

		result := make([]map[string]interface{}, 0, len(infos))
		for _, item := range infos {
			var categoryID interface{}
			if item.CategoryID != nil {
				categoryID = *item.CategoryID
			}
			categoryName := ""
			if item.CategoryID != nil {
				if lbl, ok := catLabels[*item.CategoryID]; ok {
					categoryName = lbl
				}
			}
			authorDisplay := strings.TrimSpace(item.Author)
			if item.CreatedBy != nil {
				if nick, ok := publisherLabels[*item.CreatedBy]; ok && nick != "" {
					authorDisplay = nick
				}
			}
			result = append(result, map[string]interface{}{
				"id":            item.ID,
				"title":         item.Title,
				"categoryId":    categoryID,
				"categoryName":  categoryName,
				"author":        authorDisplay,
				"tenantId":      item.TenantID,
				"ownerNickName": ascription[item.TenantID],
				"summary":       item.Summary,
				"content":       item.Content,
				"status":        item.Status,
				"viewCount":     0,
				"createTime":    item.CreatedAt.Format("2006-01-02 15:04:05"),
				"updateTime":    item.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		return result, total, nil
	}, "获取信息列表失败")
}

// GetInfo 获取信息详情
func (h *InfoHandler) GetInfo(c *gin.Context) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		response.Error(c, "无效的信息ID")
		return
	}

	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	info, err := h.infoRepo.GetByID(idStr, tenantID)
	if err != nil {
		response.Error(c, "获取信息详情失败")
		return
	}

	if info.CategoryID != nil {
		if cat, err := h.infoCategoryRepo.GetByID(*info.CategoryID, tenantID); err == nil {
			if cat.Code != "" {
				info.CategoryName = cat.Name + "（" + cat.Code + "）"
			} else {
				info.CategoryName = cat.Name
			}
		}
	}

	response.Success(c, info)
}

// CreateInfo 创建信息
func (h *InfoHandler) CreateInfo(c *gin.Context) {
	var req struct {
		Title       string  `json:"title" binding:"required"`
		Content     string  `json:"content" binding:"required"`
		CategoryID  *string `json:"categoryId"`
		Status      int8    `json:"status"`
		PublishDate string  `json:"publishDate"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	userID, _ := c.Get("user_id")
	userIDStr, _ := userID.(string)
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)

	err := h.infoRepo.Create(
		tenantID,
		userIDStr,
		req.Title,
		req.Content,
		req.CategoryID,
		req.Status,
		req.PublishDate,
	)

	if err != nil {
		response.Error(c, "创建信息失败")
		return
	}

	response.SuccessWithMessage(c, "创建信息成功", nil)
}

// UpdateInfo 更新信息
func (h *InfoHandler) UpdateInfo(c *gin.Context) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		response.Error(c, "无效的信息ID")
		return
	}

	var req struct {
		Title       string  `json:"title" binding:"required"`
		Content     string  `json:"content" binding:"required"`
		CategoryID  *string `json:"categoryId"`
		Status      int8    `json:"status"`
		PublishDate string  `json:"publishDate"`
		Author      string  `json:"author"`
		Summary     string  `json:"summary"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	err := h.infoRepo.Update(
		idStr,
		tenantID,
		req.Title,
		req.Content,
		req.CategoryID,
		req.Status,
		req.PublishDate,
		req.Author,
		req.Summary,
	)

	if err != nil {
		response.Error(c, "更新信息失败")
		return
	}

	response.SuccessWithMessage(c, "更新信息成功", nil)
}

// DeleteInfo 删除信息
func (h *InfoHandler) DeleteInfo(c *gin.Context) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		response.Error(c, "无效的信息ID")
		return
	}

	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	err := h.infoRepo.Delete(idStr, tenantID)
	if err != nil {
		response.Error(c, "删除信息失败")
		return
	}

	response.SuccessWithMessage(c, "删除信息成功", nil)
}

// GetInfoDetailByPath 兼容前端 GET /info/detail/:id
func (h *InfoHandler) GetInfoDetailByPath(c *gin.Context) {
	h.GetInfo(c)
}

// UpdateInfoByBody 兼容前端 POST /info/update
func (h *InfoHandler) UpdateInfoByBody(c *gin.Context) {
	var req struct {
		ID          string  `json:"id" binding:"required"`
		Title       string  `json:"title" binding:"required"`
		Content     string  `json:"content" binding:"required"`
		CategoryID  *string `json:"categoryId"`
		Status      int8    `json:"status"`
		PublishDate string  `json:"publishDate"`
		Author      string  `json:"author"`
		Summary     string  `json:"summary"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的信息ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.infoRepo.Update(req.ID, tenantID, req.Title, req.Content, req.CategoryID, req.Status, req.PublishDate, req.Author, req.Summary); err != nil {
		response.Error(c, "更新信息失败")
		return
	}
	response.SuccessWithMessage(c, "更新信息成功", nil)
}

// DeleteInfoByBody 兼容前端 POST /info/delete
func (h *InfoHandler) DeleteInfoByBody(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的信息ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.infoRepo.Delete(req.ID, tenantID); err != nil {
		response.Error(c, "删除信息失败")
		return
	}
	response.SuccessWithMessage(c, "删除信息成功", nil)
}

// UpdateInfoStatus 兼容前端 POST /info/status
func (h *InfoHandler) UpdateInfoStatus(c *gin.Context) {
	var req struct {
		ID     string `json:"id" binding:"required"`
		Status int8   `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的信息ID")
		return
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.infoRepo.UpdateStatus(req.ID, tenantID, req.Status); err != nil {
		response.Error(c, "更新状态失败")
		return
	}
	response.SuccessWithMessage(c, "更新状态成功", nil)
}

// BatchDeleteInfo 兼容前端 POST /info/batch-delete
func (h *InfoHandler) BatchDeleteInfo(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	for _, id := range req.IDs {
		if !ids.Valid(id) {
			response.Error(c, "无效的信息ID")
			return
		}
	}
	tenantIDVal, _ := c.Get("tenant_id")
	tenantID, _ := tenantIDVal.(string)
	if err := h.infoRepo.BatchDelete(req.IDs, tenantID); err != nil {
		response.Error(c, "批量删除信息失败")
		return
	}
	response.SuccessWithMessage(c, "批量删除信息成功", nil)
}
