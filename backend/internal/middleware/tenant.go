package middleware

import (
	"adcms-backend/internal/pkg/ids"

	"github.com/gin-gonic/gin"
)

// TenantMiddleware 租户隔离中间件（tenant_id 为 UUID 字符串）
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, exists := c.Get("tenant_id")
		if !exists {
			c.Next()
			return
		}

		tenantIDStr, _ := tenantID.(string)
		if superAdmin, _ := c.Get("is_platform_super_admin"); superAdmin == true {
			if v := c.GetHeader("X-Tenant-Id"); v != "" && ids.Valid(v) {
				tenantIDStr = v
			}
		}
		c.Set("tenant_id", tenantIDStr)
		c.Next()
	}
}
