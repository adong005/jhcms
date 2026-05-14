package middleware

import (
	"adcms-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

var permissionChecker func(c *gin.Context, code string) bool

func SetPermissionChecker(checker func(c *gin.Context, code string) bool) {
	permissionChecker = checker
}

// CheckPermission 是否拥有指定权限码（不修改响应；供路由组合鉴权使用）。
func CheckPermission(c *gin.Context, code string) bool {
	if code == "" {
		return true
	}
	if permissionChecker == nil {
		return false
	}
	return permissionChecker(c, code)
}

func PermissionMiddleware(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if code == "" {
			c.Next()
			return
		}
		if permissionChecker == nil {
			// 未配置 checker 时默认拒绝，避免误放通。
			response.Forbidden(c, "权限校验未初始化")
			c.Abort()
			return
		}
		if !permissionChecker(c, code) {
			response.Forbidden(c, "无权限访问")
			c.Abort()
			return
		}
		c.Next()
	}
}

// AllPermissionMiddleware 需同时满足多个权限码（与 CheckPermission 组合使用）。
func AllPermissionMiddleware(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, code := range codes {
			if !CheckPermission(c, code) {
				response.Forbidden(c, "无权限访问")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
