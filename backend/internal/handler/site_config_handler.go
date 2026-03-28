package handler

import (
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/repository"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type SiteConfigHandler struct {
	userRepo *repository.UserRepository
}

func NewSiteConfigHandler(userRepo *repository.UserRepository) *SiteConfigHandler {
	return &SiteConfigHandler{
		userRepo: userRepo,
	}
}

// GetSiteConfig 获取网站配置
func (h *SiteConfigHandler) GetSiteConfig(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	roleStr := role.(string)
	userIDVal, _ := c.Get("user_id")
	userIDStr, _ := userIDVal.(string)

	var config *model.User
	var err error

	switch roleStr {
	case "super_admin", "admin":
		config, err = h.userRepo.GetSiteConfigByUserID(userIDStr)
	case "user":
		tenantIDVal, _ := c.Get("tenant_id")
		tenantID, _ := tenantIDVal.(string)
		config, err = h.userRepo.GetSiteConfigByTenantAdmin(tenantID)
	default:
		response.Forbidden(c, "无权获取网站配置")
		return
	}

	if err != nil {
		response.Error(c, "获取网站配置失败")
		return
	}

	response.Success(c, gin.H{
		"title":          config.Title,
		"keywords":       config.Keywords,
		"description":    config.Description,
		"domain":         config.Domain,
		"logo":           config.Logo,
		"icpCode":        config.ICPCode,
		"contactPhone":   config.ContactPhone,
		"contactAddress": config.ContactAddress,
		"contactEmail":   config.ContactEmail,
	})
}

// UpdateSiteConfig 更新网站配置
func (h *SiteConfigHandler) UpdateSiteConfig(c *gin.Context) {
	var req struct {
		UserID         string `json:"userId"`
		Title          string `json:"title" binding:"required"`
		Keywords       string `json:"keywords" binding:"required"`
		Description    string `json:"description" binding:"required"`
		Domain         string `json:"domain" binding:"required"`
		Logo           string `json:"logo"`
		ICPCode        string `json:"icpCode"`
		ContactPhone   string `json:"contactPhone"`
		ContactAddress string `json:"contactAddress"`
		ContactEmail   string `json:"contactEmail"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	role, exists := c.Get("role")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	roleStr := role.(string)
	userIDVal, _ := c.Get("user_id")
	targetUserID, _ := userIDVal.(string)

	if roleStr != "admin" && roleStr != "super_admin" {
		response.Forbidden(c, "只有管理员可以更新网站配置")
		return
	}

	err := h.userRepo.UpdateSiteConfig(
		targetUserID,
		req.Title,
		req.Keywords,
		req.Description,
		req.Domain,
		req.Logo,
		req.ICPCode,
		req.ContactPhone,
		req.ContactAddress,
		req.ContactEmail,
	)

	if err != nil {
		response.Error(c, "更新网站配置失败")
		return
	}

	response.SuccessWithMessage(c, "网站配置更新成功", nil)
}

// UploadSiteLogo 上传网站 Logo，并把路径保存到 users.logo
func (h *SiteConfigHandler) UploadSiteLogo(c *gin.Context) {
	roleVal, exists := c.Get("role")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	role, _ := roleVal.(string)
	if role != "admin" && role != "super_admin" {
		response.Forbidden(c, "只有管理员可以上传网站Logo")
		return
	}
	userIDVal, _ := c.Get("user_id")
	userID, _ := userIDVal.(string)
	if userID == "" || !ids.Valid(userID) {
		response.Error(c, "无效用户")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, "请上传图片文件")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allow := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".svg": true,
	}
	if !allow[ext] {
		response.Error(c, "仅支持 jpg/jpeg/png/gif/webp/svg 格式")
		return
	}

	if err := os.MkdirAll("uploads/site-logo", 0o755); err != nil {
		response.Error(c, "创建上传目录失败")
		return
	}
	filename := fmt.Sprintf("%s_%s%s", userID, ids.New(), ext)
	savePath := filepath.Join("uploads/site-logo", filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.Error(c, "保存图片失败")
		return
	}

	urlPath := "/uploads/site-logo/" + filename
	if err := h.userRepo.UpdateSiteLogo(userID, urlPath); err != nil {
		response.Error(c, "保存Logo路径失败")
		return
	}

	response.Success(c, gin.H{
		"logo": urlPath,
	})
}
