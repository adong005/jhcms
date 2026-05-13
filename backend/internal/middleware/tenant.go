package middleware

import (
	"strings"

	"adcms-backend/internal/pkg/ids"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var tenantDB *gorm.DB

// SetTenantDB 注入 DB，用于从数据库读取用户 path。
func SetTenantDB(db *gorm.DB) {
	tenantDB = db
}

// TenantMiddleware 租户隔离中间件（tenant_id 为 UUID 字符串），同时注入 user_path。
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

		// 注入 user_path，用于子树权限过滤。
		if tenantDB != nil {
			if userIDVal, ok := c.Get("user_id"); ok {
				userID, _ := userIDVal.(string)
				if userID != "" {
					var path string
					if err := tenantDB.Table("users").Select("path").Where("id = ?", userID).Scan(&path).Error; err == nil {
						c.Set("user_path", path)
						// 从 path 第一段重新推导 tenant_id（以 path 为准）。
						if path != "" && path != "/" {
							parts := strings.Split(strings.Trim(path, "/"), "/")
							if len(parts) > 0 && parts[0] != "" {
								c.Set("tenant_id", parts[0])
							}
						}
					}
				}
			}
		}

		c.Next()
	}
}
