package handler

import (
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Login 登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	result, err := h.authService.Login(&req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, result)
}

// Logout 登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 前端清除 Token 即可，后端无需处理
	response.SuccessWithMessage(c, "登出成功", nil)
}

// RefreshToken 刷新令牌
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	accessToken, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"accessToken": accessToken,
	})
}

// GetUserInfo 获取用户信息
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	uid, ok := userID.(string)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	userInfo, err := h.authService.GetUserInfo(uid)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, userInfo)
}

// GetAccessCodes 获取权限码
func (h *AuthHandler) GetAccessCodes(c *gin.Context) {
	role, exists := c.Get("role")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}

	uid, ok := userID.(string)
	if !ok {
		response.Unauthorized(c, "未登录")
		return
	}
	codes := h.authService.GetAccessCodes(uid, role.(string))
	response.Success(c, codes)
}

// Impersonate 模拟登录为指定用户（授权规则见 AuthService.canImpersonate，不仅依赖 JWT 里的 is_platform_super_admin）
func (h *AuthHandler) Impersonate(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return
	}
	actorID, ok := userIDVal.(string)
	if !ok || actorID == "" {
		response.Unauthorized(c, "未登录")
		return
	}

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	result, err := h.authService.Impersonate(actorID, req.UserID)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	response.Success(c, result)
}
