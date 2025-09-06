package models

import "time"

// TicketResponse 工单响应模型（用于Swagger）
// @Description 工单响应结构
type TicketResponse struct {
	ID            uint      `json:"id" example:"1"`                            // 工单ID
	OrderNum      string    `json:"order_num" example:"TICKET-20250127143000"` // 工单号
	Status        string    `json:"status" example:"running"`                  // 工单状态
	Uid           string    `json:"uid" example:"uid-1234567890"`              // 工单唯一标识
	Step          string    `json:"step" example:"submit"`                     // 当前步骤
	Memo          string    `json:"memo" example:"需要请假3天"`                     // 工单备注
	TemplateID    uint      `json:"template_id" example:"1"`                   // 模板ID
	CreatedAt     time.Time `json:"created_at" example:"2025-01-27T14:30:00Z"` // 创建时间
	UpdatedAt     time.Time `json:"updated_at" example:"2025-01-27T14:30:00Z"` // 更新时间
	Operators     []string  `json:"operators" example:"张三,李四"`                 // 操作人列表
	OperatedUsers []string  `json:"operated_users" example:"王五"`               // 已操作用户列表
}

// TicketTemplateResponse 工单模板响应模型（用于Swagger）
// @Description 工单模板响应结构
type TicketTemplateResponse struct {
	ID          uint                 `json:"id" example:"1"`                            // 模板ID
	Name        string               `json:"name" example:"请假申请模板"`                     // 模板名称
	Memo        string               `json:"memo" example:"请假申请流程"`                     // 模板备注
	Version     string               `json:"version" example:"1.0"`                     // 版本号
	Creator     string               `json:"creator" example:"管理员"`                     // 创建者
	Uid         string               `json:"uid" example:"template-123"`                // 模板唯一标识
	StartStep   string               `json:"start_step" example:"submit"`               // 起始步骤
	Builtin     bool                 `json:"builtin" example:"false"`                   // 是否内置
	CreatedAt   time.Time            `json:"created_at" example:"2025-01-27T14:30:00Z"` // 创建时间
	UpdatedAt   time.Time            `json:"updated_at" example:"2025-01-27T14:30:00Z"` // 更新时间
	EndSteps    []string             `json:"end_steps" example:"approve,reject"`        // 结束步骤列表
	StepConfigs []StepConfigResponse `json:"step_configs"`                              // 步骤配置列表
}

// StepConfigResponse 步骤配置响应模型（用于Swagger）
// @Description 步骤配置响应结构
type StepConfigResponse struct {
	ID            uint               `json:"id" example:"1"`                  // 配置ID
	Step          string             `json:"step" example:"submit"`           // 步骤名
	State         string             `json:"state" example:"pending"`         // 步骤状态
	SignType      string             `json:"sign_type" example:"serial_sign"` // 签名类型
	JointSignRate float32            `json:"joint_sign_rate" example:"0.5"`   // 联合签名比例
	Operators     []string           `json:"operators" example:"张三,李四"`       // 操作人列表
	NextSteps     []NextStepResponse `json:"next_steps"`                      // 下一步骤列表
}

// NextStepResponse 下一步骤响应模型（用于Swagger）
// @Description 下一步骤响应结构
type NextStepResponse struct {
	ID        uint   `json:"id" example:"1"`              // 步骤ID
	ToStep    string `json:"to_step" example:"review"`    // 目标步骤
	Operation string `json:"operation" example:"approve"` // 操作类型
}

// SuccessResponse 成功响应模型（用于Swagger）
// @Description 成功响应结构
type SuccessResponse struct {
	Message string      `json:"message" example:"操作成功"` // 成功消息
	Data    interface{} `json:"data"`                   // 响应数据
}

// PaginationResponse 分页响应模型（用于Swagger）
// @Description 分页响应结构
type PaginationResponse struct {
	Total int         `json:"total" example:"100"` // 总记录数
	Page  int         `json:"page" example:"1"`    // 当前页码
	Size  int         `json:"size" example:"10"`   // 每页大小
	List  interface{} `json:"list"`                // 数据列表
}
