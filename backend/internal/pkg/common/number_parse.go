package common

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseOptionalInt 解析可选数字（兼容 JSON number/string）。
func ParseOptionalInt(v interface{}) (*int, error) {
	if v == nil {
		return nil, nil
	}
	switch x := v.(type) {
	case float64:
		n := int(x)
		return &n, nil
	case float32:
		n := int(x)
		return &n, nil
	case int:
		n := x
		return &n, nil
	case int64:
		n := int(x)
		return &n, nil
	case int32:
		n := int(x)
		return &n, nil
	case string:
		s := strings.TrimSpace(x)
		if s == "" {
			return nil, nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return &n, nil
	default:
		s := strings.TrimSpace(fmt.Sprint(x))
		if s == "" {
			return nil, nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		return &n, nil
	}
}
