package models

// API 请求结构体定义

// StepConfigAPI 步骤配置 API 请求结构
type StepConfigAPI struct {
	Step          string   `json:"step" validate:"required,min=1,max=100" example:"submit"`                    // 步骤名称
	State         string   `json:"state" validate:"omitempty,max=100" example:"pending"`                       // 步骤状态
	SignType      string   `json:"sign_type" validate:"required,oneof=jointly serial anyone" example:"serial"` // 签批类型
	JointSignRate float32  `json:"joint_sign_rate" validate:"omitempty,min=0,max=1" example:"0.5"`             // 联合签批比例
	OperatorList  []string `json:"operators" validate:"required,min=1" example:"manager1,manager2"`            // 操作人列表
	NextStepList  []string `json:"next_steps" validate:"required,min=1" example:"review,approve"`              // 下一步骤列表
}

// CreateTicketTemplateAPI 创建工单模板 API 请求
type CreateTicketTemplateAPI struct {
	Name        string           `json:"name" validate:"required,min=1,max=100" example:"请假申请模板"`       // 模板名称
	Memo        string           `json:"memo" validate:"omitempty,max=1000" example:"请假申请流程"`           // 模板备注
	Version     string           `json:"version" validate:"required,min=1,max=50" example:"1.0"`        // 版本号
	Creator     string           `json:"creator" validate:"required,min=1,max=50" example:"管理员"`        // 创建者
	StartStep   string           `json:"start_step" validate:"required,min=1,max=100" example:"submit"` // 起始步骤
	EndStepList []string         `json:"end_step" validate:"required,min=1" example:"approve,reject"`   // 结束步骤列表
	ConfigList  []*StepConfigAPI `json:"config" validate:"required,min=1"`                              // 步骤配置列表
}

// UpdateTicketTemplateAPI 更新工单模板 API 请求
type UpdateTicketTemplateAPI struct {
	ID          int              `json:"id" validate:"required,gt=0"`
	Name        *string          `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Memo        *string          `json:"memo,omitempty" validate:"omitempty,max=1000"`
	Creator     *string          `json:"creator,omitempty" validate:"omitempty,min=1,max=50"`
	StartStep   *string          `json:"start_step,omitempty" validate:"omitempty,min=1,max=100"`
	EndStepList []string         `json:"end_step,omitempty" validate:"omitempty,min=1"`
	ConfigList  []*StepConfigAPI `json:"config,omitempty" validate:"omitempty,min=1"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"admin"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"123456"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int64  `json:"expires_in" example:"86400"`
	TokenType    string `json:"token_type" example:"Bearer"`
}

// RefreshTokenRequest 刷新token请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RefreshTokenResponse 刷新token响应
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	ExpiresIn    int64  `json:"expires_in" example:"86400"`
	TokenType    string `json:"token_type" example:"Bearer"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50" example:"newuser"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"123456"`
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Nickname string `json:"nickname" validate:"omitempty,max=50" example:"新用户"`
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	User User `json:"user"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error" example:"错误信息"`
	Details string `json:"details,omitempty" example:"详细错误信息"`
}
