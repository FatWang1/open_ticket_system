package manager

import (
	"context"
	"fmt"

	"github.com/FatWang1/open_ticket_system/internal/helper"
	"github.com/FatWang1/open_ticket_system/internal/models"
	punched_tape "github.com/FatWang1/punched-tape"
	pt_models "github.com/FatWang1/punched-tape/models"
	"gorm.io/gorm"
)

// PunchedTapeIntegration punched-tape集成层
type PunchedTapeIntegration struct {
	db *gorm.DB
}

// NewPunchedTapeIntegration 创建punched-tape集成实例
func NewPunchedTapeIntegration(db *gorm.DB) *PunchedTapeIntegration {
	return &PunchedTapeIntegration{db: db}
}

// CreateTicketFromTemplate 从模板创建工单
func (pti *PunchedTapeIntegration) CreateTicketFromTemplate(ctx context.Context, templateID int, creator string, memo string) (*models.Ticket, error) {
	// 获取模板
	template, err := pti.getTemplateByID(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// 使用punched-tape构建工单
	ticketBuilder := punched_tape.NewTicketBuilder(
		helper.GenerateUID(),
		helper.GenerateOrderNum(),
		template.StartStep,
	)

	// 设置工单属性
	ticketBuilder.SetMemo(memo)
	ticketBuilder.SetStatus(pt_models.Running)

	// 从模板获取当前步骤的操作人
	currentStepConfig := pti.getStepConfig(template, template.StartStep)
	if currentStepConfig != nil {
		// 获取操作人列表
		operators := pti.getStepOperators(currentStepConfig)
		if len(operators) > 0 {
			// 设置操作人列表
			for _, operator := range operators {
				ticketBuilder.AddOperator(operator)
			}
		}
	}

	// 构建工单
	punchedTicket, err := ticketBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build ticket: %w", err)
	}

	// 转换为内部模型
	ticket := &models.Ticket{
		OrderNum:   punchedTicket.OrderNum,
		Status:     punchedTicket.Status,
		Uid:        punchedTicket.Uid,
		Step:       punchedTicket.Step,
		Memo:       punchedTicket.Memo,
		TemplateID: uint(templateID),
	}

	return ticket, nil
}

// GetNextStep 获取下一步骤
func (pti *PunchedTapeIntegration) GetNextStep(ctx context.Context, ticketID int, currentStep string, operation string) (*pt_models.NextStep, error) {
	// 获取工单模板
	ticket, err := pti.getTicketByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	template, err := pti.getTemplateByID(ctx, int(ticket.TemplateID))
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// 获取当前步骤配置
	stepConfig := pti.getStepConfig(template, currentStep)
	if stepConfig == nil {
		return nil, fmt.Errorf("step config not found for step: %s", currentStep)
	}

	// 查找匹配的下一步骤
	for _, next := range stepConfig.NextSteps {
		if next.Operation == operation {
			return &pt_models.NextStep{
				Step:      next.ToStep, // 使用ToStep字段
				Operation: next.Operation,
			}, nil
		}
	}

	return nil, fmt.Errorf("next step not found for operation: %s", operation)
}

// ValidateStepOperation 验证步骤操作
func (pti *PunchedTapeIntegration) ValidateStepOperation(ctx context.Context, ticketID int, step string, operation string) error {
	// 获取工单模板
	ticket, err := pti.getTicketByID(ctx, ticketID)
	if err != nil {
		return fmt.Errorf("failed to get ticket: %w", err)
	}

	template, err := pti.getTemplateByID(ctx, int(ticket.TemplateID))
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// 获取步骤配置
	stepConfig := pti.getStepConfig(template, step)
	if stepConfig == nil {
		return fmt.Errorf("step config not found for step: %s", step)
	}

	// 验证操作是否有效
	validOperation := false
	for _, next := range stepConfig.NextSteps {
		if next.Operation == operation {
			validOperation = true
			break
		}
	}

	if !validOperation {
		return fmt.Errorf("invalid operation %s for step %s", operation, step)
	}

	return nil
}

// GetStepOperators 获取步骤操作人
func (pti *PunchedTapeIntegration) GetStepOperators(ctx context.Context, ticketID int, step string) ([]string, error) {
	// 获取工单模板
	ticket, err := pti.getTicketByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	template, err := pti.getTemplateByID(ctx, int(ticket.TemplateID))
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// 获取步骤配置
	stepConfig := pti.getStepConfig(template, step)
	if stepConfig == nil {
		return nil, fmt.Errorf("step config not found for step: %s", step)
	}

	return pti.getStepOperators(stepConfig), nil
}

