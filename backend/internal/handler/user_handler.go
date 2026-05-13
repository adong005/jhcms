package handler

import (
	"bytes"
	"encoding/json"
	"strings"

	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/common"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/pkg/utils"
	"adcms-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

// GetUserList 获取用户列表
func (h *UserHandler) GetUserList(c *gin.Context) {
	roleVal, _ := c.Get("role")
	roleStr, _ := roleVal.(string)
	if roleStr != "super_admin" && roleStr != "admin" && roleStr != "user" {
		response.Forbidden(c, "无权访问用户列表")
		return
	}

	var req struct {
		Page     int                       `json:"page"`
		PageSize int                       `json:"pageSize"`
		Username string                    `json:"username"`
		RealName string                    `json:"realName"`
		Status   common.OptionalListStatus `json:"status"`
	}

	common.HandleListRequest(c, &req, 10, func() (interface{}, int64, error) {
		tenantID, _ := c.Get("tenant_id")
		currentUserIDVal, _ := c.Get("user_id")

		tenantIDStr, _ := tenantID.(string)
		currentUserID, _ := currentUserIDVal.(string)

		users, total, err := h.userRepo.List(tenantIDStr, roleStr, currentUserID, req.Page, req.PageSize, req.Username, req.Status.Ptr())
		if err != nil {
			return nil, 0, err
		}

		items := make([]map[string]interface{}, 0, len(users))
		for _, user := range users {
			item := map[string]interface{}{
				"id":         user.ID,
				"username":   user.Username,
				"realName":   user.RealName,
				"nickName":   user.NickName,
				"email":      user.Email,
				"phone":      user.Phone,
				"role":       user.Role,
				"status":     user.Status,
				"createTime": user.CreatedAt.Format("2006-01-02 15:04:05"),
				"updateTime": user.UpdatedAt.Format("2006-01-02 15:04:05"),
			}

			if user.CreatedBy != nil {
				item["createdBy"] = *user.CreatedBy
			}

			if user.LastLoginDate != nil {
				item["lastLoginDate"] = user.LastLoginDate.Format("2006-01-02 15:04:05")
			}

			if user.ExpireDate != nil {
				item["expireDate"] = user.ExpireDate.Format("2006-01-02 15:04:05")
			} else {
				item["expireDate"] = defaultUserExpireDate().Format("2006-01-02 15:04:05")
			}

			items = append(items, item)
		}

		return items, total, nil
	}, "获取用户列表失败")
}

// GetUser 获取用户详情
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		response.Error(c, "无效的用户ID")
		return
	}

	user, err := h.userRepo.FindByID(idStr)
	if err != nil {
		response.Error(c, "获取用户详情失败")
		return
	}

	item := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"realName":   user.RealName,
		"nickName":   user.NickName,
		"email":      user.Email,
		"phone":      user.Phone,
		"role":       user.Role,
		"status":     user.Status,
		"createTime": user.CreatedAt.Format("2006-01-02 15:04:05"),
		"updateTime": user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if user.CreatedBy != nil {
		item["createdBy"] = *user.CreatedBy
	}
	if user.ExpireDate != nil {
		item["expireDate"] = user.ExpireDate.Format("2006-01-02 15:04:05")
	} else {
		item["expireDate"] = defaultUserExpireDate().Format("2006-01-02 15:04:05")
	}

	response.Success(c, item)
}

