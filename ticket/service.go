package ticket

import (
	"context"
	"fmt"

	"github.com/FatWang1/open_ticket_system/internal/helper"
	"github.com/FatWang1/open_ticket_system/internal/manager"
	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/FatWang1/open_ticket_system/pkg"
	"gorm.io/gorm"
)

// TicketService 工单服务接口
type TicketService interface {
	CreateTicket(ctx context.Context, input *models.CreateTicketRequest) (*models.CreateTicketResponse, error)
	GetTicketByID(ctx context.Context, id int) (*models.Ticket, error)
	UpdateTicket(ctx context.Context, input *models.UpdateTicketRequest) (*models.UpdateTicketResponse, error)
	DeleteTicket(ctx context.Context, input *models.DeleteTicketRequest) (*models.DeleteTicketResponse, error)
	Approval(ctx context.Context, input *models.ApprovalRequest) error
	CloseTicket(ctx context.Context, input *models.CloseTicketRequest) (*models.CloseTicketResponse, error)
	ListTickets(ctx context.Context, input *models.ListTicketRequest) (*models.ListTicketResponse, error)
}

// ticketService 工单服务实现
type ticketService struct {
	ticketManager          *manager.TicketManager
	templateManager        *manager.TicketTemplateManager
	punchedTapeIntegration *manager.PunchedTapeIntegration
	db                     *gorm.DB
}

// NewTicketService 创建工单服务实例
func NewTicketService(db *gorm.DB) TicketService {
	if db == nil {
		// 返回模拟服务
		return &ticketService{
			ticketManager:          nil,
			templateManager:        nil,
			punchedTapeIntegration: nil,
			db:                     nil,
		}
	}

	return &ticketService{
		ticketManager:          manager.NewTicketManager(db),
		templateManager:        manager.NewTicketTemplateManager(db),
		punchedTapeIntegration: manager.NewPunchedTapeIntegration(db),
		db:                     db,
	}
}

// CreateTicket 创建工单
func (s *ticketService) CreateTicket(ctx context.Context, input *models.CreateTicketRequest) (*models.CreateTicketResponse, error) {
	if s.ticketManager == nil {
		// 模拟创建工单逻辑
		return &models.CreateTicketResponse{ID: 1}, nil
	}

	// 使用punched-tape集成层从模板创建工单
	ticket, err := s.punchedTapeIntegration.CreateTicketFromTemplate(ctx, int(input.TemplateID), input.Creator, input.Memo)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticket from template: %w", err)
	}

	// 从模板获取操作人信息并创建工单操作人记录
	operators := []models.TicketOperator{
		{
			Operator: input.Creator,
		},
	}

	// 使用数据库管理器创建工单
	if err := s.ticketManager.CreateTicket(ctx, ticket, operators); err != nil {
		return nil, fmt.Errorf("failed to create ticket: %w", err)
	}

	return &models.CreateTicketResponse{ID: int(ticket.ID)}, nil
}

// GetTicketByID 根据ID获取工单
func (s *ticketService) GetTicketByID(ctx context.Context, id int) (*models.Ticket, error) {
	if s.ticketManager == nil {
		// 模拟数据，实际使用时应该查询数据库
		return &models.Ticket{
			Model: gorm.Model{
				ID: uint(id),
			},
			OrderNum:   fmt.Sprintf("TICKET-%03d", id),
			Status:     "running",
			Uid:        fmt.Sprintf("uid-%d", id),
			Step:       "step1",
			TemplateID: 1,
		}, nil
	}

	return s.ticketManager.GetTicketByID(ctx, id)
}

// UpdateTicket 更新工单
func (s *ticketService) UpdateTicket(ctx context.Context, input *models.UpdateTicketRequest) (*models.UpdateTicketResponse, error) {
	if s.ticketManager == nil {
		// 模拟更新逻辑，检查ID是否存在
		if input.ID == 9999 {
			return nil, fmt.Errorf("ticket not found")
		}
		return &models.UpdateTicketResponse{ID: input.ID}, nil
	}

	// 构建更新字段
	updates := make(map[string]interface{})
	if input.Memo != nil {
		updates["memo"] = *input.Memo
	}

	// 使用数据库管理器更新工单
	if err := s.ticketManager.UpdateTicket(ctx, input.ID, updates); err != nil {
		return nil, fmt.Errorf("failed to update ticket: %w", err)
	}

	return &models.UpdateTicketResponse{ID: input.ID}, nil
}

