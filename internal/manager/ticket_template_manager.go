package manager

import (
	"context"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/FatWang1/open_ticket_system/internal/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// TicketTemplateManager 工单模板管理器
type TicketTemplateManager struct {
	db *gorm.DB
}

// NewTicketTemplateManager 创建新的工单模板管理器
func NewTicketTemplateManager(db *gorm.DB) *TicketTemplateManager {
	return &TicketTemplateManager{db: db}
}

// CreateTicketTemplate 创建工单模板
func (m *TicketTemplateManager) CreateTicketTemplate(ctx context.Context, template *models.TicketTemplate, endSteps []models.TemplateEndStep, stepConfigs []models.StepConfigDB) error {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Creating ticket template - name: %s, uid: %s", template.Name, template.Uid)

	if err := m.db.Create(template).Error; err != nil {
		logger.Printf("[ERROR] Failed to create ticket template - name: %s, uid: %s, error: %v", template.Name, template.Uid, err)
		return errors.Wrap(err, "failed to create ticket template")
	}

	// 创建结束步骤
	if len(endSteps) > 0 {
		for _, endStep := range endSteps {
			endStep.TicketTemplateID = template.ID
			if err := m.db.Create(&endStep).Error; err != nil {
				logger.Printf("[ERROR] Failed to create template end step - step: %s, error: %v", endStep.EndStep, err)
				return errors.Wrap(err, "failed to create template end step")
			}
		}
	}

	// 创建步骤配置
	if len(stepConfigs) > 0 {
		for _, stepConfig := range stepConfigs {
			stepConfig.TicketTemplateID = template.ID
			if err := m.db.Create(&stepConfig).Error; err != nil {
				logger.Printf("[ERROR] Failed to create step config - step: %s, error: %v", stepConfig.Step, err)
				return errors.Wrap(err, "failed to create step config")
			}

			// 创建操作人
			if len(stepConfig.Operators) > 0 {
				for _, operator := range stepConfig.Operators {
					operator.StepConfigID = stepConfig.ID
					if err := m.db.Create(&operator).Error; err != nil {
						logger.Printf("[ERROR] Failed to create step operator - operator: %s, error: %v", operator.Operator, err)
						return errors.Wrap(err, "failed to create step operator")
					}
				}
			}

			// 创建下一步骤
			if len(stepConfig.NextSteps) > 0 {
				for _, nextStep := range stepConfig.NextSteps {
					nextStep.StepConfigID = stepConfig.ID
					if err := m.db.Create(&nextStep).Error; err != nil {
						logger.Printf("[ERROR] Failed to create next step - to_step: %s, error: %v", nextStep.ToStep, err)
						return errors.Wrap(err, "failed to create next step")
					}
				}
			}
		}
	}

	logger.Printf("[INFO] Successfully created ticket template - name: %s, uid: %s, id: %d", template.Name, template.Uid, template.ID)
	return nil
}

// GetTicketTemplateByID 根据ID获取工单模板
func (m *TicketTemplateManager) GetTicketTemplateByID(ctx context.Context, id int) (*models.TicketTemplate, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Getting ticket template by ID - id: %d", id)

	var template models.TicketTemplate
	if err := m.db.First(&template, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Printf("[WARN] Ticket template not found - id: %d", id)
			return nil, errors.Wrap(err, "ticket template not found")
		}
		logger.Printf("[ERROR] Failed to get ticket template by ID - id: %d, error: %v", id, err)
		return nil, errors.Wrap(err, "failed to get ticket template")
	}

	logger.Printf("[INFO] Successfully retrieved ticket template - id: %d, name: %s", id, template.Name)
	return &template, nil
}

// GetTemplateByUID 根据UID获取工单模板
func (m *TicketTemplateManager) GetTemplateByUID(ctx context.Context, uid string) (*models.TicketTemplate, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Getting ticket template by UID - uid: %s", uid)

	var template models.TicketTemplate
	if err := m.db.Where("uid = ?", uid).First(&template).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Printf("[WARN] Ticket template not found - uid: %s", uid)
			return nil, errors.Wrap(err, "ticket template not found")
		}
		logger.Printf("[ERROR] Failed to get ticket template by UID - uid: %s, error: %v", uid, err)
		return nil, errors.Wrap(err, "failed to get ticket template")
	}

	logger.Printf("[INFO] Successfully retrieved ticket template - uid: %s, name: %s", uid, template.Name)
	return &template, nil
}

// ListTicketTemplates 获取工单模板列表
func (m *TicketTemplateManager) ListTicketTemplates(ctx context.Context, page, size int) ([]models.TicketTemplate, int64, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Listing ticket templates - page: %d, size: %d", page, size)

	var templates []models.TicketTemplate
	var total int64

	// 获取总数
	if err := m.db.Model(&models.TicketTemplate{}).Count(&total).Error; err != nil {
		logger.Printf("[ERROR] Failed to count ticket templates - error: %v", err)
		return nil, 0, errors.Wrap(err, "failed to count ticket templates")
	}

	// 获取分页数据
	offset := (page - 1) * size
	if err := m.db.
		Select("id, created_at, updated_at, deleted_at, name, memo, version, creator, uid, start_step, builtin").
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&templates).Error; err != nil {
		logger.Printf("[ERROR] Failed to query ticket templates - error: %v", err)
		return nil, 0, errors.Wrap(err, "failed to query ticket templates")
	}

	logger.Printf("[INFO] Successfully listed ticket templates - count: %d, total: %d", len(templates), total)
	return templates, total, nil
}

