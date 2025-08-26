package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"gorm.io/gorm"
)

// TicketTemplateManager 工单模板数据库管理器
type TicketTemplateManager struct {
	db *gorm.DB
}

// NewTicketTemplateManager 创建工单模板管理器实例
func NewTicketTemplateManager(db *gorm.DB) *TicketTemplateManager {
	return &TicketTemplateManager{db: db}
}

// CreateTicketTemplate 创建工单模板
func (m *TicketTemplateManager) CreateTicketTemplate(ctx context.Context, template *models.TicketTemplate, endSteps []models.TemplateEndStep, stepConfigs []models.StepConfig) error {
	// 开启事务
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建工单模板
	if err := tx.Create(template).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create ticket template: %w", err)
	}

	// 创建结束步骤配置
	for i := range endSteps {
		endSteps[i].TicketTemplateID = template.ID
		if err := tx.Create(&endSteps[i]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create template end step: %w", err)
		}
	}

	// 创建步骤配置
	for i := range stepConfigs {
		stepConfigs[i].TicketTemplateID = template.ID
		if err := tx.Create(&stepConfigs[i]).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create step config: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTicketTemplateByID 根据ID获取工单模板
func (m *TicketTemplateManager) GetTicketTemplateByID(ctx context.Context, id int) (*models.TicketTemplate, error) {
	var template models.TicketTemplate
	if err := m.db.Preload("EndSteps").Preload("StepConfigs").First(&template, id).Error; err != nil {
		return nil, fmt.Errorf("ticket template not found: %v", err)
	}
	return &template, nil
}

// UpdateTicketTemplate 更新工单模板
func (m *TicketTemplateManager) UpdateTicketTemplate(ctx context.Context, id int, updates map[string]interface{}) error {
	// 开启事务
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查模板是否存在
	var template models.TicketTemplate
	if err := tx.First(&template, id).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ticket template not found: %v", err)
	}

	// 检查是否为内置模板
	if template.Builtin {
		tx.Rollback()
		return fmt.Errorf("cannot update builtin template")
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	// 更新模板信息
	if err := tx.Model(&template).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update ticket template: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteTicketTemplate 删除工单模板
func (m *TicketTemplateManager) DeleteTicketTemplate(ctx context.Context, id int) error {
	// 开启事务
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查模板是否存在
	var template models.TicketTemplate
	if err := tx.First(&template, id).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ticket template not found: %v", err)
	}

	// 检查是否为内置模板
	if template.Builtin {
		tx.Rollback()
		return fmt.Errorf("cannot delete builtin template")
	}

	// 检查是否有工单使用此模板
	var count int64
	if err := tx.Model(&models.Ticket{}).Where("template_id = ?", id).Count(&count).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to check template usage: %w", err)
	}

	if count > 0 {
		tx.Rollback()
		return fmt.Errorf("template is in use by %d tickets", count)
	}

	// 删除相关配置
	if err := tx.Where("ticket_template_id = ?", id).Delete(&models.TemplateEndStep{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete template end steps: %w", err)
	}

	if err := tx.Where("ticket_template_id = ?", id).Delete(&models.StepConfig{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete step configs: %w", err)
	}

	// 软删除模板
	if err := tx.Delete(&template).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete ticket template: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ListTicketTemplates 查询工单模板列表
func (m *TicketTemplateManager) ListTicketTemplates(ctx context.Context, filters map[string]interface{}, page, size int) ([]*models.TicketTemplate, int64, error) {
	// 构建查询条件
	query := m.db.Model(&models.TicketTemplate{})

	// 应用筛选条件
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if creator, ok := filters["creator"].(string); ok && creator != "" {
		query = query.Where("creator = ?", creator)
	}
	if version, ok := filters["version"].(string); ok && version != "" {
		query = query.Where("version = ?", version)
	}
	if builtin, ok := filters["builtin"].(bool); ok {
		query = query.Where("builtin = ?", builtin)
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count ticket templates: %w", err)
	}

	// 分页查询
	var templates []*models.TicketTemplate
	offset := (page - 1) * size
	if err := query.Preload("EndSteps").Preload("StepConfigs").
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&templates).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query ticket templates: %w", err)
	}

	return templates, total, nil
}

// GetTemplateByUID 根据UID获取工单模板
func (m *TicketTemplateManager) GetTemplateByUID(ctx context.Context, uid string) (*models.TicketTemplate, error) {
	var template models.TicketTemplate
	if err := m.db.Where("uid = ?", uid).Preload("EndSteps").Preload("StepConfigs").First(&template).Error; err != nil {
		return nil, fmt.Errorf("ticket template not found: %v", err)
	}
	return &template, nil
}

// GetTemplateEndSteps 获取模板结束步骤
func (m *TicketTemplateManager) GetTemplateEndSteps(ctx context.Context, templateID int) ([]models.TemplateEndStep, error) {
	var endSteps []models.TemplateEndStep
	if err := m.db.Where("ticket_template_id = ?", templateID).Find(&endSteps).Error; err != nil {
		return nil, fmt.Errorf("failed to get template end steps: %w", err)
	}
	return endSteps, nil
}

// GetTemplateStepConfigs 获取模板步骤配置
func (m *TicketTemplateManager) GetTemplateStepConfigs(ctx context.Context, templateID int) ([]models.StepConfig, error) {
	var stepConfigs []models.StepConfig
	if err := m.db.Where("ticket_template_id = ?", templateID).Find(&stepConfigs).Error; err != nil {
		return nil, fmt.Errorf("failed to delete step configs: %w", err)
	}
	return stepConfigs, nil
}

// CheckTemplateExists 检查模板是否存在
func (m *TicketTemplateManager) CheckTemplateExists(ctx context.Context, id int) (bool, error) {
	var count int64
	if err := m.db.Model(&models.TicketTemplate{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check template existence: %w", err)
	}
	return count > 0, nil
}
