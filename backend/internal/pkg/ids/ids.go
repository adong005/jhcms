package ids

import "github.com/google/uuid"

// DefaultTenantUUID 默认租户标识（用于系统模板数据与 JWT 兼容）。
const DefaultTenantUUID = "00000000-0000-0000-0000-000000000001"

// New 生成标准 UUID 字符串（36 字符含连字符），用于主键与外键。
func New() string {
	return uuid.New().String()
}

// Valid 是否为合法 UUID 字符串。
func Valid(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
