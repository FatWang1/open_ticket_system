package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/FatWang1/open_ticket_system/internal/middleware"
	"github.com/FatWang1/open_ticket_system/internal/models"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	// 自动迁移
	db.AutoMigrate(&models.User{})

	return db
}

func TestAuthService_Login(t *testing.T) {
	db := setupTestDB()
	jwtConfig := &middleware.JWTConfig{
		Secret:          "test-secret",
		AccessTokenExp:  time.Hour,
		RefreshTokenExp: time.Hour * 24 * 7,
	}

	service := NewAuthService(db, jwtConfig)

	// 创建测试用户
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Username: "testuser",
		Password: string(hashedPassword),
		Email:    "test@example.com",
		Status:   1,
		RoleID:   1,
	}
	db.Create(&user)

	// 测试登录
	req := &models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	response, err := service.Login(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, "Bearer", response.TokenType)
	assert.Equal(t, "testuser", response.User.Username)
}

func TestAuthService_Register(t *testing.T) {
	db := setupTestDB()
	jwtConfig := &middleware.JWTConfig{
		Secret:          "test-secret",
		AccessTokenExp:  time.Hour,
		RefreshTokenExp: time.Hour * 24 * 7,
	}

	service := NewAuthService(db, jwtConfig)

	// 测试注册
	req := &models.RegisterRequest{
		Username: "newuser",
		Password: "password123",
		Email:    "newuser@example.com",
		Nickname: "新用户",
	}

	response, err := service.Register(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "newuser", response.User.Username)
	assert.Equal(t, "newuser@example.com", response.User.Email)
	assert.Equal(t, "新用户", response.User.Nickname)
}

func TestAuthService_RefreshToken(t *testing.T) {
	db := setupTestDB()
	jwtConfig := &middleware.JWTConfig{
		Secret:          "test-secret",
		AccessTokenExp:  time.Hour,
		RefreshTokenExp: time.Hour * 24 * 7,
	}

	service := NewAuthService(db, jwtConfig)

	// 创建测试用户
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := models.User{
		Username: "testuser",
		Password: string(hashedPassword),
		Email:    "test@example.com",
		Status:   1,
		RoleID:   1,
	}
	db.Create(&user)

	// 生成refresh token
	refreshToken, err := middleware.GenerateRefreshToken(jwtConfig, user.ID, user.Username, "user")
	assert.NoError(t, err)

	// 测试刷新token
	req := &models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	response, err := service.RefreshToken(req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.AccessToken)
	assert.NotEmpty(t, response.RefreshToken)
	assert.Equal(t, "Bearer", response.TokenType)
}
