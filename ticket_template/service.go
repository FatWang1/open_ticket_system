package ticket_template

import (
	"context"
	"fmt"
	"time"

	"github.com/FatWang1/open_ticket_system/internal/manager"
	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/FatWang1/open_ticket_system/internal/utils"
	"github.com/pkg/errors"
)

// TicketTemplateService 工单模板服务
type TicketTemplateService struct {
	templateManager        *manager.TicketTemplateManager
	punchedTapeIntegration *manager.PunchedTapeIntegration
}

// NewTicketTemplateService 创建新的工单模板服务
func NewTicketTemplateService(templateManager *manager.TicketTemplateManager, punchedTapeIntegration *manager.PunchedTapeIntegration) *TicketTemplateService {
	return &TicketTemplateService{
		templateManager:        templateManager,
		punchedTapeIntegration: punchedTapeIntegration,
	}
}

// CreateTicketTemplate 创建工单模板
func (s *TicketTemplateService) CreateTicketTemplate(ctx context.Context, input *models.CreateTicketTemplateAPI) (*models.TicketTemplateResponse, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Creating ticket template - name: %s", input.Name)

	// 生成UID
	uid := fmt.Sprintf("template-%d", time.Now().Unix())

	// 构建punched-tape模板
	template, err := s.punchedTapeIntegration.BuildTemplateWithPunchedTape(ctx, input.Name, input.Memo, input.Version, input.Creator, uid, input.StartStep, input.EndStepList, input.ConfigList)
	if err != nil {
		logger.Printf("[ERROR] Failed to build template with punched-tape - name: %s, error: %v", input.Name, err)
		return nil, errors.Wrap(err, "failed to build template with punched-tape")
	}

	// 转换EndStepList为TemplateEndStep
	endSteps := make([]models.TemplateEndStep, 0, len(input.EndStepList))
	for _, endStep := range input.EndStepList {
		endSteps = append(endSteps, models.TemplateEndStep{
			EndStep: endStep,
		})
	}

	// 转换ConfigList为StepConfigDB
	stepConfigs := make([]models.StepConfigDB, 0, len(input.ConfigList))
	for _, config := range input.ConfigList {
		stepConfig := models.StepConfigDB{
			Step:          config.Step,
			State:         config.State,
			SignType:      config.SignType,
			JointSignRate: config.JointSignRate,
		}
		stepConfigs = append(stepConfigs, stepConfig)
	}

	// 创建模板
	if err := s.templateManager.CreateTicketTemplate(ctx, template, endSteps, stepConfigs); err != nil {
		logger.Printf("[ERROR] Failed to create ticket template - name: %s, error: %v", input.Name, err)
		return nil, errors.Wrap(err, "failed to create ticket template")
	}

	// 转换为响应格式
	response := s.convertTemplateToResponse(template)

	logger.Printf("[INFO] Successfully created ticket template - name: %s, id: %d", input.Name, template.ID)
	return response, nil
}

// GetTicketTemplateByID 根据ID获取工单模板
func (s *TicketTemplateService) GetTicketTemplateByID(ctx context.Context, id int) (*models.TicketTemplateResponse, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Getting ticket template by ID - id: %d", id)

	template, err := s.templateManager.GetTicketTemplateByID(ctx, id)
	if err != nil {
		logger.Printf("[ERROR] Failed to get ticket template by ID - id: %d, error: %v", id, err)
		return nil, errors.Wrap(err, "failed to get ticket template")
	}

	response := s.convertTemplateToResponse(template)

	logger.Printf("[INFO] Successfully retrieved ticket template - id: %d, name: %s", id, template.Name)
	return response, nil
}

// GetTicketTemplateByUID 根据UID获取工单模板
func (s *TicketTemplateService) GetTicketTemplateByUID(ctx context.Context, uid string) (*models.TicketTemplateResponse, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Getting ticket template by UID - uid: %s", uid)

	template, err := s.templateManager.GetTemplateByUID(ctx, uid)
	if err != nil {
		logger.Printf("[ERROR] Failed to get ticket template by UID - uid: %s, error: %v", uid, err)
		return nil, errors.Wrap(err, "failed to get ticket template")
	}

	response := s.convertTemplateToResponse(template)

	logger.Printf("[INFO] Successfully retrieved ticket template - uid: %s, name: %s", uid, template.Name)
	return response, nil
}

