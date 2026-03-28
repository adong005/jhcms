package handler

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// 与 curl / 浏览器一致的 JSON，确保手动解析可用
func TestCreateUserJSON_decodeAndCoerce(t *testing.T) {
	raw := `{"status":1,"username":"dong005","nickName":"ad","email":"adong005@gmail.com","phone":"18503290567","role":1,"expireDate":"2026-03-29","password":"022018"}`
	dec := json.NewDecoder(bytes.NewReader([]byte(raw)))
	dec.UseNumber()
	var m map[string]interface{}
	if err := dec.Decode(&m); err != nil {
		t.Fatal(err)
	}
	u := strings.TrimSpace(jsonStringField(m["username"]))
	p := jsonStringField(m["password"])
	if u != "dong005" || p != "022018" {
		t.Fatalf("username/password got %q %q", u, p)
	}
	r, err := coerceRoleFromInterface(m["role"])
	if err != nil || r != "super_admin" {
		t.Fatalf("role got %q err %v", r, err)
	}
	exp := jsonStringField(m["expireDate"])
	if _, err := parseUserExpireDate(exp); err != nil {
		t.Fatal(err)
	}
}
