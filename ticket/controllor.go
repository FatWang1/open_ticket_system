package ticket

import (
	"net/http"
	"strconv"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/FatWang1/open_ticket_system/ticket/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TicketController 工单控制器
type TicketController struct {
	ticketService TicketService
}

// NewTicketController 创建工单控制器实例
func NewTicketController(db *gorm.DB) *TicketController {
	return &TicketController{
		ticketService: NewTicketService(db),
	}
}

// CreateTicket 创建工单
// @Summary      创建工单
// @Description  根据模板创建新的工单
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        request body models.CreateTicketRequest true "创建工单请求"
// @Success      200  {object}  models.CreateTicketResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /tickets [post]
func (c *TicketController) CreateTicket(ctx *gin.Context) {
	var req models.CreateTicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用验证器校验请求
	if err := validator.ValidateCreateTicketRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.ticketService.CreateTicket(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTicketByID 根据ID获取工单
// @Summary      获取工单详情
// @Description  根据工单ID获取工单详细信息
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        id path int true "工单ID"
// @Success      200  {object}  models.TicketResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      404  {object}  map[string]interface{} "工单不存在"
// @Router       /tickets/{id} [get]
func (c *TicketController) GetTicketByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	ticket, err := c.ticketService.GetTicketByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, ticket)
}

// UpdateTicket 更新工单
// @Summary      更新工单
// @Description  更新工单信息
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        id path int true "工单ID"
// @Param        request body models.UpdateTicketRequest true "更新工单请求"
// @Success      200  {object}  models.UpdateTicketResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      404  {object}  map[string]interface{} "工单不存在"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /tickets/{id} [put]
func (c *TicketController) UpdateTicket(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateTicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id

	// 使用验证器校验请求
	if err := validator.ValidateUpdateTicketRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.ticketService.UpdateTicket(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteTicket 删除工单
// @Summary      删除工单
// @Description  根据ID删除工单
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        id path int true "工单ID"
// @Success      200  {object}  models.DeleteTicketResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /tickets/{id} [delete]
func (c *TicketController) DeleteTicket(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	req := models.DeleteTicketRequest{ID: id}
	response, err := c.ticketService.DeleteTicket(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// Approval 工单审批
// @Summary      工单审批
// @Description  审批工单，支持通过、拒绝等操作
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        id path int true "工单ID"
// @Param        request body models.ApprovalRequest true "审批请求"
// @Success      200  {object}  map[string]interface{} "审批成功"
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /tickets/{id}/approval [post]
func (c *TicketController) Approval(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.ApprovalRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id

	// 使用验证器校验请求
	if err := validator.ValidateApprovalRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.ticketService.Approval(ctx, &req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "approval successful"})
}

// CloseTicket 关闭工单
// @Summary      关闭工单
// @Description  关闭工单
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        id path int true "工单ID"
// @Param        request body models.CloseTicketRequest true "关闭工单请求"
// @Success      200  {object}  models.CloseTicketResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /tickets/{id}/close [post]
func (c *TicketController) CloseTicket(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.CloseTicketRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id

	// 使用验证器校验请求
	if err := validator.ValidateCloseTicketRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.ticketService.CloseTicket(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ListTickets 查询工单列表
// @Summary      查询工单列表
// @Description  分页查询工单列表，支持筛选和排序
// @Tags         tickets
// @Accept       json
// @Produce      json
// @Param        page query int false "页码" default(1)
// @Param        size query int false "每页数量" default(10)
// @Param        name query string false "工单名称"
// @Param        creator query string false "创建者"
// @Param        status query string false "工单状态"
// @Param        template_id query int false "模板ID"
// @Success      200  {object}  models.ListTicketResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /tickets [get]
func (c *TicketController) ListTickets(ctx *gin.Context) {
	var req models.ListTicketRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用验证器校验请求
	if err := validator.ValidateListTicketRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.ticketService.ListTickets(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