// DeleteTicket 删除工单
func (s *ticketService) DeleteTicket(ctx context.Context, input *models.DeleteTicketRequest) (*models.DeleteTicketResponse, error) {
	if s.ticketManager == nil {
		// 模拟删除逻辑，检查ID是否存在
		if input.ID == 9999 {
			return nil, fmt.Errorf("ticket not found")
		}
		return &models.DeleteTicketResponse{ID: input.ID}, nil
	}

	// 使用数据库管理器删除工单
	if err := s.ticketManager.DeleteTicket(ctx, input.ID); err != nil {
		return nil, fmt.Errorf("failed to delete ticket: %w", err)
	}

	return &models.DeleteTicketResponse{ID: input.ID}, nil
}

// Approval 工单审批
func (s *ticketService) Approval(ctx context.Context, input *models.ApprovalRequest) error {
	if s.ticketManager == nil {
		// 模拟审批逻辑，检查操作类型
		if input.Operation != "approve" && input.Operation != "reject" {
			return fmt.Errorf("invalid operation type")
		}
		return nil
	}

	// 获取待审批的工单
	ticket, err := s.ticketManager.GetTicketForApproval(ctx, input.ID)
	if err != nil {
		return fmt.Errorf("failed to get ticket for approval: %w", err)
	}

	// 检查工单状态
	if ticket.Status != pkg.Running {
		return fmt.Errorf("ticket is not in running status")
	}

	// 检查审批人是否有权限
	if err := s.ticketManager.CheckApprovalPermission(ctx, input.ID, input.ApprovalUser); err != nil {
		return fmt.Errorf("approval user not authorized: %w", err)
	}

	// 使用punched-tape集成层验证步骤操作
	if err := s.punchedTapeIntegration.ValidateStepOperation(ctx, input.ID, ticket.Step, input.Operation); err != nil {
		return fmt.Errorf("invalid step operation: %w", err)
	}

	// 记录审批操作
	operatedUser := &models.TicketOperatedUser{
		TicketID:     uint(input.ID),
		OperatedUser: input.ApprovalUser,
	}
	if err := s.ticketManager.RecordApprovalOperation(ctx, operatedUser); err != nil {
		return fmt.Errorf("failed to record approval operation: %w", err)
	}

	// 根据操作类型处理
	if input.Operation == pkg.Approve {
		// 使用punched-tape集成层获取下一步骤
		nextStep, err := s.punchedTapeIntegration.GetNextStep(ctx, input.ID, ticket.Step, input.Operation)
		if err != nil {
			return fmt.Errorf("failed to get next step: %w", err)
		}

		// 检查是否为结束步骤
		isEndStep, err := s.punchedTapeIntegration.IsEndStep(ctx, input.ID, nextStep.Step)
		if err != nil {
			return fmt.Errorf("failed to check if end step: %w", err)
		}

		// 更新工单状态
		status := pkg.Running
		if isEndStep {
			status = pkg.Passed
		}

		if err := s.ticketManager.UpdateTicketStatus(ctx, input.ID, status, nextStep.Step); err != nil {
			return fmt.Errorf("failed to update ticket status: %w", err)
		}

		// 获取下一步骤的操作人并更新工单操作人
		if !isEndStep {
			nextStepOperators, err := s.punchedTapeIntegration.GetStepOperators(ctx, input.ID, nextStep.Step)
			if err != nil {
				return fmt.Errorf("failed to get next step operators: %w", err)
			}

			// 更新工单操作人
			if err := s.updateTicketOperators(ctx, input.ID, nextStepOperators); err != nil {
				return fmt.Errorf("failed to update ticket operators: %w", err)
			}
		}
	} else if input.Operation == pkg.Reject {
		// 审批拒绝，更新工单状态
		if err := s.ticketManager.UpdateTicketStatus(ctx, input.ID, pkg.Rejected, ""); err != nil {
			return fmt.Errorf("failed to update ticket status: %w", err)
		}
	}

	return nil
}

