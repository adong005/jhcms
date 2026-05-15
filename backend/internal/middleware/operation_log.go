package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	case strings.Contains(p, "/purge"), strings.Contains(p, "/clear"):
		return "delete"
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

func maskErrorMsg(raw string) string {
	if raw == "" {
		return ""
	}
	if len(raw) > 500 {
		raw = raw[:500] + "...(truncated)"
	}
	for _, kw := range []string{"password", "token", "secret", "authorization"} {
		raw = strings.ReplaceAll(strings.ToLower(raw), kw, "***")
	}
	return strings.TrimSpace(strings.ReplaceAll(raw, "\n", " | "))
}

// OperationLogMiddleware 将已登录用户的 API 操作落库到 system_logs（同步写入，避免通道丢数）。
func OperationLogMiddleware(db *gorm.DB, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := c.GetHeader("X-Request-Id")
		if reqID == "" || !ids.Valid(reqID) {
			reqID = ids.New()
		}
		c.Set("request_id", reqID)
		c.Header("X-Request-Id", reqID)

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

		if strings.HasPrefix(path, "/health") {
			return
		}

		userIDVal, hasUser := c.Get("user_id")
		usernameVal, hasUsername := c.Get("username")
		tenantIDVal, hasTenant := c.Get("tenant_id")
		if !hasUser || !hasUsername || !hasTenant {
			logger.Warn("operation log skipped: missing auth context",
				zap.Bool("hasUser", hasUser),
				zap.Bool("hasUsername", hasUsername),
				zap.Bool("hasTenant", hasTenant),
				zap.String("path", c.Request.URL.Path))
			return
		}

		userID, _ := userIDVal.(string)
		username, _ := usernameVal.(string)
		tenantID, _ := tenantIDVal.(string)
		if tenantID == "" {
			if rv, ok := c.Get("role"); ok {
				if r, _ := rv.(string); r == "super_admin" {
					tenantID = ids.DefaultTenantUUID
				}
			}
			if tenantID == "" {
				if sup, ok := c.Get("is_platform_super_admin"); ok && sup == true {
					tenantID = ids.DefaultTenantUUID
				}
			}
		}
		if userID == "" || username == "" || tenantID == "" {
			logger.Warn("operation log skipped: empty user identity",
				zap.String("userID", userID),
				zap.String("username", username),
				zap.String("tenantID", tenantID),
				zap.String("path", c.Request.URL.Path))
			return
		}

		userPath := ""
		if pv, ok := c.Get("user_path"); ok {
			userPath, _ = pv.(string)
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

		errorMsg := maskErrorMsg(c.Errors.String())
		if errorMsg == "" && status == "fail" {
			errorMsg = strconv.Itoa(c.Writer.Status())
		}

		entry := &model.SystemLog{
			ID:          ids.New(),
			RequestID:   reqID,
			TenantID:    tenantID,
			UserID:      userID,
			Path:        userPath,
			Username:    username,
			Action:      resolveAction(c.Request.Method, path),
			Module:      resolveModule(path),
			Description: c.Request.Method + " " + path,
			IP:          c.ClientIP(),
			Method:      c.Request.Method,
			URL:         path,
			UserAgent:   c.Request.UserAgent(),
			Status:      status,
			LogType:     "api",
			StatusCode:  c.Writer.Status(),
			Duration:    int(time.Since(start).Milliseconds()),
			ErrorMsg:    errorMsg,
			RequestJSON: requestPayload,
			ParentID:    parentID,
		}

		if err := db.Create(entry).Error; err != nil {
			logger.Error("operation log write failed",
				zap.Error(err),
				zap.String("logId", entry.ID),
				zap.String("path", path),
				zap.String("tenantID", entry.TenantID))
		}
	}
}
