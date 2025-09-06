package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/FatWang1/open_ticket_system/pkg"
	"gorm.io/gorm"
)

// TicketManager 工单数据库管理器
type TicketManager struct {
	db *gorm.DB
}

// NewTicketManager 创建工单管理器实例
func NewTicketManager(db *gorm.DB) *TicketManager {
	return &TicketManager{db: db}
}

// CreateTicket 创建工单
func (m *TicketManager) CreateTicket(ctx context.Context, ticket *models.Ticket, operators []models.TicketOperator) error {
	// 开启事务
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建工单记录
	if err := tx.Create(ticket).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	// 创建工单操作人记录
	for _, op := range operators {
		op.TicketID = ticket.ID
		if err := tx.Create(&op).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create ticket operator: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTicketByID 根据ID获取工单
func (m *TicketManager) GetTicketByID(ctx context.Context, id int) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := m.db.First(&ticket, id).Error; err != nil {
		return nil, fmt.Errorf("ticket not found: %v", err)
	}
	return &ticket, nil
}

// UpdateTicket 更新工单
func (m *TicketManager) UpdateTicket(ctx context.Context, id int, updates map[string]interface{}) error {
	// 开启事务
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查工单是否存在
	var ticket models.Ticket
	if err := tx.First(&ticket, id).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ticket not found: %v", err)
	}

	// 检查工单状态是否允许更新
	if ticket.Status != pkg.Running {
		tx.Rollback()
		return fmt.Errorf("cannot update ticket with status: %s", ticket.Status)
	}

	// 添加更新时间
	updates["updated_at"] = time.Now()

	// 更新工单信息
	if err := tx.Model(&ticket).Updates(updates).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// DeleteTicket 删除工单
func (m *TicketManager) DeleteTicket(ctx context.Context, id int) error {
	// 开启事务
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 检查工单是否存在
	var ticket models.Ticket
	if err := tx.First(&ticket, id).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("ticket not found: %v", err)
	}

	// 检查工单状态是否允许删除
	if ticket.Status == pkg.Passed || ticket.Status == pkg.Rejected {
		tx.Rollback()
		return fmt.Errorf("cannot delete completed ticket")
	}

	// 软删除工单
	if err := tx.Delete(&ticket).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	// 删除相关记录
	if err := tx.Where("ticket_id = ?", id).Delete(&models.TicketOperator{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete ticket operators: %w", err)
	}

	if err := tx.Where("ticket_id = ?", id).Delete(&models.TicketOperatedUser{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete ticket operated users: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetTicketForApproval 获取待审批的工单
func (m *TicketManager) GetTicketForApproval(ctx context.Context, id int) (*models.Ticket, error) {
	var ticket models.Ticket
	if err := m.db.First(&ticket, id).Error; err != nil {
		return nil, fmt.Errorf("ticket not found: %v", err)
	}
	return &ticket, nil
}

// CheckApprovalPermission 检查审批权限
func (m *TicketManager) CheckApprovalPermission(ctx context.Context, ticketID int, approvalUser string) error {
	var operator models.TicketOperator
	if err := m.db.Where("ticket_id = ? AND operator = ?", ticketID, approvalUser).First(&operator).Error; err != nil {
		return fmt.Errorf("approval user not authorized: %v", err)
	}
	return nil
}

// RecordApprovalOperation 记录审批操作
func (m *TicketManager) RecordApprovalOperation(ctx context.Context, operatedUser *models.TicketOperatedUser) error {
	return m.db.Create(operatedUser).Error
}

// UpdateTicketStatus 更新工单状态
func (m *TicketManager) UpdateTicketStatus(ctx context.Context, id int, status string, step string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}
	if step != "" {
		updates["step"] = step
	}

	return m.db.Model(&models.Ticket{}).Where("id = ?", id).Updates(updates).Error
}

// ListTickets 查询工单列表
func (m *TicketManager) ListTickets(ctx context.Context, filters map[string]interface{}, page, size int) ([]*models.Ticket, int64, error) {
	// 构建查询条件
	query := m.db.Model(&models.Ticket{})

	// 应用筛选条件
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("order_num LIKE ?", "%"+name+"%")
	}
	if creator, ok := filters["creator"].(string); ok && creator != "" {
		query = query.Joins("JOIN ticket_operators ON tickets.id = ticket_operators.ticket_id").
			Where("ticket_operators.operator = ?", creator)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if templateID, ok := filters["template_id"].(int); ok && templateID > 0 {
		query = query.Where("template_id = ?", templateID)
	}
	if orderNum, ok := filters["order_num"].(string); ok && orderNum != "" {
		query = query.Where("order_num LIKE ?", "%"+orderNum+"%")
	}

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	// 分页查询
	var tickets []*models.Ticket
	offset := (page - 1) * size
	if err := query.
		Offset(offset).Limit(size).
		Order("created_at DESC").
		Find(&tickets).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to query tickets: %w", err)
	}

	return tickets, total, nil
}

// GetTicketOperators 获取工单操作人
func (m *TicketManager) GetTicketOperators(ctx context.Context, ticketID int) ([]models.TicketOperator, error) {
	var operators []models.TicketOperator
	if err := m.db.Where("ticket_id = ?", ticketID).Find(&operators).Error; err != nil {
		return nil, fmt.Errorf("failed to get ticket operators: %w", err)
	}
	return operators, nil
}

// GetTicketOperatedUsers 获取工单操作记录
func (m *TicketManager) GetTicketOperatedUsers(ctx context.Context, ticketID int) ([]models.TicketOperatedUser, error) {
	var operatedUsers []models.TicketOperatedUser
	if err := m.db.Where("ticket_id = ?", ticketID).Find(&operatedUsers).Error; err != nil {
		return nil, fmt.Errorf("failed to get ticket operated users: %w", err)
	}
	return operatedUsers, nil
}