// IsEndStep 检查是否为结束步骤
func (pti *PunchedTapeIntegration) IsEndStep(ctx context.Context, ticketID int, step string) (bool, error) {
	// 获取工单模板
	ticket, err := pti.getTicketByID(ctx, ticketID)
	if err != nil {
		return false, fmt.Errorf("failed to get ticket: %w", err)
	}

	template, err := pti.getTemplateByID(ctx, int(ticket.TemplateID))
	if err != nil {
		return false, fmt.Errorf("failed to get template: %w", err)
	}

	// 检查是否为结束步骤
	for _, endStep := range template.EndSteps {
		if endStep.EndStep == step {
			return true, nil
		}
	}

	return false, nil
}

// GetStepDisposal 获取步骤处置方式
func (pti *PunchedTapeIntegration) GetStepDisposal(ctx context.Context, ticketID int, step string) (*pt_models.Disposal, error) {
	// 获取工单模板
	ticket, err := pti.getTicketByID(ctx, ticketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	template, err := pti.getTemplateByID(ctx, int(ticket.TemplateID))
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	// 获取步骤配置
	stepConfig := pti.getStepConfig(template, step)
	if stepConfig == nil {
		return nil, fmt.Errorf("step config not found for step: %s", step)
	}

	// 转换为punched-tape模型
	disposal := &pt_models.Disposal{
		SignType:      stepConfig.SignType,
		JointSignRate: stepConfig.JointSignRate,
	}

	return disposal, nil
}

// CreateTemplateFromPunchedTape 从punched-tape模型创建模板
func (pti *PunchedTapeIntegration) CreateTemplateFromPunchedTape(ctx context.Context, punchedTemplate *pt_models.TicketTemplate, name string, memo string, version string, creator string) (*models.TicketTemplate, error) {
	// 转换为内部模型
	template := &models.TicketTemplate{
		Name:      name,
		Memo:      memo,
		Version:   version,
		Creator:   creator,
		Uid:       punchedTemplate.Uid,
		StartStep: punchedTemplate.StartStep,
		Builtin:   punchedTemplate.Builtin,
	}

	// 转换结束步骤
	var endSteps []models.TemplateEndStep
	for _, endStep := range punchedTemplate.EndStep {
		endSteps = append(endSteps, models.TemplateEndStep{
			EndStep: endStep,
		})
	}

	// 转换步骤配置
	var stepConfigs []models.StepConfig
	for _, config := range punchedTemplate.Config {
		stepConfig := models.StepConfig{
			Step:          config.Step,
			State:         config.State,
			SignType:      config.Disposal.SignType,
			JointSignRate: config.Disposal.JointSignRate,
		}

		// 转换操作人
		for _, operator := range config.Operator {
			stepConfig.Operators = append(stepConfig.Operators, models.StepOperator{
				Operator: operator,
			})
		}

		// 转换下一步骤
		for _, next := range config.Next {
			stepConfig.NextSteps = append(stepConfig.NextSteps, models.NextStep{
				ToStep:    next.Step,
				Operation: next.Operation,
			})
		}

		stepConfigs = append(stepConfigs, stepConfig)
	}

	// 使用模板管理器创建
	templateManager := NewTicketTemplateManager(pti.db)
	if err := templateManager.CreateTicketTemplate(ctx, template, endSteps, stepConfigs); err != nil {
		return nil, fmt.Errorf("failed to create template: %w", err)
	}

	return template, nil
}

// BuildTemplateWithPunchedTape 使用punched-tape构建器创建模板
func (pti *PunchedTapeIntegration) BuildTemplateWithPunchedTape(ctx context.Context, uid string, startStep string, name string, memo string, version string, creator string) (*models.TicketTemplate, error) {
	// 使用punched-tape构建器
	templateBuilder := punched_tape.NewTemplateBuilder(uid, startStep)
	templateBuilder.SetBuiltin(false)

	// 构建punched-tape模板
	punchedTemplate, err := templateBuilder.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build punched-tape template: %w", err)
	}

	// 转换为内部模型并保存
	return pti.CreateTemplateFromPunchedTape(ctx, punchedTemplate, name, memo, version, creator)
}

// 辅助方法

func (pti *PunchedTapeIntegration) getTicketByID(ctx context.Context, id int) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := pti.db.First(&ticket, id).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (pti *PunchedTapeIntegration) getTemplateByID(ctx context.Context, id int) (*models.TicketTemplate, error) {
	var template models.TicketTemplate
	if err := pti.db.Preload("EndSteps").Preload("StepConfigs").First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

func (pti *PunchedTapeIntegration) getStepConfig(template *models.TicketTemplate, step string) *models.StepConfig {
	for _, config := range template.StepConfigs {
		if config.Step == step {
			return &config
		}
	}
	return nil
}

func (pti *PunchedTapeIntegration) getStepOperators(stepConfig *models.StepConfig) []string {
	var operators []string
	for _, op := range stepConfig.Operators {
		operators = append(operators, op.Operator)
	}
	return operators
}