// CreateUser 创建用户（用手动 JSON 解析，兼容 role 为数字/字符串，避免 ShouldBindJSON 与旧部署不一致）
func (h *UserHandler) CreateUser(c *gin.Context) {
	dec := json.NewDecoder(c.Request.Body)
	dec.UseNumber()
	var raw map[string]interface{}
	if err := dec.Decode(&raw); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	username := strings.TrimSpace(jsonStringField(raw["username"]))
	password := jsonStringField(raw["password"])
	if username == "" || password == "" {
		response.Error(c, "请求参数错误")
		return
	}

	roleStr, err := coerceRoleFromInterface(raw["role"])
	if err != nil {
		response.Error(c, "role 无效：请传 1/2/3 或 super_admin/admin/user")
		return
	}
	operatorRoleVal, _ := c.Get("role")
	operatorRole, _ := operatorRoleVal.(string)

	currentUserIDVal, _ := c.Get("user_id")
	currentUserID, _ := currentUserIDVal.(string)
	if !h.canAssignRole(operatorRole, currentUserID, roleStr) {
		response.Error(c, "无权分配该角色")
		return
	}

	// 创建用户对象
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		response.Error(c, "密码处理失败")
		return
	}
	status := int8(1)
	if v, ok := raw["status"]; ok {
		if s, err := coerceInt8FromJSON(v); err == nil && s != 0 {
			status = s
		}
	}

	user := &model.User{
		Username:        username,
		Password:        hashedPwd,
		RealName:        jsonStringField(raw["realName"]),
		NickName:        jsonStringField(raw["nickName"]),
		Email:           jsonStringField(raw["email"]),
		Phone:           jsonStringField(raw["phone"]),
		Role:            roleStr,
		Status:          status,
		CreatorOptional: model.CreatorOptional{CreatedBy: &currentUserID},
		DataScope:       dataScopeForRole(roleStr),
	}
	expStr := jsonStringField(raw["expireDate"])
	if expStr != "" {
		exp, perr := parseUserExpireDate(expStr)
		if perr != nil {
			response.Error(c, "expireDate 格式错误，支持 yyyy-MM-dd 或 yyyy-MM-dd HH:mm:ss")
			return
		}
		user.ExpireDate = &exp
	}

	// 如果是管理员创建用户，设置 parent_id
	if roleStr == "user" {
		user.ParentID = &currentUserID
	}

	err = h.userRepo.Create(user)
	if err != nil {
		em := strings.ToLower(err.Error())
		if strings.Contains(em, "duplicate") || strings.Contains(em, "1062") {
			response.Error(c, "用户名已存在")
			return
		}
		response.Error(c, "创建用户失败")
		return
	}

	response.SuccessWithMessage(c, "创建用户成功", nil)
}

// resetPasswordForUser 将 user 重置为随机密码（租户校验），成功时写入 JSON 响应
func (h *UserHandler) resetPasswordForUser(c *gin.Context, user *model.User) {
	if tidVal, ok := c.Get("tenant_id"); ok {
		if tid, ok2 := tidVal.(string); ok2 && tid != "" && user.EffectiveTenantID() != tid {
			response.Error(c, "无权操作该用户")
			return
		}
	}
	plain, err := randomTempPassword(10)
	if err != nil {
		response.Error(c, "生成密码失败")
		return
	}
	hashed, err := utils.HashPassword(plain)
	if err != nil {
		response.Error(c, "密码处理失败")
		return
	}
	if err := h.userRepo.UpdatePasswordHash(user.ID, hashed); err != nil {
		response.Error(c, "重置密码失败")
		return
	}
	response.SuccessWithMessage(c, "重置密码成功", gin.H{
		"newPassword": plain,
		"email":       user.Email,
	})
}

