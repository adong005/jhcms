package middleware

import (
	"adcms-backend/internal/pkg/jwt"
	"adcms-backend/internal/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Unauthorized(c, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("tenant_id", claims.TenantID)
		isAdmin := claims.IsAdmin || claims.LegacySuperAdmin
		c.Set("is_admin", isAdmin)
		c.Set("is_platform_super_admin", isAdmin)
		if claims.ParentID != nil {
			c.Set("parent_id", *claims.ParentID)
		}

		c.Next()
	}
}
