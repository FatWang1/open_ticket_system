package models

// CreateTicketRequest 创建工单请求
// @Description 创建工单请求结构
type CreateTicketRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=100" example:"请假申请"` // 工单名称
	Creator    string `json:"creator" validate:"required,min=1,max=50" example:"张三"` // 创建者
	TemplateID uint   `json:"template_id" validate:"required,gt=0" example:"1"`      // 模板ID
	Memo       string `json:"memo" validate:"omitempty,max=1000" example:"需要请假3天"`   // 工单备注
}

// CreateTicketResponse 创建工单响应
// @Description 创建工单响应结构
type CreateTicketResponse struct {
	ID int `json:"id" example:"1"` // 工单ID
}

// UpdateTicketRequest 更新工单请求
type UpdateTicketRequest struct {
	ID   int     `json:"id" validate:"required,gt=0"`
	Memo *string `json:"memo" validate:"omitempty,max=1000"`
}

// UpdateTicketResponse 更新工单响应
type UpdateTicketResponse struct {
	ID int `json:"id"`
}

// DeleteTicketRequest 删除工单请求
type DeleteTicketRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

// DeleteTicketResponse 删除工单响应
type DeleteTicketResponse struct {
	ID int `json:"id"`
}

// ApprovalRequest 工单审批请求
// @Description 工单审批请求结构
type ApprovalRequest struct {
	ID           int     `json:"id" validate:"required,gt=0" example:"1"`                              // 工单ID
	ApprovalUser string  `json:"approval_user" validate:"required,min=1,max=50" example:"李四"`          // 审批人
	Memo         *string `json:"memo,omitempty" validate:"omitempty,max=1000" example:"同意"`            // 审批备注
	Operation    string  `json:"operation" validate:"required,oneof=approve reject" example:"approve"` // 审批操作
	NextStep     string  `json:"next_step" validate:"required,min=1,max=100" example:"step2"`          // 下一步骤
}

// CloseTicketRequest 关闭工单请求
type CloseTicketRequest struct {
	ID       int     `json:"id" validate:"required,gt=0"`
	Memo     *string `json:"memo,omitempty" validate:"omitempty,max=1000"`
	Operator string  `json:"operator" validate:"required,min=1,max=50"`
}

// CloseTicketResponse 关闭工单响应
type CloseTicketResponse struct {
	ID int `json:"id"`
}

// ListTicketRequest 查询工单列表请求
type ListTicketRequest struct {
	Page       int     `json:"page" validate:"min=1,max=1000"`
	Size       int     `json:"size" validate:"min=1,max=100"`
	Name       *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Creator    *string `json:"creator,omitempty" validate:"omitempty,min=1,max=50"`
	Status     *string `json:"status,omitempty" validate:"omitempty,oneof=running passed rejected closed"`
	TemplateID *int    `json:"template_id,omitempty" validate:"omitempty,gt=0"`
	OrderNum   *string `json:"order_num,omitempty" validate:"omitempty,min=1,max=255"`
}

// ListTicketResponse 查询工单列表响应
// @Description 查询工单列表响应结构
type ListTicketResponse struct {
	Total int64             `json:"total" example:"100"` // 总记录数
	List  []*TicketResponse `json:"list"`                // 工单列表
}

// CreateTicketTemplateRequest 创建工单模板请求
// @Description 创建工单模板请求结构
type CreateTicketTemplateRequest struct {
	Name        string           `json:"name" validate:"required,min=1,max=100" example:"请假申请模板"`       // 模板名称
	Memo        string           `json:"memo" validate:"omitempty,max=1000" example:"请假申请流程"`           // 模板备注
	Version     string           `json:"version" validate:"required,min=1,max=50" example:"1.0"`        // 版本号
	Creator     string           `json:"creator" validate:"required,min=1,max=50" example:"管理员"`        // 创建者
	StartStep   string           `json:"start_step" validate:"required,min=1,max=100" example:"submit"` // 起始步骤
	EndStepList []string         `json:"end_step" validate:"required,min=1" example:"approve,reject"`   // 结束步骤列表
	ConfigList  []*StepConfigAPI `json:"config" validate:"required,min=1"`                              // 步骤配置列表
}

// CreateTicketTemplateResponse 创建工单模板响应
type CreateTicketTemplateResponse struct {
	ID int `json:"id"`
}

// GetTicketTemplateRequest 获取工单模板请求
type GetTicketTemplateRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

// UpdateTicketTemplateRequest 更新工单模板请求
type UpdateTicketTemplateRequest struct {
	ID          int              `json:"id" validate:"required,gt=0"`
	Name        *string          `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Memo        *string          `json:"memo,omitempty" validate:"omitempty,max=1000"`
	Creator     *string          `json:"creator,omitempty" validate:"omitempty,min=1,max=50"`
	StartStep   *string          `json:"start_step,omitempty" validate:"omitempty,min=1,max=100"`
	EndStepList []string         `json:"end_step,omitempty" validate:"omitempty,min=1"`
	ConfigList  []*StepConfigAPI `json:"config,omitempty" validate:"omitempty,min=1"`
}

// UpdateTicketTemplateResponse 更新工单模板响应
type UpdateTicketTemplateResponse struct {
	ID int `json:"id"`
}

// ListTicketTemplateRequest 查询工单模板列表请求
type ListTicketTemplateRequest struct {
	Page    int     `json:"page" validate:"min=1,max=1000"`
	Size    int     `json:"size" validate:"min=1,max=100"`
	Name    *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Creator *string `json:"creator,omitempty" validate:"omitempty,min=1,max=50"`
	Builtin *bool   `json:"builtin,omitempty"`
}

// ListTicketTemplateResponse 查询工单模板列表响应
// @Description 查询工单模板列表响应结构
type ListTicketTemplateResponse struct {
	Total int64                     `json:"total" example:"50"` // 总记录数
	List  []*TicketTemplateResponse `json:"list"`               // 模板列表
}

// DeleteTicketTemplateRequest 删除工单模板请求
type DeleteTicketTemplateRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

// DeleteTicketTemplateResponse 删除工单模板响应
type DeleteTicketTemplateResponse struct {
	ID int `json:"id"`
}
