package middleware

import (
	"adcms-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

var permissionChecker func(c *gin.Context, code string) bool

func SetPermissionChecker(checker func(c *gin.Context, code string) bool) {
	permissionChecker = checker
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
