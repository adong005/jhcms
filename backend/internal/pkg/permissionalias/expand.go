package permissionalias

// EffectiveCodes 在权限校验与菜单过滤中使用：将历史/兼容权限码展开为当前生效的码集合（去重）。
func EffectiveCodes(codes []string) map[string]struct{} {
	out := make(map[string]struct{}, len(codes)+16)
	for _, c := range codes {
		if c == "" {
			continue
		}
		out[c] = struct{}{}
	}
	// 旧：单一 system:menu:update 覆盖写操作
	if _, ok := out["system:menu:update"]; ok {
		for _, x := range []string{
			"system:menu:create",
			"system:menu:update",
			"system:menu:delete",
			"system:menu:status",
			"system:menu:show",
		} {
			out[x] = struct{}{}
		}
	}
	// 旧：log:* 与 system:log:* 互认
	if _, ok := out["log:list"]; ok {
		out["system:log:list"] = struct{}{}
	}
	if _, ok := out["system:log:list"]; ok {
		out["log:list"] = struct{}{}
	}
	if _, ok := out["log:delete"]; ok {
		out["system:log:delete"] = struct{}{}
	}
	if _, ok := out["system:log:delete"]; ok {
		out["log:delete"] = struct{}{}
	}
	if _, ok := out["log:clear"]; ok {
		out["system:log:clear"] = struct{}{}
	}
	if _, ok := out["system:log:clear"]; ok {
		out["log:clear"] = struct{}{}
	}
	return out
}

// HasAny 若用户 codes 展开后包含 required 中任意一个则 true。
func HasAny(userCodes []string, required ...string) bool {
	if len(required) == 0 {
		return true
	}
	eff := EffectiveCodes(userCodes)
	for _, r := range required {
		if r == "" {
			return true
		}
		if _, ok := eff[r]; ok {
			return true
		}
	}
	return false
}
