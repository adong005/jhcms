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

func (h *UserHandler) currentUser(c *gin.Context) (*model.User, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "未登录")
		return nil, false
	}
	userID, ok := userIDVal.(string)
	if !ok || userID == "" {
		response.Unauthorized(c, "未登录")
		return nil, false
	}
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		response.Error(c, "用户不存在")
		return nil, false
	}
	return user, true
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	var req struct {
		RealName string `json:"realName"`
		NickName string `json:"nickName"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if req.RealName != "" {
		user.RealName = req.RealName
	}
	if req.NickName != "" {
		user.NickName = req.NickName
	}
	// 个人中心不允许修改用户名，忽略 req.Username
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, "更新个人资料失败")
		return
	}
	response.SuccessWithMessage(c, "更新个人资料成功", nil)
}

func (h *UserHandler) GetSecuritySettings(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{
		"accountPassword":      user.Password != "",
		"securityPhone":        user.Phone != "",
		"securityPhoneNumber":  user.Phone,
		"securityQuestion":     user.SecurityQuestion1 != "" || user.SecurityQuestion2 != "",
		"securityEmail":        user.Email != "",
		"securityEmailAddress": user.Email,
		"securityMfa":          user.GoogleAuthBound,
		"passwordStrength":     "中",
	})
}

func (h *UserHandler) UpdateSecuritySettings(c *gin.Context) {
	// 兼容前端占位调用：当前使用各子项独立更新接口
	response.SuccessWithMessage(c, "更新成功", nil)
}

func (h *UserHandler) UpdatePassword(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	var req struct {
		OldPassword string `json:"oldPassword" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		response.Error(c, "旧密码错误")
		return
	}
	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		response.Error(c, "密码处理失败")
		return
	}
	user.Password = hashed
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, "修改密码失败")
		return
	}
	response.SuccessWithMessage(c, "修改密码成功", nil)
}

func (h *UserHandler) GetPhoneSetting(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"phone": user.Phone})
}

func (h *UserHandler) UpdatePhoneSetting(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	var req struct {
		NewPhone string `json:"newPhone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	user.Phone = strings.TrimSpace(req.NewPhone)
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, "更新手机号失败")
		return
	}
	response.SuccessWithMessage(c, "更新手机号成功", nil)
}

func (h *UserHandler) GetQuestionSetting(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{
		"question1": user.SecurityQuestion1,
		"answer1":   user.SecurityAnswer1,
		"question2": user.SecurityQuestion2,
		"answer2":   user.SecurityAnswer2,
	})
}

func (h *UserHandler) UpdateQuestionSetting(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	var req struct {
		Question1 string `json:"question1" binding:"required"`
		Answer1   string `json:"answer1" binding:"required"`
		Question2 string `json:"question2" binding:"required"`
		Answer2   string `json:"answer2" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	user.SecurityQuestion1 = req.Question1
	user.SecurityAnswer1 = req.Answer1
	user.SecurityQuestion2 = req.Question2
	user.SecurityAnswer2 = req.Answer2
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, "更新密保问题失败")
		return
	}
	response.SuccessWithMessage(c, "更新密保问题成功", nil)
}

func (h *UserHandler) GetEmailSetting(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{"email": user.Email})
}

func (h *UserHandler) UpdateEmailSetting(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	var req struct {
		NewEmail string `json:"newEmail" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	user.Email = strings.TrimSpace(req.NewEmail)
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, "更新邮箱失败")
		return
	}
	response.SuccessWithMessage(c, "更新邮箱成功", nil)
}

func (h *UserHandler) GetGoogleAuthSetting(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	secret := user.GoogleAuthSecret
	if secret == "" {
		s, err := randomTempPassword(16)
		if err == nil {
			secret = s
			user.GoogleAuthSecret = s
			_ = h.userRepo.Update(user)
		}
	}
	qrURL := ""
	if secret != "" {
		qrURL = fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=otpauth://totp/ADCMS:%s?secret=%s&issuer=ADCMS", user.Username, secret)
	}
	response.Success(c, gin.H{
		"isBound":   user.GoogleAuthBound,
		"qrCodeUrl": qrURL,
		"secretKey": secret,
	})
}

func (h *UserHandler) BindGoogleAuth(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	var req struct {
		VerifyCode string `json:"verifyCode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	user.GoogleAuthBound = true
	if user.GoogleAuthSecret == "" {
		user.GoogleAuthSecret, _ = randomTempPassword(16)
	}
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, "绑定谷歌验证器失败")
		return
	}
	response.SuccessWithMessage(c, "绑定成功", nil)
}

func (h *UserHandler) UnbindGoogleAuth(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	user.GoogleAuthBound = false
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, "解绑谷歌验证器失败")
		return
	}
	response.SuccessWithMessage(c, "解绑成功", nil)
}

func (h *UserHandler) GetNotificationSettings(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	response.Success(c, gin.H{
		"accountPassword": user.NotifyAccountPassword,
		"systemMessage":   user.NotifySystemMessage,
		"todoTask":        user.NotifyTodoTask,
	})
}

func (h *UserHandler) UpdateNotificationSettings(c *gin.Context) {
	user, ok := h.currentUser(c)
	if !ok {
		return
	}
	var req struct {
		AccountPassword *bool `json:"accountPassword"`
		SystemMessage   *bool `json:"systemMessage"`
		TodoTask        *bool `json:"todoTask"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}
	if req.AccountPassword != nil {
		user.NotifyAccountPassword = *req.AccountPassword
	}
	if req.SystemMessage != nil {
		user.NotifySystemMessage = *req.SystemMessage
	}
	if req.TodoTask != nil {
		user.NotifyTodoTask = *req.TodoTask
	}
	if err := h.userRepo.Update(user); err != nil {
		response.Error(c, "更新通知设置失败")
		return
	}
	response.SuccessWithMessage(c, "更新通知设置成功", nil)
}

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
