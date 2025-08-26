package ticket_template

import (
	"context"
	"fmt"

	"github.com/FatWang1/open_ticket_system/internal/helper"
	"github.com/FatWang1/open_ticket_system/internal/manager"
	"github.com/FatWang1/open_ticket_system/internal/models"
	"gorm.io/gorm"
)

// TicketTemplateService 工单模板服务接口
type TicketTemplateService interface {
	CreateTicketTemplate(ctx context.Context, input *models.CreateTicketTemplateRequest) (*models.CreateTicketTemplateResponse, error)
	GetTicketTemplateByID(ctx context.Context, id int) (*models.TicketTemplate, error)
	UpdateTicketTemplate(ctx context.Context, input *models.UpdateTicketTemplateRequest) (*models.UpdateTicketTemplateResponse, error)
	DeleteTicketTemplate(ctx context.Context, input *models.DeleteTicketTemplateRequest) (*models.DeleteTicketTemplateResponse, error)
	ListTicketTemplates(ctx context.Context, input *models.ListTicketTemplateRequest) (*models.ListTicketTemplateResponse, error)
}

// ticketTemplateService 工单模板服务实现
type ticketTemplateService struct {
	templateManager        *manager.TicketTemplateManager
	punchedTapeIntegration *manager.PunchedTapeIntegration
	db                     *gorm.DB
}

// NewTicketTemplateService 创建工单模板服务实例
func NewTicketTemplateService(db *gorm.DB) TicketTemplateService {
	if db == nil {
		// 返回模拟服务
		return &ticketTemplateService{
			templateManager:        nil,
			punchedTapeIntegration: nil,
			db:                     nil,
		}
	}

	return &ticketTemplateService{
		templateManager:        manager.NewTicketTemplateManager(db),
		punchedTapeIntegration: manager.NewPunchedTapeIntegration(db),
		db:                     db,
	}
}

// CreateTicketTemplate 创建工单模板
func (s *ticketTemplateService) CreateTicketTemplate(ctx context.Context, input *models.CreateTicketTemplateRequest) (*models.CreateTicketTemplateResponse, error) {
	if s.templateManager == nil {
		// 模拟创建逻辑
		return &models.CreateTicketTemplateResponse{ID: 1}, nil
	}

	// 使用punched-tape集成层构建模板
	template, err := s.punchedTapeIntegration.BuildTemplateWithPunchedTape(
		ctx,
		helper.GenerateUID(),
		input.StartStep,
		input.Name,
		input.Memo,
		input.Version,
		input.Creator,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build template with punched-tape: %w", err)
	}

	return &models.CreateTicketTemplateResponse{ID: int(template.ID)}, nil
}

// GetTicketTemplateByID 根据ID获取工单模板
func (s *ticketTemplateService) GetTicketTemplateByID(ctx context.Context, id int) (*models.TicketTemplate, error) {
	if s.templateManager == nil {
		// 模拟数据，实际使用时应该查询数据库
		return &models.TicketTemplate{
			Uid:       fmt.Sprintf("template-%d", id),
			StartStep: "step1",
			Builtin:   false,
		}, nil
	}

	return s.templateManager.GetTicketTemplateByID(ctx, id)
}

// UpdateTicketTemplate 更新工单模板
func (s *ticketTemplateService) UpdateTicketTemplate(ctx context.Context, input *models.UpdateTicketTemplateRequest) (*models.UpdateTicketTemplateResponse, error) {
	if s.templateManager == nil {
		// 模拟更新逻辑
		return &models.UpdateTicketTemplateResponse{ID: input.ID}, nil
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Memo != nil {
		updates["memo"] = *input.Memo
	}
	if input.StartStep != nil {
		updates["start_step"] = *input.StartStep
	}

	// 使用数据库管理器更新工单模板
	if err := s.templateManager.UpdateTicketTemplate(ctx, input.ID, updates); err != nil {
		return nil, fmt.Errorf("failed to update ticket template: %w", err)
	}

	return &models.UpdateTicketTemplateResponse{ID: input.ID}, nil
}

// DeleteTicketTemplate 删除工单模板
func (s *ticketTemplateService) DeleteTicketTemplate(ctx context.Context, input *models.DeleteTicketTemplateRequest) (*models.DeleteTicketTemplateResponse, error) {
	if s.templateManager == nil {
		// 模拟删除逻辑
		return &models.DeleteTicketTemplateResponse{ID: input.ID}, nil
	}

	// 使用数据库管理器删除工单模板
	if err := s.templateManager.DeleteTicketTemplate(ctx, input.ID); err != nil {
		return nil, fmt.Errorf("failed to delete ticket template: %w", err)
	}

	return &models.DeleteTicketTemplateResponse{ID: input.ID}, nil
}

// ListTicketTemplates 查询工单模板列表
func (s *ticketTemplateService) ListTicketTemplates(ctx context.Context, input *models.ListTicketTemplateRequest) (*models.ListTicketTemplateResponse, error) {
	if s.templateManager == nil {
		// 模拟查询逻辑
		return &models.ListTicketTemplateResponse{
			Total: 0,
			List:  []*models.TicketTemplateResponse{},
		}, nil
	}

	// 构建查询过滤器
	filters := helper.BuildTemplateFilters(input.Name, input.Creator, nil, input.Builtin)

	// 使用数据库管理器查询工单模板列表
	templates, total, err := s.templateManager.ListTicketTemplates(ctx, filters, input.Page, input.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to list ticket templates: %w", err)
	}

	// 转换为响应模型
	templateResponses := make([]*models.TicketTemplateResponse, len(templates))
	for i, template := range templates {
		templateResponses[i] = s.convertTemplateToResponse(template)
	}

	return &models.ListTicketTemplateResponse{
		Total: total,
		List:  templateResponses,
	}, nil
}

// convertTemplateToResponse 将TicketTemplate转换为TicketTemplateResponse
func (s *ticketTemplateService) convertTemplateToResponse(template *models.TicketTemplate) *models.TicketTemplateResponse {
	// 获取结束步骤列表
	endSteps := make([]string, 0)
	for _, step := range template.EndSteps {
		endSteps = append(endSteps, step.EndStep)
	}

	// 获取步骤配置列表
	stepConfigs := make([]models.StepConfigResponse, 0)
	for _, config := range template.StepConfigs {
		// 获取操作人列表
		operators := make([]string, 0)
		for _, op := range config.Operators {
			operators = append(operators, op.Operator)
		}

		// 获取下一步骤列表
		nextSteps := make([]models.NextStepResponse, 0)
		for _, next := range config.NextSteps {
			nextSteps = append(nextSteps, models.NextStepResponse{
				ID:        next.ID,
				ToStep:    next.ToStep,
				Operation: next.Operation,
			})
		}

		stepConfigs = append(stepConfigs, models.StepConfigResponse{
			ID:            config.ID,
			Step:          config.Step,
			State:         config.State,
			SignType:      config.SignType,
			JointSignRate: config.JointSignRate,
			Operators:     operators,
			NextSteps:     nextSteps,
		})
	}

	return &models.TicketTemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Memo:        template.Memo,
		Version:     template.Version,
		Creator:     template.Creator,
		Uid:         template.Uid,
		StartStep:   template.StartStep,
		Builtin:     template.Builtin,
		CreatedAt:   template.CreatedAt,
		UpdatedAt:   template.UpdatedAt,
		EndSteps:    endSteps,
		StepConfigs: stepConfigs,
	}
}
