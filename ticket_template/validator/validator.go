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

// ValidateCreateTicketTemplateRequest 验证创建工单模板请求
func ValidateCreateTicketTemplateRequest(req *models.CreateTicketTemplateRequest) error {
	v := NewCustomValidator()
	return v.Validate(req)
}

// ValidateUpdateTicketTemplateRequest 验证更新工单模板请求
func ValidateUpdateTicketTemplateRequest(req *models.UpdateTicketTemplateRequest) error {
	v := NewCustomValidator()
	return v.Validate(req)
}

// ValidateListTicketTemplateRequest 验证查询工单模板列表请求
func ValidateListTicketTemplateRequest(req *models.ListTicketTemplateRequest) error {
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
