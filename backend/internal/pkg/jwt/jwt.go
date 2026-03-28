package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID               string  `json:"user_id"`
	Username             string  `json:"username"`
	Role                 string  `json:"role"`
	TenantID             string  `json:"tenant_id"`
	IsAdmin              bool    `json:"is_admin"`
	LegacySuperAdmin     bool    `json:"is_platform_super_admin,omitempty"`
	ParentID             *string `json:"parent_id,omitempty"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// Init 初始化 JWT 密钥
func Init(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateToken 生成 Token（用户/租户主键为 UUID 字符串）
func GenerateToken(userID, username, role string, tenantID string, isAdmin bool, parentID *string, expireSeconds int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(expireSeconds) * time.Second)

	claims := Claims{
		UserID:               userID,
		Username:             username,
		Role:                 role,
		TenantID:             tenantID,
		IsAdmin:              isAdmin,
		LegacySuperAdmin:     isAdmin, // 兼容旧前端/旧 token 读取字段
		ParentID:             parentID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireTime),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			NotBefore: jwt.NewNumericDate(nowTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 Token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
