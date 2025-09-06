package validator

import (
	"errors"
	"regexp"
	"strings"

	"github.com/FatWang1/open_ticket_system/internal/models"
)

// ValidateLoginRequest 验证登录请求
func ValidateLoginRequest(req *models.LoginRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return errors.New("用户名不能为空")
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("用户名长度必须在3-50个字符之间")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("密码不能为空")
	}

	if len(req.Password) < 6 || len(req.Password) > 100 {
		return errors.New("密码长度必须在6-100个字符之间")
	}

	return nil
}

// ValidateRegisterRequest 验证注册请求
func ValidateRegisterRequest(req *models.RegisterRequest) error {
	if strings.TrimSpace(req.Username) == "" {
		return errors.New("用户名不能为空")
	}

	if len(req.Username) < 3 || len(req.Username) > 50 {
		return errors.New("用户名长度必须在3-50个字符之间")
	}

	// 验证用户名格式（只允许字母、数字、下划线）
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(req.Username) {
		return errors.New("用户名只能包含字母、数字和下划线")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("密码不能为空")
	}

	if len(req.Password) < 6 || len(req.Password) > 100 {
		return errors.New("密码长度必须在6-100个字符之间")
	}

	// 验证密码强度（至少包含字母和数字）
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(req.Password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(req.Password)
	if !hasLetter || !hasNumber {
		return errors.New("密码必须包含至少一个字母和一个数字")
	}

	if strings.TrimSpace(req.Email) == "" {
		return errors.New("邮箱不能为空")
	}

	// 验证邮箱格式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("邮箱格式不正确")
	}

	if len(req.Email) > 100 {
		return errors.New("邮箱长度不能超过100个字符")
	}

	if req.Nickname != "" && len(req.Nickname) > 50 {
		return errors.New("昵称长度不能超过50个字符")
	}

	return nil
}

// ValidateRefreshTokenRequest 验证刷新token请求
func ValidateRefreshTokenRequest(req *models.RefreshTokenRequest) error {
	if strings.TrimSpace(req.RefreshToken) == "" {
		return errors.New("刷新令牌不能为空")
	}

	return nil
}
