package handler

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func jsonStringField(v interface{}) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return strings.TrimSpace(x)
	case json.Number:
		return x.String()
	case float64:
		if x == float64(int64(x)) {
			return strconv.FormatInt(int64(x), 10)
		}
		return strings.TrimSpace(fmt.Sprint(x))
	default:
		return strings.TrimSpace(fmt.Sprint(x))
	}
}

func coerceInt8FromJSON(v interface{}) (int8, error) {
	switch x := v.(type) {
	case float64:
		return int8(x), nil
	case json.Number:
		n, err := x.Int64()
		if err != nil {
			return 0, err
		}
		return int8(n), nil
	case string:
		n, err := strconv.ParseInt(strings.TrimSpace(x), 10, 8)
		if err != nil {
			return 0, err
		}
		return int8(n), nil
	default:
		return 0, fmt.Errorf("unsupported status type %T", x)
	}
}

func coerceRoleFromInterface(v interface{}) (string, error) {
	if v == nil {
		return "", errors.New("missing role")
	}
	switch x := v.(type) {
	case string:
		s := strings.TrimSpace(x)
		switch s {
		case "super_admin", "admin", "user":
			return s, nil
		}
		if isRoleCodeLike(s) {
			return s, nil
		}
		if n, err := strconv.Atoi(s); err == nil {
			return roleCodeFromInt(n)
		}
		return "", fmt.Errorf("unknown role %q", s)
	case json.Number:
		n, err := x.Int64()
		if err != nil {
			return "", err
		}
		return roleCodeFromInt(int(n))
	case float64:
		return roleCodeFromInt(int(x))
	default:
		return "", fmt.Errorf("unsupported role type %T", x)
	}
}

// parseUserRoleJSON 兼容前端数字 role：1=超级管理员 2=管理员 3=用户；也接受字符串 super_admin/admin/user
func parseUserRoleJSON(raw json.RawMessage) (string, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return "", errors.New("empty role")
	}
	if raw[0] == '"' {
		var s string
		if err := json.Unmarshal(raw, &s); err != nil {
			return "", err
		}
		s = strings.TrimSpace(s)
		switch s {
		case "super_admin", "admin", "user":
			return s, nil
		}
		if isRoleCodeLike(s) {
			return s, nil
		}
		if n, err := strconv.Atoi(s); err == nil {
			return roleCodeFromInt(n)
		}
		return "", fmt.Errorf("unknown role %q", s)
	}
	var num json.Number
	if err := json.Unmarshal(raw, &num); err == nil {
		n64, err := num.Int64()
		if err != nil {
			return "", err
		}
		return roleCodeFromInt(int(n64))
	}
	var f float64
	if err := json.Unmarshal(raw, &f); err != nil {
		return "", err
	}
	return roleCodeFromInt(int(f))
}

var roleCodeRegexp = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_:-]{1,49}$`)

func isRoleCodeLike(s string) bool {
	return roleCodeRegexp.MatchString(strings.TrimSpace(s))
}

func roleCodeFromInt(n int) (string, error) {
	switch n {
	case 1:
		return "super_admin", nil
	case 2:
		return "admin", nil
	case 3:
		return "user", nil
	default:
		return "", fmt.Errorf("invalid role %d", n)
	}
}

func dataScopeForRole(role string) string {
	switch role {
	case "super_admin", "admin":
		return "TENANT_ALL"
	default:
		return "SELF"
	}
}

const tempPasswordChars = "abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func randomTempPassword(n int) (string, error) {
	b := make([]byte, n)
	for i := range b {
		var rb [1]byte
		if _, err := rand.Read(rb[:]); err != nil {
			return "", err
		}
		b[i] = tempPasswordChars[int(rb[0])%len(tempPasswordChars)]
	}
	return string(b), nil
}

func parseUserExpireDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("empty")
	}
	if t, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local); err == nil {
		return t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local), nil
	}
	return time.Time{}, errors.New("parse expireDate")
}

func defaultUserExpireDate() time.Time {
	return time.Now().AddDate(0, 1, 0)
}
