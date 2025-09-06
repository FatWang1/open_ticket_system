package auth

import (
	"errors"
	"time"

	"github.com/FatWang1/fatwang-go-utils/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/FatWang1/open_ticket_system/internal/middleware"
	"github.com/FatWang1/open_ticket_system/internal/models"
)

// AuthService 认证服务
type AuthService struct {
	db        *gorm.DB
	logger    utils.Logger
	jwtConfig *middleware.JWTConfig
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, jwtConfig *middleware.JWTConfig, logger utils.Logger) *AuthService {
	return &AuthService{
		db:        db,
		logger:    logger,
		jwtConfig: jwtConfig,
	}
}

// Login 用户登录
func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	var user models.User

	// 查找用户
	if err := s.db.Where("username = ? AND status = 1", req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Errorf("用户名或密码错误, err = %+v", err)
		return nil, errors.New("用户名或密码错误")
	}

	// 生成token
	accessToken, err := middleware.GenerateAccessToken(s.jwtConfig, user.ID, user.Username, "user")
	if err != nil {
		return nil, err
	}

	refreshToken, err := middleware.GenerateRefreshToken(s.jwtConfig, user.ID, user.Username, "user")
	if err != nil {
		return nil, err
	}

	// 更新最后登录时间
	now := time.Now()
	s.db.Model(&user).Update("last_login_at", now)

	// 清除密码字段
	user.Password = ""

	return &models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtConfig.AccessTokenExp.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// RefreshToken 刷新token
func (s *AuthService) RefreshToken(req *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	// 验证refresh token
	claims, err := middleware.ValidateRefreshToken(s.jwtConfig, req.RefreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	// 查找用户
	var user models.User
	if err := s.db.Where("id = ? AND status = 1", claims.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在或已被禁用")
		}
		return nil, err
	}

	// 生成新的token
	accessToken, err := middleware.GenerateAccessToken(s.jwtConfig, user.ID, user.Username, "user")
	if err != nil {
		return nil, err
	}

	refreshToken, err := middleware.GenerateRefreshToken(s.jwtConfig, user.ID, user.Username, "user")
	if err != nil {
		return nil, err
	}

	return &models.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtConfig.AccessTokenExp.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// Register 用户注册
func (s *AuthService) Register(req *models.RegisterRequest) (*models.RegisterResponse, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if req.Email != "" {
		if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Nickname: req.Nickname,
		Status:   1, // 正常状态
		RoleID:   1, // 默认角色ID
	}

	if err := s.db.Create(&user).Error; err != nil {
		s.logger.Errorf("Failed to create user - error: %v", err)
		return nil, err
	}

	// 清除密码字段
	user.Password = ""

	return &models.RegisterResponse{
		User: user,
	}, nil
}

// GetUserByID 根据ID获取用户信息
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ? AND status = 1", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// 清除密码字段
	user.Password = ""
	return &user, nil
}
