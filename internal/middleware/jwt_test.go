package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestJWTAuthMiddleware 测试JWT认证中间件
func TestJWTAuthMiddleware(t *testing.T) {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupAuth      func() string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "should_allow_access_with_valid_jwt",
			setupAuth: func() string {
				config := &JWTConfig{Secret: "test-secret"}
				token, _ := GenerateToken(config, 1, "user1", "admin")
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"user_id":1,"username":"user1","role":"admin"}`,
		},
		{
			name: "should_reject_access_without_token",
			setupAuth: func() string {
				return ""
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Authorization header is required"}`,
		},
		{
			name: "should_reject_access_with_invalid_token",
			setupAuth: func() string {
				return "Bearer invalid-token"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid or expired token"}`,
		},
		{
			name: "should_reject_access_with_expired_token",
			setupAuth: func() string {
				config := &JWTConfig{Secret: "test-secret"}
				// 创建过期的token
				claims := &Claims{
					UserID:   1,
					Username: "user1",
					Role:     "admin",
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // 1小时前过期
					},
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(config.Secret))
				return "Bearer " + tokenString
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid or expired token"}`,
		},
		{
			name: "should_reject_access_with_malformed_header",
			setupAuth: func() string {
				return "InvalidHeaderFormat"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid authorization header format"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建Gin引擎
			router := gin.New()

			// 创建JWT配置
			config := &JWTConfig{Secret: "test-secret"}

			// 添加JWT中间件
			router.Use(JWTAuthMiddleware(config))

			// 添加测试路由
			router.GET("/test", func(c *gin.Context) {
				userID := c.GetUint("user_id")
				username := c.GetString("username")
				role := c.GetString("role")
				c.JSON(http.StatusOK, gin.H{
					"user_id":  userID,
					"username": username,
					"role":     role,
				})
			})

			// 创建请求
			req, err := http.NewRequest("GET", "/test", nil)
			require.NoError(t, err)

			// 设置认证头
			if authHeader := tt.setupAuth(); authHeader != "" {
				req.Header.Set("Authorization", authHeader)
			}

			// 创建响应记录器
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 验证响应
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}

// TestGenerateToken 测试JWT token生成
func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		role        string
		userID      uint
		expectError bool
	}{
		{
			name:        "should_generate_valid_token",
			username:    "user1",
			role:        "admin",
			userID:      1,
			expectError: false,
		},
		{
			name:        "should_generate_token_with_empty_username",
			username:    "",
			role:        "user",
			userID:      2,
			expectError: false,
		},
		{
			name:        "should_generate_token_with_empty_role",
			username:    "user3",
			role:        "",
			userID:      3,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &JWTConfig{Secret: "test-secret"}
			token, err := GenerateToken(config, tt.userID, tt.username, tt.role)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// 验证token可以正确解析
				parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(config.Secret), nil
				})
				require.NoError(t, err)

				claims, ok := parsedToken.Claims.(*Claims)
				require.True(t, ok)
				assert.Equal(t, tt.username, claims.Username)
				assert.Equal(t, tt.role, claims.Role)
				assert.Equal(t, tt.userID, claims.UserID)
			}
		})
	}
}

// TestJWTConfig 测试JWT配置
func TestJWTConfig(t *testing.T) {
	t.Run("should_create_valid_jwt_config", func(t *testing.T) {
		config := &JWTConfig{
			Secret: "test-secret-key",
		}

		assert.Equal(t, "test-secret-key", config.Secret)
	})

	t.Run("should_validate_jwt_config", func(t *testing.T) {
		config := &JWTConfig{
			Secret: "test-secret-key",
		}

		// 测试token生成和验证
		token, err := GenerateToken(config, 1, "testuser", "admin")
		require.NoError(t, err)

		// 验证token
		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Secret), nil
		})
		require.NoError(t, err)

		claims, ok := parsedToken.Claims.(*Claims)
		require.True(t, ok)
		assert.Equal(t, "testuser", claims.Username)
		assert.Equal(t, "admin", claims.Role)
		assert.Equal(t, uint(1), claims.UserID)
	})
}

// TestJWTMiddlewareIntegration 测试JWT中间件集成
func TestJWTMiddlewareIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should_protect_protected_routes", func(t *testing.T) {
		router := gin.New()
		config := &JWTConfig{Secret: "test-secret"}

		// 公开路由
		router.GET("/public", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "public"})
		})

		// 受保护的路由组
		protected := router.Group("/protected")
		protected.Use(JWTAuthMiddleware(config))
		{
			protected.GET("/profile", func(c *gin.Context) {
				userID := c.GetUint("user_id")
				username := c.GetString("username")
				role := c.GetString("role")
				c.JSON(http.StatusOK, gin.H{
					"user_id":  userID,
					"username": username,
					"role":     role,
				})
			})
		}

		// 测试公开路由 - 应该可以访问
		req, _ := http.NewRequest("GET", "/public", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// 测试受保护路由 - 没有token应该被拒绝
		req, _ = http.NewRequest("GET", "/protected/profile", nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		// 测试受保护路由 - 有有效token应该可以访问
		token, _ := GenerateToken(config, 1, "testuser", "admin")
		req, _ = http.NewRequest("GET", "/protected/profile", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestJWTSecurity 测试JWT安全性
func TestJWTSecurity(t *testing.T) {
	t.Run("should_reject_tokens_with_wrong_secret", func(t *testing.T) {
		config1 := &JWTConfig{Secret: "secret1"}
		config2 := &JWTConfig{Secret: "secret2"}

		// 用config1生成token
		token, err := GenerateToken(config1, 1, "user1", "admin")
		require.NoError(t, err)

		// 用config2验证token应该失败
		parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config2.Secret), nil
		})
		// 即使解析成功，验证也应该失败
		if err == nil {
			assert.False(t, parsedToken.Valid)
		} else {
			assert.Error(t, err)
		}
	})

	t.Run("should_reject_tampered_tokens", func(t *testing.T) {
		config := &JWTConfig{Secret: "test-secret"}
		token, _ := GenerateToken(config, 1, "user1", "admin")

		// 篡改token
		tamperedToken := token + "tampered"

		// 验证篡改的token应该失败
		parsedToken, err := jwt.ParseWithClaims(tamperedToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Secret), nil
		})
		// 篡改的token应该无法解析或验证失败
		if err == nil {
			assert.False(t, parsedToken.Valid)
		} else {
			assert.Error(t, err)
		}
	})
}

// TestJWTPermissions 测试JWT权限控制
func TestJWTPermissions(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should_enforce_role_based_access_control", func(t *testing.T) {
		router := gin.New()
		config := &JWTConfig{Secret: "test-secret"}

		// 管理员路由
		admin := router.Group("/admin")
		admin.Use(JWTAuthMiddleware(config))
		admin.Use(func(c *gin.Context) {
			role := c.GetString("role")
			if role != "admin" {
				c.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
				c.Abort()
				return
			}
		})
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "admin dashboard"})
			})
		}

		// 普通用户路由
		user := router.Group("/user")
		user.Use(JWTAuthMiddleware(config))
		{
			user.GET("/profile", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "user profile"})
			})
		}

		// 测试管理员访问管理员路由
		adminToken, _ := GenerateToken(config, 1, "admin1", "admin")
		req, _ := http.NewRequest("GET", "/admin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		// 测试普通用户访问管理员路由
		userToken, _ := GenerateToken(config, 2, "user1", "user")
		req, _ = http.NewRequest("GET", "/admin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusForbidden, w.Code)

		// 测试普通用户访问普通用户路由
		req, _ = http.NewRequest("GET", "/user/profile", nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
