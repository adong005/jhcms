package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SuccessCode = 0
	ErrorCode   = 1
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PageResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（带自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Message: message,
		Data:    data,
	})
}

// PageSuccess 分页成功响应
func PageSuccess(c *gin.Context, items interface{}, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:    SuccessCode,
		Message: "success",
		Data: map[string]interface{}{
			"items": items,
			"total": total,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    ErrorCode,
		Message: message,
	})
}

// ErrorWithCode 错误响应（自定义错误码）
func ErrorWithCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: message,
	})
}
