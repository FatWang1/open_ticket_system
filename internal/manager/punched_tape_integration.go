package manager

import (
	"context"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/FatWang1/open_ticket_system/internal/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// PunchedTapeIntegration punched-tape集成
type PunchedTapeIntegration struct {
	db *gorm.DB
}

// NewPunchedTapeIntegration 创建新的punched-tape集成实例
func NewPunchedTapeIntegration(db *gorm.DB) *PunchedTapeIntegration {
	return &PunchedTapeIntegration{db: db}
}

// BuildTemplateWithPunchedTape 使用punched-tape构建模板
func (pti *PunchedTapeIntegration) BuildTemplateWithPunchedTape(ctx context.Context, name, memo, version, creator, uid, startStep string, endStepList []string, configList []*models.StepConfigAPI) (*models.TicketTemplate, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Building template with punched-tape - name: %s, uid: %s", name, uid)

	// 创建内部模板模型
	template := &models.TicketTemplate{
		Name:      name,
		Memo:      memo,
		Version:   version,
		Creator:   creator,
		Uid:       uid,
		StartStep: startStep,
		Builtin:   false,
	}

	// 验证模板数据
	if err := pti.validateTemplateData(template, endStepList, configList); err != nil {
		logger.Printf("[ERROR] Template validation failed - name: %s, error: %v", name, err)
		return nil, errors.Wrap(err, "template validation failed")
	}

	logger.Printf("[INFO] Successfully built template with punched-tape - name: %s, uid: %s", name, uid)
	return template, nil
}

// CreateTemplateFromPunchedTape 从punched-tape模板创建内部模板
func (pti *PunchedTapeIntegration) CreateTemplateFromPunchedTape(ctx context.Context, name, memo, version, creator, uid, startStep string) (*models.TicketTemplate, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Creating template from punched-tape - uid: %s", uid)

	template := &models.TicketTemplate{
		Name:      name,
		Memo:      memo,
		Version:   version,
		Creator:   creator,
		Uid:       uid,
		StartStep: startStep,
		Builtin:   false,
	}

	logger.Printf("[INFO] Successfully created template from punched-tape - uid: %s", uid)
	return template, nil
}

// IsEndStep 检查是否为结束步骤
func (pti *PunchedTapeIntegration) IsEndStep(ctx context.Context, template *models.TicketTemplate, step string) (bool, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Checking if step is end step - template_id: %d, step: %s", template.ID, step)

	// 检查是否为结束步骤
	var endSteps []models.TemplateEndStep
	if err := pti.db.Where("ticket_template_id = ?", template.ID).Find(&endSteps).Error; err != nil {
		logger.Printf("[ERROR] Failed to get template end steps - template_id: %d, error: %v", template.ID, err)
		return false, errors.Wrap(err, "failed to get template end steps")
	}

	for _, endStep := range endSteps {
		if endStep.EndStep == step {
			logger.Printf("[INFO] Step is end step - template_id: %d, step: %s", template.ID, step)
			return true, nil
		}
	}

	logger.Printf("[INFO] Step is not end step - template_id: %d, step: %s", template.ID, step)
	return false, nil
}

// getStepConfig 获取步骤配置
func (pti *PunchedTapeIntegration) getStepConfig(template *models.TicketTemplate, step string) *models.StepConfigDB {
	var stepConfigs []models.StepConfigDB
	if err := pti.db.Where("ticket_template_id = ?", template.ID).Find(&stepConfigs).Error; err != nil {
		return nil
	}

	for _, config := range stepConfigs {
		if config.Step == step {
			return &config
		}
	}
	return nil
}

// getStepOperators 获取步骤操作人
func (pti *PunchedTapeIntegration) getStepOperators(stepConfig *models.StepConfigDB) []string {
	var operators []models.StepOperator
	if err := pti.db.Where("step_config_id = ?", stepConfig.ID).Find(&operators).Error; err != nil {
		return nil
	}

	operatorList := make([]string, 0, len(operators))
	for _, operator := range operators {
		operatorList = append(operatorList, operator.Operator)
	}
	return operatorList
}

// ValidateTemplate 验证模板
func (pti *PunchedTapeIntegration) ValidateTemplate(ctx context.Context, template *models.TicketTemplate) error {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Validating template - uid: %s", template.Uid)

	// 检查起始步骤是否存在
	stepConfig := pti.getStepConfig(template, template.StartStep)
	if stepConfig == nil {
		logger.Printf("[ERROR] Start step not found in template - uid: %s, start_step: %s", template.Uid, template.StartStep)
		return errors.New("start step not found in template")
	}

	// 检查所有步骤是否可达
	if err := pti.checkStepsReachable(ctx, template); err != nil {
		logger.Printf("[ERROR] Template validation failed - uid: %s, error: %v", template.Uid, err)
		return errors.Wrap(err, "template validation failed")
	}

	logger.Printf("[INFO] Template validation successful - uid: %s", template.Uid)
	return nil
}

// validateTemplateData 验证模板数据
func (pti *PunchedTapeIntegration) validateTemplateData(template *models.TicketTemplate, endStepList []string, configList []*models.StepConfigAPI) error {
	// 验证基本信息
	if template.Name == "" {
		return errors.New("template name is required")
	}
	if template.StartStep == "" {
		return errors.New("start step is required")
	}
	if len(endStepList) == 0 {
		return errors.New("at least one end step is required")
	}
	if len(configList) == 0 {
		return errors.New("at least one step config is required")
	}

	// 验证步骤配置
	for _, config := range configList {
		if config.Step == "" {
			return errors.New("step name is required")
		}
		if len(config.OperatorList) == 0 {
			return errors.New("at least one operator is required for step: " + config.Step)
		}
		if len(config.NextStepList) == 0 {
			return errors.New("at least one next step is required for step: " + config.Step)
		}
	}

	return nil
}

// checkStepsReachable 检查所有步骤是否可达
func (pti *PunchedTapeIntegration) checkStepsReachable(ctx context.Context, template *models.TicketTemplate) error {
	// 这里可以实现更复杂的可达性检查逻辑
	// 目前简单返回nil表示检查通过
	return nil
}
