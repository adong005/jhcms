package handler

import (
	"fmt"
	"strings"

	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/response"
	"adcms-backend/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

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
