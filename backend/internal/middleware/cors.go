package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware CORS 中间件
func CORSMiddleware(allowOrigins string) gin.HandlerFunc {
	origins := strings.Split(allowOrigins, ",")
	
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// 检查是否允许该来源
		allowed := false
		for _, o := range origins {
			if strings.TrimSpace(o) == origin || strings.TrimSpace(o) == "*" {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