// ListTicketTemplates 获取工单模板列表
func (s *TicketTemplateService) ListTicketTemplates(ctx context.Context, page, size int) ([]models.TicketTemplateResponse, int64, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Listing ticket templates - page: %d, size: %d", page, size)

	templates, total, err := s.templateManager.ListTicketTemplates(ctx, page, size)
	if err != nil {
		logger.Printf("[ERROR] Failed to list ticket templates - error: %v", err)
		return nil, 0, errors.Wrap(err, "failed to list ticket templates")
	}

	// 转换为响应格式
	responses := make([]models.TicketTemplateResponse, 0, len(templates))
	for _, template := range templates {
		response := s.convertTemplateToResponse(&template)
		responses = append(responses, *response)
	}

	logger.Printf("[INFO] Successfully listed ticket templates - count: %d, total: %d", len(templates), total)
	return responses, total, nil
}

// UpdateTicketTemplate 更新工单模板
func (s *TicketTemplateService) UpdateTicketTemplate(ctx context.Context, id int, input *models.UpdateTicketTemplateAPI) (*models.TicketTemplateResponse, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Updating ticket template - id: %d", id)

	// 获取现有模板
	template, err := s.templateManager.GetTicketTemplateByID(ctx, id)
	if err != nil {
		logger.Printf("[ERROR] Failed to get ticket template for update - id: %d, error: %v", id, err)
		return nil, errors.Wrap(err, "failed to get ticket template")
	}

	// 更新字段
	if input.Name != nil {
		template.Name = *input.Name
	}
	if input.Memo != nil {
		template.Memo = *input.Memo
	}
	if input.StartStep != nil {
		template.StartStep = *input.StartStep
	}

	// 更新模板
	if err := s.templateManager.UpdateTicketTemplate(ctx, template); err != nil {
		logger.Printf("[ERROR] Failed to update ticket template - id: %d, error: %v", id, err)
		return nil, errors.Wrap(err, "failed to update ticket template")
	}

	response := s.convertTemplateToResponse(template)

	logger.Printf("[INFO] Successfully updated ticket template - id: %d", id)
	return response, nil
}

// DeleteTicketTemplate 删除工单模板
func (s *TicketTemplateService) DeleteTicketTemplate(ctx context.Context, id int) error {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Deleting ticket template - id: %d", id)

	if err := s.templateManager.DeleteTicketTemplate(ctx, id); err != nil {
		logger.Printf("[ERROR] Failed to delete ticket template - id: %d, error: %v", id, err)
		return errors.Wrap(err, "failed to delete ticket template")
	}

	logger.Printf("[INFO] Successfully deleted ticket template - id: %d", id)
	return nil
}

// convertTemplateToResponse 转换模板为响应格式
func (s *TicketTemplateService) convertTemplateToResponse(template *models.TicketTemplate) *models.TicketTemplateResponse {
	// 获取结束步骤列表
	endSteps := make([]string, 0)
	if s.templateManager != nil {
		if endStepModels, err := s.templateManager.GetTemplateEndSteps(context.Background(), int(template.ID)); err == nil {
			for _, step := range endStepModels {
				endSteps = append(endSteps, step.EndStep)
			}
		}
	}

	// 获取步骤配置列表
	stepConfigs := make([]models.StepConfigResponse, 0)
	if s.templateManager != nil {
		if stepConfigModels, err := s.templateManager.GetTemplateStepConfigs(context.Background(), int(template.ID)); err == nil {
			for _, config := range stepConfigModels {
				// 获取操作人列表
				operators := make([]string, 0)
				if operatorModels, err := s.templateManager.GetStepOperators(context.Background(), int(config.ID)); err == nil {
					for _, op := range operatorModels {
						operators = append(operators, op.Operator)
					}
				}

				// 获取下一步骤列表
				nextSteps := make([]models.NextStepResponse, 0)
				if nextStepModels, err := s.templateManager.GetNextSteps(context.Background(), int(config.ID)); err == nil {
					for _, next := range nextStepModels {
						nextSteps = append(nextSteps, models.NextStepResponse{
							ID:        next.ID,
							ToStep:    next.ToStep,
							Operation: next.Operation,
						})
					}
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
		}
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