// ResetPassword POST /user/reset-password（与 UpdateUserByBody+resetPassword 二选一即可）
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的用户ID")
		return
	}
	user, err := h.userRepo.FindByID(req.ID)
	if err != nil {
		response.Error(c, "用户不存在")
		return
	}
	h.resetPasswordForUser(c, user)
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		response.Error(c, "无效的用户ID")
		return
	}

	var req struct {
		RealName string `json:"realName"`
		NickName string `json:"nickName"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Status   *int8  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	user, err := h.userRepo.FindByID(idStr)
	if err != nil {
		response.Error(c, "用户不存在")
		return
	}

	// 更新字段
	user.RealName = req.RealName
	user.NickName = req.NickName
	user.Email = req.Email
	user.Phone = req.Phone
	if req.Status != nil {
		user.Status = *req.Status
	}

	err = h.userRepo.Update(user)
	if err != nil {
		response.Error(c, "更新用户失败")
		return
	}

	response.SuccessWithMessage(c, "更新用户成功", nil)
}

// UpdateUserByBody 兼容前端 POST /user/update（body 传 id，id 仅接受字符串以与前端统一）
func (h *UserHandler) UpdateUserByBody(c *gin.Context) {
	var req struct {
		ID            string          `json:"id" binding:"required"`
		Username      string          `json:"username"`
		RealName      string          `json:"realName"`
		NickName      string          `json:"nickName"`
		Email         string          `json:"email"`
		Phone         string          `json:"phone"`
		Role          json.RawMessage `json:"role"`
		Status        *int8           `json:"status"`
		Password      string          `json:"password"`
		ExpireDate    string          `json:"expireDate"`
		ResetPassword *bool           `json:"resetPassword"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	if !ids.Valid(req.ID) {
		response.Error(c, "无效的用户ID")
		return
	}

	user, err := h.userRepo.FindByID(req.ID)
	if err != nil {
		response.Error(c, "用户不存在")
		return
	}

	// 走已有 /user/update，避免新路径在旧部署或反向代理上 404
	if req.ResetPassword != nil && *req.ResetPassword {
		h.resetPasswordForUser(c, user)
		return
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if len(bytes.TrimSpace(req.Role)) > 0 {
		r, rerr := parseUserRoleJSON(req.Role)
		if rerr != nil {
			response.Error(c, "role 无效：请传 1/2/3 或 super_admin/admin/user")
			return
		}
		operatorRoleVal, _ := c.Get("role")
		operatorRole, _ := operatorRoleVal.(string)
		currentUserIDVal, _ := c.Get("user_id")
		currentUserID, _ := currentUserIDVal.(string)
		if !h.canAssignRole(operatorRole, currentUserID, r) {
			response.Error(c, "无权分配该角色")
			return
		}
		user.Role = r
		user.DataScope = dataScopeForRole(r)
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if req.Password != "" {
		hashed, hashErr := utils.HashPassword(req.Password)
		if hashErr != nil {
			response.Error(c, "密码处理失败")
			return
		}
		user.Password = hashed
	}

	if req.ExpireDate == "" {
		def := defaultUserExpireDate()
		user.ExpireDate = &def
	} else {
		parsed, parseErr := parseUserExpireDate(req.ExpireDate)
		if parseErr != nil {
			response.Error(c, "expireDate 格式错误，支持 yyyy-MM-dd 或 yyyy-MM-dd HH:mm:ss")
			return
		}
		user.ExpireDate = &parsed
	}

	if err = h.userRepo.Update(user); err != nil {
		response.Error(c, "更新用户失败")
		return
	}

	response.SuccessWithMessage(c, "更新用户成功", nil)
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		response.Error(c, "无效的用户ID")
		return
	}

	_, err := h.userRepo.FindByID(idStr)
	if err != nil {
		response.Error(c, "用户不存在")
		return
	}

	if err = h.userRepo.Delete(idStr); err != nil {
		response.Error(c, "删除用户失败")
		return
	}
	response.SuccessWithMessage(c, "删除用户成功", nil)
}

// UpdateUserStatus 兼容前端 POST /user/status
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	var req struct {
		ID     string `json:"id" binding:"required"`
		Status *int8  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的用户ID")
		return
	}
	if err := h.userRepo.UpdateStatus(req.ID, *req.Status); err != nil {
		response.Error(c, "更新用户状态失败")
		return
	}
	response.SuccessWithMessage(c, "更新用户状态成功", nil)
}

// DeleteUserByBody 兼容前端 POST /user/delete
func (h *UserHandler) DeleteUserByBody(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !ids.Valid(req.ID) {
		response.Error(c, "无效的用户ID")
		return
	}
	if err := h.userRepo.Delete(req.ID); err != nil {
		response.Error(c, "删除用户失败")
		return
	}
	response.SuccessWithMessage(c, "删除用户成功", nil)
}

// BatchDeleteUsers 兼容前端 POST /user/batch-delete
func (h *UserHandler) BatchDeleteUsers(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	validIDs := make([]string, 0, len(req.IDs))
	for _, s := range req.IDs {
		if ids.Valid(s) {
			validIDs = append(validIDs, s)
		}
	}
	if len(validIDs) == 0 {
		response.Error(c, "缺少有效的用户ID")
		return
	}
	if err := h.userRepo.BatchDelete(validIDs); err != nil {
		response.Error(c, "批量删除用户失败")
		return
	}
	response.SuccessWithMessage(c, "批量删除用户成功", nil)
}

func (h *UserHandler) canAssignRole(operatorRole, operatorUserID, targetRole string) bool {
	targetRole = strings.TrimSpace(targetRole)
	if targetRole == "" {
		return false
	}
	switch operatorRole {
	case "super_admin":
		ok, err := h.userRepo.RoleExistsByCode(targetRole)
		return err == nil && ok
	case "admin":
		if targetRole == "user" {
			return true
		}
		if targetRole == "admin" || targetRole == "super_admin" {
			return false
		}
		ok, err := h.userRepo.RoleExistsByCodeAndCreator(targetRole, operatorUserID)
		return err == nil && ok
	default:
		return false
	}
}
