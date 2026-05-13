package service

import (
	"errors"
	"time"

	"adcms-backend/internal/config"
	"adcms-backend/internal/model"
	"adcms-backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResult 登录结果
type LoginResult struct {
	AccessToken  string      `json:"accessToken"`
	RefreshToken string      `json:"refreshToken"`
	Expires      int64       `json:"expires"`
	UserInfo     interface{} `json:"userInfo"`
}

// UserInfo 用户信息
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RealName string `json:"realName"`
	NickName string `json:"nickName"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Status   int8   `json:"status"`
	Avatar   string `json:"avatar"`
}

// TokenClaims JWT claims
type TokenClaims struct {
	UserID               string `json:"user_id"`
	Username             string `json:"username"`
	Role                 string `json:"role"`
	IsPlatformSuperAdmin bool   `json:"is_platform_super_admin"`
	jwt.RegisteredClaims
}

// AuthService 认证服务
type AuthService struct {
	userRepo *repository.UserRepository
	config   *config.Config
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repository.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// Login 用户登录
func (s *AuthService) Login(req *LoginRequest) (*LoginResult, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 检查状态
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	// 生成token
	accessToken, refreshToken, expires, err := s.generateTokens(user)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	// 更新最后登录时间
	_ = s.userRepo.UpdateLastLogin(user.ID)

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expires:      expires,
		UserInfo:     s.buildUserInfo(user),
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	// 解析refresh token
	claims, err := s.parseToken(refreshToken)
	if err != nil {
		return "", errors.New("刷新令牌无效")
	}

	// 查找用户
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return "", errors.New("用户不存在")
	}

	// 检查状态
	if user.Status != 1 {
		return "", errors.New("账号已被禁用")
	}

	// 生成新的access token
	accessToken, _, _, err := s.generateTokens(user)
	if err != nil {
		return "", errors.New("生成令牌失败")
	}

	return accessToken, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(userID string) (*UserInfo, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return s.buildUserInfo(user), nil
}

// GetAccessCodes 获取用户权限码
func (s *AuthService) GetAccessCodes(userID string, role string) []string {
	codes, err := s.userRepo.GetAccessCodesByUser(userID, role)
	if err != nil {
		return []string{}
	}
	return codes
}

// Impersonate 模拟登录为指定用户
func (s *AuthService) Impersonate(actorID string, targetUserID string) (*LoginResult, error) {
	// 检查操作者是否有权限模拟登录
	actor, err := s.userRepo.FindByID(actorID)
	if err != nil {
		return nil, errors.New("操作者不存在")
	}

	// 检查是否有模拟登录权限
	if !s.canImpersonate(actor) {
		return nil, errors.New("无权限执行模拟登录")
	}

	// 查找目标用户
	targetUser, err := s.userRepo.FindByID(targetUserID)
	if err != nil {
		return nil, errors.New("目标用户不存在")
	}

	// 检查目标用户状态
	if targetUser.Status != 1 {
		return nil, errors.New("目标用户账号已被禁用")
	}

	// 生成token
	accessToken, refreshToken, expires, err := s.generateTokens(targetUser)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expires:      expires,
		UserInfo:     s.buildUserInfo(targetUser),
	}, nil
}

// canImpersonate 判断用户是否有模拟登录权限
func (s *AuthService) canImpersonate(user *model.User) bool {
	// 平台超级管理员可以模拟登录
	if user.Role == "super_admin" {
		return true
	}
	return false
}

// generateTokens 生成JWT令牌
func (s *AuthService) generateTokens(user *model.User) (accessToken, refreshToken string, expires int64, err error) {
	now := time.Now()
	accessExpire := now.Add(time.Duration(s.config.JWT.AccessTokenExpire) * time.Second)
	refreshExpire := now.Add(time.Duration(s.config.JWT.RefreshTokenExpire) * time.Second)

	isSuperAdmin := user.Role == "super_admin"

	// Access Token
	accessClaims := TokenClaims{
		UserID:               user.ID,
		Username:             user.Username,
		Role:                 user.Role,
		IsPlatformSuperAdmin: isSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpire),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", 0, err
	}

	// Refresh Token
	refreshClaims := TokenClaims{
		UserID:               user.ID,
		Username:             user.Username,
		Role:                 user.Role,
		IsPlatformSuperAdmin: isSuperAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpire),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", "", 0, err
	}

	return accessToken, refreshToken, accessExpire.Unix(), nil
}

// parseToken 解析JWT令牌
func (s *AuthService) parseToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无效的令牌")
}

// buildUserInfo 构建用户信息
func (s *AuthService) buildUserInfo(user *model.User) *UserInfo {
	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		RealName: user.RealName,
		NickName: user.NickName,
		Email:    user.Email,
		Phone:    user.Phone,
		Role:     user.Role,
		Status:   user.Status,
	}
}
