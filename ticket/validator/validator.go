package validator

import (
	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/go-playground/validator/v10"
)

// CustomValidator 自定义验证器
type CustomValidator struct {
	validate *validator.Validate
}

// NewCustomValidator 创建新的自定义验证器
func NewCustomValidator() *CustomValidator {
	v := validator.New()
	return &CustomValidator{
		validate: v,
	}
}

// Validate 验证结构体
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validate.Struct(i)
}

// ValidateVar 验证单个字段
func (cv *CustomValidator) ValidateVar(field interface{}, tag string) error {
	return cv.validate.Var(field, tag)
}

// 具体的请求验证函数

// ValidateCreateTicketRequest 验证创建工单请求
func ValidateCreateTicketRequest(req *models.CreateTicketRequest) error {
	v := NewCustomValidator()
	return v.Validate(req)
}

// ValidateUpdateTicketRequest 验证更新工单请求
func ValidateUpdateTicketRequest(req *models.UpdateTicketRequest) error {
	v := NewCustomValidator()
	return v.Validate(req)
}

// ValidateApprovalRequest 验证工单审批请求
func ValidateApprovalRequest(req *models.ApprovalRequest) error {
	v := NewCustomValidator()
	return v.Validate(req)
}

// ValidateCloseTicketRequest 验证关闭工单请求
func ValidateCloseTicketRequest(req *models.CloseTicketRequest) error {
	v := NewCustomValidator()
	return v.Validate(req)
}

// ValidateListTicketRequest 验证查询工单列表请求
func ValidateListTicketRequest(req *models.ListTicketRequest) error {
	v := NewCustomValidator()

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	return v.Validate(req)
}
