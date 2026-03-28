package middleware

import (
	"bytes"
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"
	"encoding/json"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func resolveAction(method, path string) string {
	p := strings.ToLower(path)
	switch {
	case strings.Contains(p, "/login"):
		return "login"
	case strings.Contains(p, "/logout"):
		return "logout"
	case strings.Contains(p, "/create"):
		return "create"
	case strings.Contains(p, "/update"), strings.Contains(p, "/status"):
		return "update"
	case strings.Contains(p, "/delete"), strings.Contains(p, "/batch-delete"):
		return "delete"
	case strings.Contains(p, "/export"):
		return "export"
	case strings.Contains(p, "/list"), method == "GET":
		return "query"
	default:
		switch strings.ToUpper(method) {
		case "POST":
			return "create"
		case "PUT", "PATCH":
			return "update"
		case "DELETE":
			return "delete"
		default:
			return "query"
		}
	}
}

func resolveModule(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 && parts[0] == "api" {
		return parts[1]
	}
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

func maskSensitivePayload(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	if len(s) > 8000 {
		s = s[:8000] + "...(truncated)"
	}

	var anyJSON interface{}
	if err := json.Unmarshal([]byte(s), &anyJSON); err != nil {
		return s
	}

	var walk func(v interface{}) interface{}
	walk = func(v interface{}) interface{} {
		switch x := v.(type) {
		case map[string]interface{}:
			out := make(map[string]interface{}, len(x))
			keys := make([]string, 0, len(x))
			for k := range x {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				low := strings.ToLower(k)
				if strings.Contains(low, "password") ||
					strings.Contains(low, "token") ||
					strings.Contains(low, "secret") ||
					strings.Contains(low, "authorization") ||
					strings.Contains(low, "verifycode") ||
					strings.Contains(low, "answer") {
					out[k] = "***"
					continue
				}
				out[k] = walk(x[k])
			}
			return out
		case []interface{}:
			out := make([]interface{}, 0, len(x))
			for _, it := range x {
				out = append(out, walk(it))
			}
			return out
		default:
			return v
		}
	}

	masked, _ := json.Marshal(walk(anyJSON))
	return string(masked)
}

// OperationLogMiddleware 将已登录用户的 API 操作落库到 system_logs。
func OperationLogMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestPayload := ""
		if c.Request != nil && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				_ = c.Request.Body.Close()
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				contentType := strings.ToLower(c.GetHeader("Content-Type"))
				if strings.Contains(contentType, "application/json") && len(bodyBytes) > 0 {
					requestPayload = maskSensitivePayload(string(bodyBytes))
				}
			}
		}
		c.Next()

		path := c.Request.URL.Path
		if !strings.HasPrefix(path, "/api/") {
			return
		}

		// 排除健康检查和日志自身写操作，避免噪音。
		if strings.HasPrefix(path, "/api/system-logs") || strings.HasPrefix(path, "/health") {
			return
		}

		userIDVal, hasUser := c.Get("user_id")
		usernameVal, hasUsername := c.Get("username")
		tenantIDVal, hasTenant := c.Get("tenant_id")
		if !hasUser || !hasUsername || !hasTenant {
			return
		}

		userID, _ := userIDVal.(string)
		username, _ := usernameVal.(string)
		tenantID, _ := tenantIDVal.(string)
		if userID == "" || username == "" || tenantID == "" {
			return
		}

		var parentID *string
		if pidVal, ok := c.Get("parent_id"); ok {
			if pid, ok2 := pidVal.(string); ok2 && pid != "" {
				parentID = &pid
			}
		}

		status := "success"
		if c.Writer.Status() >= 400 {
			status = "fail"
		}

		errorMsg := c.Errors.String()
		if errorMsg == "" && status == "fail" {
			errorMsg = strconv.Itoa(c.Writer.Status())
		}

		_ = db.Create(&model.SystemLog{
			ID:          ids.New(),
			TenantID:    tenantID,
			Username:    username,
			Action:      resolveAction(c.Request.Method, path),
			Module:      resolveModule(path),
			Description: c.Request.Method + " " + path,
			IP:          c.ClientIP(),
			Status:      status,
			Duration:    int(time.Since(start).Milliseconds()),
			ErrorMsg:    strings.TrimSpace(strings.ReplaceAll(errorMsg, "\n", " | ")),
			RequestJSON: requestPayload,
			ParentID:    parentID,
		}).Error
	}
}