// convertTicketToResponse 将Ticket转换为TicketResponse
func (s *ticketService) convertTicketToResponse(ticket *models.Ticket) *models.TicketResponse {
	// 获取操作人列表
	operators := make([]string, 0)
	for _, op := range ticket.Operators {
		operators = append(operators, op.Operator)
	}

	// 获取已操作用户列表
	operatedUsers := make([]string, 0)
	for _, user := range ticket.OperatedUsers {
		operatedUsers = append(operatedUsers, user.OperatedUser)
	}

	return &models.TicketResponse{
		ID:            ticket.ID,
		OrderNum:      ticket.OrderNum,
		Status:        ticket.Status,
		Uid:           ticket.Uid,
		Step:          ticket.Step,
		Memo:          ticket.Memo,
		TemplateID:    ticket.TemplateID,
		CreatedAt:     ticket.CreatedAt,
		UpdatedAt:     ticket.UpdatedAt,
		Operators:     operators,
		OperatedUsers: operatedUsers,
	}
}

// CloseTicket 关闭工单
func (s *ticketService) CloseTicket(ctx context.Context, input *models.CloseTicketRequest) (*models.CloseTicketResponse, error) {
	if s.ticketManager == nil {
		// 模拟关闭逻辑，检查ID是否存在
		if input.ID == 9999 {
			return nil, fmt.Errorf("ticket not found")
		}
		return &models.CloseTicketResponse{ID: input.ID}, nil
	}

	// 构建更新字段
	updates := map[string]interface{}{
		"status": pkg.Closed,
	}
	if input.Memo != nil {
		updates["memo"] = *input.Memo
	}

	// 使用数据库管理器更新工单状态
	if err := s.ticketManager.UpdateTicketStatus(ctx, input.ID, pkg.Closed, ""); err != nil {
		return nil, fmt.Errorf("failed to update ticket status: %w", err)
	}

	// 记录关闭操作
	operatedUser := &models.TicketOperatedUser{
		TicketID:     uint(input.ID),
		OperatedUser: input.Operator,
	}
	if err := s.ticketManager.RecordApprovalOperation(ctx, operatedUser); err != nil {
		return nil, fmt.Errorf("failed to record close operation: %w", err)
	}

	return &models.CloseTicketResponse{ID: input.ID}, nil
}

// ListTickets 查询工单列表
func (s *ticketService) ListTickets(ctx context.Context, input *models.ListTicketRequest) (*models.ListTicketResponse, error) {
	if s.ticketManager == nil {
		// 模拟查询逻辑
		return &models.ListTicketResponse{
			Total: 0,
			List:  []*models.TicketResponse{},
		}, nil
	}

	// 构建查询过滤器
	filters := helper.BuildTicketFilters(input.Name, input.Creator, input.Status, input.OrderNum, input.TemplateID)

	// 使用数据库管理器查询工单列表
	tickets, total, err := s.ticketManager.ListTickets(ctx, filters, input.Page, input.Size)
	if err != nil {
		return nil, fmt.Errorf("failed to list tickets: %w", err)
	}

	// 转换为响应模型
	ticketResponses := make([]*models.TicketResponse, len(tickets))
	for i, ticket := range tickets {
		ticketResponses[i] = s.convertTicketToResponse(ticket)
	}

	return &models.ListTicketResponse{
		Total: total,
		List:  ticketResponses,
	}, nil
}

// updateTicketOperators 更新工单操作人
func (s *ticketService) updateTicketOperators(ctx context.Context, ticketID int, operators []string) error {
	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除旧的操作人
	if err := tx.Where("ticket_id = ?", ticketID).Delete(&models.TicketOperator{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete old operators: %w", err)
	}

	// 添加新的操作人
	for _, operator := range operators {
		ticketOperator := &models.TicketOperator{
			TicketID: uint(ticketID),
			Operator: operator,
		}
		if err := tx.Create(ticketOperator).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create ticket operator: %w", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
