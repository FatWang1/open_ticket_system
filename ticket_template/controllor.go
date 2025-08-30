package ticket_template

import (
	"net/http"
	"strconv"

	"github.com/FatWang1/open_ticket_system/internal/manager"
	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/FatWang1/open_ticket_system/internal/utils"
	"github.com/FatWang1/open_ticket_system/ticket_template/validator"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// TicketTemplateController 工单模板控制器
type TicketTemplateController struct {
	service *TicketTemplateService
}

// NewTicketTemplateController 创建新的工单模板控制器
func NewTicketTemplateController(db interface{}) *TicketTemplateController {
	// 创建管理器
	templateManager := manager.NewTicketTemplateManager(db.(*gorm.DB))
	punchedTapeIntegration := manager.NewPunchedTapeIntegration(db.(*gorm.DB))

	// 创建服务
	service := NewTicketTemplateService(templateManager, punchedTapeIntegration)

	return &TicketTemplateController{
		service: service,
	}
}

// Register 注册路由
func (c *TicketTemplateController) Register(router *gin.RouterGroup) {
	templates := router.Group("/ticket_templates")
	{
		templates.POST("", c.CreateTicketTemplate)
		templates.GET("", c.ListTicketTemplates)
		templates.GET("/:id", c.GetTicketTemplateByID)
		templates.PUT("/:id", c.UpdateTicketTemplate)
		templates.DELETE("/:id", c.DeleteTicketTemplate)
	}
}

// CreateTicketTemplate 创建工单模板
func (c *TicketTemplateController) CreateTicketTemplate(ctx *gin.Context) {
	logger := utils.GetLogger()

	logger.Printf("[INFO] Creating ticket template - request from: %s", ctx.ClientIP())

	var request models.CreateTicketTemplateAPI
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Printf("[ERROR] Failed to bind JSON request - error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 验证请求
	if err := validator.ValidateCreateTicketTemplateRequest(&request); err != nil {
		logger.Printf("[ERROR] Request validation failed - error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用服务
	template, err := c.service.CreateTicketTemplate(ctx, &request)
	if err != nil {
		logger.Printf("[ERROR] Failed to create ticket template - error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Printf("[INFO] Successfully created ticket template - id: %d", template.ID)
	ctx.JSON(http.StatusCreated, template)
}

// GetTicketTemplateByID 根据ID获取工单模板
func (c *TicketTemplateController) GetTicketTemplateByID(ctx *gin.Context) {
	logger := utils.GetLogger()

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Printf("[ERROR] Invalid ID parameter - id: %s, error: %v", idStr, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	logger.Printf("[INFO] Getting ticket template by ID - id: %d, request from: %s", id, ctx.ClientIP())

	template, err := c.service.GetTicketTemplateByID(ctx, id)
	if err != nil {
		logger.Printf("[ERROR] Failed to get ticket template - id: %d, error: %v", id, err)
		if errors.Is(err, errors.New("ticket template not found")) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Ticket template not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Printf("[INFO] Successfully retrieved ticket template - id: %d", id)
	ctx.JSON(http.StatusOK, template)
}

// ListTicketTemplates 获取工单模板列表
func (c *TicketTemplateController) ListTicketTemplates(ctx *gin.Context) {
	logger := utils.GetLogger()

	// 获取分页参数
	pageStr := ctx.DefaultQuery("page", "1")
	sizeStr := ctx.DefaultQuery("size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size < 1 || size > 100 {
		size = 10
	}

	logger.Printf("[INFO] Listing ticket templates - page: %d, size: %d, request from: %s", page, size, ctx.ClientIP())

	templates, total, err := c.service.ListTicketTemplates(ctx, page, size)
	if err != nil {
		logger.Printf("[ERROR] Failed to list ticket templates - error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := gin.H{
		"templates": templates,
		"total":     total,
		"page":      page,
		"size":      size,
	}

	logger.Printf("[INFO] Successfully listed ticket templates - count: %d, total: %d", len(templates), total)
	ctx.JSON(http.StatusOK, result)
}

// UpdateTicketTemplate 更新工单模板
func (c *TicketTemplateController) UpdateTicketTemplate(ctx *gin.Context) {
	logger := utils.GetLogger()

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Printf("[ERROR] Invalid ID parameter - id: %s, error: %v", idStr, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	var request models.UpdateTicketTemplateAPI
	if err := ctx.ShouldBindJSON(&request); err != nil {
		logger.Printf("[ERROR] Failed to bind JSON request - id: %d, error: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// 验证请求
	if err := validator.ValidateUpdateTicketTemplateRequest(&request); err != nil {
		logger.Printf("[ERROR] Request validation failed - id: %d, error: %v", id, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	logger.Printf("[INFO] Updating ticket template - id: %d, request from: %s", id, ctx.ClientIP())

	// 调用服务
	template, err := c.service.UpdateTicketTemplate(ctx, id, &request)
	if err != nil {
		logger.Printf("[ERROR] Failed to update ticket template - id: %d, error: %v", id, err)
		if errors.Is(err, errors.New("ticket template not found")) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Ticket template not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Printf("[INFO] Successfully updated ticket template - id: %d", id)
	ctx.JSON(http.StatusOK, template)
}

// DeleteTicketTemplate 删除工单模板
func (c *TicketTemplateController) DeleteTicketTemplate(ctx *gin.Context) {
	logger := utils.GetLogger()

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Printf("[ERROR] Invalid ID parameter - id: %s, error: %v", idStr, err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	logger.Printf("[INFO] Deleting ticket template - id: %d, request from: %s", id, ctx.ClientIP())

	// 调用服务
	if err := c.service.DeleteTicketTemplate(ctx, id); err != nil {
		logger.Printf("[ERROR] Failed to delete ticket template - id: %d, error: %v", id, err)
		if errors.Is(err, errors.New("ticket template not found")) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Ticket template not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Printf("[INFO] Successfully deleted ticket template - id: %d", id)
	ctx.JSON(http.StatusOK, gin.H{"message": "Ticket template deleted successfully"})
}
