package auth

import (
	"github.com/FatWang1/fatwang-go-utils/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/FatWang1/open_ticket_system/auth/validator"
	"github.com/FatWang1/open_ticket_system/internal/middleware"
	"github.com/FatWang1/open_ticket_system/internal/models"
)

// AuthController 认证控制器
type AuthController struct {
	authService *AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(db *gorm.DB, jwtConfig *middleware.JWTConfig, logger utils.Logger) *AuthController {
	return &AuthController{
		authService: NewAuthService(db, jwtConfig, logger),
	}
}

// RegisterRoutes 注册路由
func (ac *AuthController) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", ac.Login)
		auth.POST("/register", ac.Register)
		auth.POST("/refresh", ac.RefreshToken)
		auth.GET("/profile", middleware.JWTAuthMiddleware(ac.authService.jwtConfig), ac.GetProfile)
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取JWT token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "登录请求"
// @Success 200 {object} models.LoginResponse "登录成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "用户名或密码错误"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "details": err.Error()})
		return
	}

	// 验证请求参数
	if err := validator.ValidateLoginRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ac.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册新账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "注册请求"
// @Success 201 {object} models.RegisterResponse "注册成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 409 {object} models.ErrorResponse "用户名或邮箱已存在"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "details": err.Error()})
		return
	}

	// 验证请求参数
	if err := validator.ValidateRegisterRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ac.authService.Register(&req)
	if err != nil {
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已存在" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// RefreshToken 刷新token
// @Summary 刷新JWT token
// @Description 使用refresh token获取新的access token
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.RefreshTokenRequest true "刷新token请求"
// @Success 200 {object} models.RefreshTokenResponse "刷新成功"
// @Failure 400 {object} models.ErrorResponse "请求参数错误"
// @Failure 401 {object} models.ErrorResponse "无效的刷新令牌"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/refresh [post]
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "details": err.Error()})
		return
	}

	// 验证请求参数
	if err := validator.ValidateRefreshTokenRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ac.authService.RefreshToken(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetProfile 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.User "用户信息"
// @Failure 401 {object} models.ErrorResponse "未授权"
// @Failure 500 {object} models.ErrorResponse "服务器内部错误"
// @Router /api/v1/auth/profile [get]
func (ac *AuthController) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := ac.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
