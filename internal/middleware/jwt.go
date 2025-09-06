package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTConfig JWT配置
type JWTConfig struct {
	Secret          string
	AccessTokenExp  time.Duration // access token过期时间
	RefreshTokenExp time.Duration // refresh token过期时间
}

// Claims JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Type     string `json:"type"` // "access" 或 "refresh"
	jwt.RegisteredClaims
}

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware(config *JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// 验证签名方法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return []byte(config.Secret), nil
		}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithoutClaimsValidation())

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is not valid"})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		if claims, ok := token.Claims.(*Claims); ok {
			// 手动验证过期时间
			if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
				c.Abort()
				return
			}

			// 只允许access token通过认证中间件
			if claims.Type != "access" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token type"})
				c.Abort()
				return
			}
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)
		}

		c.Next()
	}
}

// GenerateAccessToken 生成访问令牌
func GenerateAccessToken(config *JWTConfig, userID uint, username, role string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(config.AccessTokenExp)

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(config *JWTConfig, userID uint, username, role string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Type:     "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshTokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Secret))
}

// ValidateRefreshToken 验证刷新令牌
func ValidateRefreshToken(config *JWTConfig, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(config.Secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenMalformed
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, jwt.ErrTokenMalformed
	}

	// 验证token类型
	if claims.Type != "refresh" {
		return nil, jwt.ErrTokenMalformed
	}

	return claims, nil
}

// GenerateToken 生成JWT令牌 (保持向后兼容)
func GenerateToken(config *JWTConfig, userID uint, username, role string) (string, error) {
	return GenerateAccessToken(config, userID, username, role)
}
