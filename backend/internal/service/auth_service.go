package service

import (
	"adcms-backend/internal/config"
	"adcms-backend/internal/model"
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/jwt"
	"adcms-backend/internal/pkg/utils"
	"adcms-backend/internal/repository"
	"errors"
	"time"
)

type AuthService struct {
	userRepo *repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	User         *UserInfo `json:"user"`
}

type UserInfo struct {
	UserID               string   `json:"userId"`
	Username             string   `json:"username"`
	RealName             string   `json:"realName"`
	NickName             string   `json:"nickName"`
	Email                string   `json:"email"`
	Role                 string   `json:"role"`
	Roles                []string `json:"roles"`
	TenantID             string   `json:"tenantId"`
	IsAdmin              bool     `json:"isAdmin"`
	IsPlatformSuperAdmin bool     `json:"isPlatformSuperAdmin"`
	DataScope            string   `json:"dataScope"`
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 超过 1 个月未登录，自动锁定账号。
	if user.LastLoginDate != nil && user.LastLoginDate.Before(time.Now().AddDate(0, -1, 0)) {
		_ = s.userRepo.UpdateStatus(user.ID, 0)
		return nil, errors.New("账号超过一个月未登录，已自动锁定")
	}

	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	accessToken, err := jwt.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		user.EffectiveTenantID(),
		user.IsAdmin,
		user.ParentID,
		s.cfg.JWT.AccessTokenExpire,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		user.EffectiveTenantID(),
		user.IsAdmin,
		user.ParentID,
		s.cfg.JWT.RefreshTokenExpire,
	)
	if err != nil {
		return nil, err
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	userInfo := userInfoFromModel(user)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userInfo,
	}, nil
}

func userInfoFromModel(user *model.User) *UserInfo {
	return &UserInfo{
		UserID:               user.ID,
		Username:             user.Username,
		RealName:             user.RealName,
		NickName:             user.NickName,
		Email:                user.Email,
		Role:                 user.Role,
		Roles:                []string{user.Role},
		TenantID:             user.EffectiveTenantID(),
		IsAdmin:              user.IsAdmin,
		IsPlatformSuperAdmin: user.IsAdmin,
		DataScope:            user.DataScope,
	}
}

func (s *AuthService) canImpersonate(actor, target *model.User) bool {
	if actor.IsAdmin {
		return true
	}
	if actor.Role == "super_admin" {
		return true
	}
	if actor.Role == "admin" &&
		actor.EffectiveTenantID() == target.EffectiveTenantID() &&
		actor.EffectiveTenantID() != "" {
		return true
	}
	return false
}

// Impersonate 以目标用户身份签发令牌
func (s *AuthService) Impersonate(actorUserID string, targetUserIDStr string) (*LoginResponse, error) {
	actor, err := s.userRepo.FindByID(actorUserID)
	if err != nil {
		return nil, errors.New("无法验证当前用户")
	}
	if actor.Status != 1 {
		return nil, errors.New("当前账号已禁用")
	}

	if !ids.Valid(targetUserIDStr) {
		return nil, errors.New("无效的用户ID")
	}
	user, err := s.userRepo.FindByID(targetUserIDStr)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.Status != 1 {
		return nil, errors.New("目标账号已禁用，无法进入")
	}
	if !s.canImpersonate(actor, user) {
		return nil, errors.New("无权进入该账户：需超级管理员角色、平台超级管理员标识，或同租户管理员")
	}

	accessToken, err := jwt.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		user.EffectiveTenantID(),
		user.IsAdmin,
		user.ParentID,
		s.cfg.JWT.AccessTokenExpire,
	)
	if err != nil {
		return nil, err
	}
	refreshToken, err := jwt.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		user.EffectiveTenantID(),
		user.IsAdmin,
		user.ParentID,
		s.cfg.JWT.RefreshTokenExpire,
	)
	if err != nil {
		return nil, err
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	userInfo := userInfoFromModel(user)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userInfo,
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := jwt.ParseToken(refreshToken)
	if err != nil {
		return "", errors.New("令牌无效或已过期")
	}

	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", errors.New("用户不存在")
	}

	if user.Status != 1 {
		return "", errors.New("账号已被禁用")
	}

	accessToken, err := jwt.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
		user.EffectiveTenantID(),
		user.IsAdmin,
		user.ParentID,
		s.cfg.JWT.AccessTokenExpire,
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID string) (*UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return userInfoFromModel(user), nil
}

// GetAccessCodes 获取权限码
func (s *AuthService) GetAccessCodes(userID string, role string) []string {
	codes, err := s.userRepo.GetAccessCodesByUser(userID, role)
	if err != nil {
		return []string{}
	}
	return codes
}
