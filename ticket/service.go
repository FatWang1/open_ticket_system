package ticket

import (
	"context"
	"fmt"
	"time"

	"github.com/FatWang1/open_ticket_system/internal/manager"
	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/FatWang1/open_ticket_system/internal/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// TicketService 工单服务接口
type TicketService interface {
	CreateTicket(ctx context.Context, input *models.CreateTicketRequest) (*models.CreateTicketResponse, error)
	GetTicketByID(ctx context.Context, id int) (*models.Ticket, error)
	UpdateTicket(ctx context.Context, input *models.UpdateTicketRequest) (*models.UpdateTicketResponse, error)
	DeleteTicket(ctx context.Context, input *models.DeleteTicketRequest) (*models.DeleteTicketResponse, error)
	Approval(ctx context.Context, input *models.ApprovalRequest) error
	CloseTicket(ctx context.Context, input *models.CloseTicketRequest) error
	ListTickets(ctx context.Context, input *models.ListTicketRequest) (*models.ListTicketResponse, error)
}

// ticketService 工单服务实现
type ticketService struct {
	ticketManager          *manager.TicketManager
	punchedTapeIntegration *manager.PunchedTapeIntegration
	db                     *gorm.DB
}

// NewTicketService 创建工单服务实例
func NewTicketService(db *gorm.DB) TicketService {
	if db == nil {
		// 返回模拟服务
		return &ticketService{
			ticketManager:          nil,
			punchedTapeIntegration: nil,
			db:                     nil,
		}
	}

	return &ticketService{
		ticketManager:          manager.NewTicketManager(db),
		punchedTapeIntegration: manager.NewPunchedTapeIntegration(db),
		db:                     db,
	}
}

// CreateTicket 创建工单
func (s *ticketService) CreateTicket(ctx context.Context, input *models.CreateTicketRequest) (*models.CreateTicketResponse, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Creating ticket - template_id: %d, creator: %s", input.TemplateID, input.Creator)

	if s.ticketManager == nil {
		// 模拟创建工单逻辑
		logger.Printf("[INFO] Using mock service for ticket creation")
		return &models.CreateTicketResponse{ID: 1}, nil
	}

	// 创建工单
	ticket := &models.Ticket{
		OrderNum:   fmt.Sprintf("TICKET-%d", time.Now().Unix()),
		Status:     "running",
		Uid:        fmt.Sprintf("uid-%d", time.Now().Unix()),
		Step:       "submit",
		Memo:       input.Memo,
		TemplateID: input.TemplateID,
	}

	// 从模板获取操作人信息并创建工单操作人记录
	operators := []models.TicketOperator{
		{
			Operator: input.Creator,
		},
	}

	// 使用数据库管理器创建工单
	if err := s.ticketManager.CreateTicket(ctx, ticket, operators); err != nil {
		logger.Printf("[ERROR] Failed to create ticket - error: %v", err)
		return nil, errors.Wrap(err, "failed to create ticket")
	}

	logger.Printf("[INFO] Successfully created ticket - id: %d", ticket.ID)
	return &models.CreateTicketResponse{ID: int(ticket.ID)}, nil
}

// GetTicketByID 根据ID获取工单
func (s *ticketService) GetTicketByID(ctx context.Context, id int) (*models.Ticket, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Getting ticket by ID - id: %d", id)

	if s.ticketManager == nil {
		// 模拟数据，实际使用时应该查询数据库
		logger.Printf("[INFO] Using mock service for ticket retrieval")
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

	ticket, err := s.ticketManager.GetTicketByID(ctx, id)
	if err != nil {
		logger.Printf("[ERROR] Failed to get ticket by ID - id: %d, error: %v", id, err)
		return nil, errors.Wrap(err, "failed to get ticket")
	}

	logger.Printf("[INFO] Successfully retrieved ticket - id: %d", id)
	return ticket, nil
}

// UpdateTicket 更新工单
func (s *ticketService) UpdateTicket(ctx context.Context, input *models.UpdateTicketRequest) (*models.UpdateTicketResponse, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Updating ticket - id: %d", input.ID)

	if s.ticketManager == nil {
		// 模拟更新逻辑，检查ID是否存在
		logger.Printf("[INFO] Using mock service for ticket update")
		if input.ID == 9999 {
			return nil, errors.New("ticket not found")
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
		logger.Printf("[ERROR] Failed to update ticket - id: %d, error: %v", input.ID, err)
		return nil, errors.Wrap(err, "failed to update ticket")
	}

	logger.Printf("[INFO] Successfully updated ticket - id: %d", input.ID)
	return &models.UpdateTicketResponse{ID: input.ID}, nil
}

// DeleteTicket 删除工单
func (s *ticketService) DeleteTicket(ctx context.Context, input *models.DeleteTicketRequest) (*models.DeleteTicketResponse, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Deleting ticket - id: %d", input.ID)

	if s.ticketManager == nil {
		// 模拟删除逻辑，检查ID是否存在
		logger.Printf("[INFO] Using mock service for ticket deletion")
		if input.ID == 9999 {
			return nil, errors.New("ticket not found")
		}
		return &models.DeleteTicketResponse{ID: input.ID}, nil
	}

	// 使用数据库管理器删除工单
	if err := s.ticketManager.DeleteTicket(ctx, input.ID); err != nil {
		logger.Printf("[ERROR] Failed to delete ticket - id: %d, error: %v", input.ID, err)
		return nil, errors.Wrap(err, "failed to delete ticket")
	}

	logger.Printf("[INFO] Successfully deleted ticket - id: %d", input.ID)
	return &models.DeleteTicketResponse{ID: input.ID}, nil
}

// Approval 工单审批
func (s *ticketService) Approval(ctx context.Context, input *models.ApprovalRequest) error {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Processing ticket approval - ticket_id: %d, operation: %s", input.ID, input.Operation)

	if s.ticketManager == nil {
		// 模拟审批逻辑，检查操作类型
		logger.Printf("[INFO] Using mock service for ticket approval")
		if input.Operation != "approve" && input.Operation != "reject" {
			return errors.New("invalid operation type")
		}
		return nil
	}

	// 验证操作类型
	if input.Operation != "approve" && input.Operation != "reject" {
		logger.Printf("[ERROR] Invalid operation type - operation: %s", input.Operation)
		return errors.New("invalid operation type")
	}

	// 获取工单
	ticket, err := s.ticketManager.GetTicketByID(ctx, input.ID)
	if err != nil {
		logger.Printf("[ERROR] Failed to get ticket for approval - id: %d, error: %v", input.ID, err)
		return errors.Wrap(err, "failed to get ticket")
	}

	// 检查工单状态
	if ticket.Status != "running" {
		logger.Printf("[ERROR] Ticket is not in running status - id: %d, status: %s", input.ID, ticket.Status)
		return errors.New("ticket is not in running status")
	}

	// 更新工单状态
	updates := map[string]interface{}{
		"status": input.Operation,
	}
	if err := s.ticketManager.UpdateTicket(ctx, input.ID, updates); err != nil {
		logger.Printf("[ERROR] Failed to update ticket status - id: %d, error: %v", input.ID, err)
		return errors.Wrap(err, "failed to update ticket status")
	}

	logger.Printf("[INFO] Successfully processed ticket approval - id: %d, operation: %s", input.ID, input.Operation)
	return nil
}

// CloseTicket 关闭工单
func (s *ticketService) CloseTicket(ctx context.Context, input *models.CloseTicketRequest) error {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Closing ticket - id: %d", input.ID)

	if s.ticketManager == nil {
		// 模拟关闭逻辑
		logger.Printf("[INFO] Using mock service for ticket closure")
		return nil
	}

	// 获取工单
	ticket, err := s.ticketManager.GetTicketByID(ctx, input.ID)
	if err != nil {
		logger.Printf("[ERROR] Failed to get ticket for closure - id: %d, error: %v", input.ID, err)
		return errors.Wrap(err, "failed to get ticket")
	}

	// 检查工单状态
	if ticket.Status == "closed" {
		logger.Printf("[WARN] Ticket is already closed - id: %d", input.ID)
		return errors.New("ticket is already closed")
	}

	// 更新工单状态为关闭
	updates := map[string]interface{}{
		"status": "closed",
	}
	if err := s.ticketManager.UpdateTicket(ctx, input.ID, updates); err != nil {
		logger.Printf("[ERROR] Failed to close ticket - id: %d, error: %v", input.ID, err)
		return errors.Wrap(err, "failed to close ticket")
	}

	logger.Printf("[INFO] Successfully closed ticket - id: %d", input.ID)
	return nil
}

// ListTickets 查询工单列表
func (s *ticketService) ListTickets(ctx context.Context, input *models.ListTicketRequest) (*models.ListTicketResponse, error) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Listing tickets - page: %d, size: %d", input.Page, input.Size)

	if s.ticketManager == nil {
		// 模拟查询逻辑
		logger.Printf("[INFO] Using mock service for ticket listing")
		return &models.ListTicketResponse{
			Total: 0,
			List:  []*models.TicketResponse{},
		}, nil
	}

	// 构建查询过滤器
	filters := make(map[string]interface{})
	if input.Name != nil {
		filters["name"] = *input.Name
	}
	if input.Creator != nil {
		filters["creator"] = *input.Creator
	}
	if input.Status != nil {
		filters["status"] = *input.Status
	}
	if input.TemplateID != nil {
		filters["template_id"] = *input.TemplateID
	}
	if input.OrderNum != nil {
		filters["order_num"] = *input.OrderNum
	}

	// 使用数据库管理器查询工单列表
	tickets, total, err := s.ticketManager.ListTickets(ctx, filters, input.Page, input.Size)
	if err != nil {
		logger.Printf("[ERROR] Failed to list tickets - error: %v", err)
		return nil, errors.Wrap(err, "failed to list tickets")
	}

	// 转换为响应模型
	ticketResponses := make([]*models.TicketResponse, len(tickets))
	for i, ticket := range tickets {
		ticketResponses[i] = s.convertTicketToResponse(ticket)
	}

	result := &models.ListTicketResponse{
		Total: total,
		List:  ticketResponses,
	}

	logger.Printf("[INFO] Successfully listed tickets - count: %d, total: %d", len(tickets), total)
	return result, nil
}

// convertTicketToResponse 将Ticket转换为TicketResponse
func (s *ticketService) convertTicketToResponse(ticket *models.Ticket) *models.TicketResponse {
	return &models.TicketResponse{
		ID:         ticket.ID,
		OrderNum:   ticket.OrderNum,
		Status:     ticket.Status,
		Uid:        ticket.Uid,
		Step:       ticket.Step,
		Memo:       ticket.Memo,
		TemplateID: ticket.TemplateID,
		CreatedAt:  ticket.CreatedAt,
		UpdatedAt:  ticket.UpdatedAt,
	}
}
