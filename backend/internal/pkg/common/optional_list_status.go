package common

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
)

// OptionalListStatus 用于列表接口的 status 筛选 JSON 绑定。
// 前端筛选「全部」时常传 "" 或 null；直接绑定 *int 会导致 ShouldBindJSON 失败。
// 支持：缺省、null、""、数字、字符串形式的数字（如 "1"）。
type OptionalListStatus struct {
	p *int
}

// Ptr 返回 nil 表示不按 status 过滤。
func (o OptionalListStatus) Ptr() *int {
	return o.p
}

func (o *OptionalListStatus) UnmarshalJSON(b []byte) error {
	o.p = nil
	if len(b) == 0 {
		return nil
	}
	if bytes.Equal(b, []byte("null")) {
		return nil
	}
	if len(b) >= 2 && b[0] == '"' {
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return err
		}
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		o.p = &n
		return nil
	}
	var n int
	if err := json.Unmarshal(b, &n); err != nil {
		return err
	}
	o.p = &n
	return nil
}