// UpdateTicketTemplate 更新工单模板
func (m *TicketTemplateManager) UpdateTicketTemplate(ctx context.Context, template *models.TicketTemplate) error {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Updating ticket template - id: %d, name: %s", template.ID, template.Name)

	if err := m.db.Save(template).Error; err != nil {
		logger.Printf("[ERROR] Failed to update ticket template - id: %d, error: %v", template.ID, err)
		return errors.Wrap(err, "failed to update ticket template")
	}

	logger.Printf("[INFO] Successfully updated ticket template - id: %d, name: %s", template.ID, template.Name)
	return nil
}

// DeleteTicketTemplate 删除工单模板
func (m *TicketTemplateManager) DeleteTicketTemplate(ctx context.Context, id int) error {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Deleting ticket template - id: %d", id)

	// 删除相关数据
	if err := m.db.Where("ticket_template_id = ?", id).Delete(&models.TemplateEndStep{}).Error; err != nil {
		logger.Printf("[ERROR] Failed to delete template end steps - template_id: %d, error: %v", id, err)
		return errors.Wrap(err, "failed to delete template end steps")
	}

	if err := m.db.Where("ticket_template_id = ?", id).Delete(&models.StepConfigDB{}).Error; err != nil {
		logger.Printf("[ERROR] Failed to delete step configs - template_id: %d, error: %v", id, err)
		return errors.Wrap(err, "failed to delete step configs")
	}

	// 删除模板
	if err := m.db.Delete(&models.TicketTemplate{}, id).Error; err != nil {
		logger.Printf("[ERROR] Failed to delete ticket template - id: %d, error: %v", id, err)
		return errors.Wrap(err, "failed to delete ticket template")
	}

	logger.Printf("[INFO] Successfully deleted ticket template - id: %d", id)
	return nil
}

// GetTemplateEndSteps 获取模板结束步骤
func (m *TicketTemplateManager) GetTemplateEndSteps(ctx context.Context, templateID int) ([]models.TemplateEndStep, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Getting template end steps - template_id: %d", templateID)

	var endSteps []models.TemplateEndStep
	if err := m.db.Where("ticket_template_id = ?", templateID).Find(&endSteps).Error; err != nil {
		logger.Printf("[ERROR] Failed to get template end steps - template_id: %d, error: %v", templateID, err)
		return nil, errors.Wrap(err, "failed to get template end steps")
	}

	logger.Printf("[INFO] Successfully retrieved template end steps - template_id: %d, count: %d", templateID, len(endSteps))
	return endSteps, nil
}

// GetTemplateStepConfigs 获取模板步骤配置
func (m *TicketTemplateManager) GetTemplateStepConfigs(ctx context.Context, templateID int) ([]models.StepConfigDB, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Getting template step configs - template_id: %d", templateID)

	var stepConfigs []models.StepConfigDB
	if err := m.db.Where("ticket_template_id = ?", templateID).Find(&stepConfigs).Error; err != nil {
		logger.Printf("[ERROR] Failed to get template step configs - template_id: %d, error: %v", templateID, err)
		return nil, errors.Wrap(err, "failed to get template step configs")
	}

	logger.Printf("[INFO] Successfully retrieved template step configs - template_id: %d, count: %d", templateID, len(stepConfigs))
	return stepConfigs, nil
}

// GetStepOperators 获取步骤操作人
func (m *TicketTemplateManager) GetStepOperators(ctx context.Context, stepConfigID int) ([]models.StepOperator, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Getting step operators - step_config_id: %d", stepConfigID)

	var operators []models.StepOperator
	if err := m.db.Where("step_config_id = ?", stepConfigID).Find(&operators).Error; err != nil {
		logger.Printf("[ERROR] Failed to get step operators - step_config_id: %d, error: %v", stepConfigID, err)
		return nil, errors.Wrap(err, "failed to get step operators")
	}

	logger.Printf("[INFO] Successfully retrieved step operators - step_config_id: %d, count: %d", stepConfigID, len(operators))
	return operators, nil
}

// GetNextSteps 获取下一步骤配置
func (m *TicketTemplateManager) GetNextSteps(ctx context.Context, stepConfigID int) ([]models.NextStep, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Getting next steps - step_config_id: %d", stepConfigID)

	var nextSteps []models.NextStep
	if err := m.db.Where("step_config_id = ?", stepConfigID).Find(&nextSteps).Error; err != nil {
		logger.Printf("[ERROR] Failed to get next steps - step_config_id: %d, error: %v", stepConfigID, err)
		return nil, errors.Wrap(err, "failed to get next steps")
	}

	logger.Printf("[INFO] Successfully retrieved next steps - step_config_id: %d, count: %d", stepConfigID, len(nextSteps))
	return nextSteps, nil
}
