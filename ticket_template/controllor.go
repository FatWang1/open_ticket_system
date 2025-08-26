package ticket_template

import (
	"net/http"
	"strconv"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"github.com/FatWang1/open_ticket_system/ticket_template/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TicketTemplateController 工单模板控制器
type TicketTemplateController struct {
	templateService TicketTemplateService
}

// NewTicketTemplateController 创建工单模板控制器实例
func NewTicketTemplateController(db *gorm.DB) *TicketTemplateController {
	return &TicketTemplateController{
		templateService: NewTicketTemplateService(db),
	}
}

// CreateTicketTemplate 创建工单模板
// @Summary      创建工单模板
// @Description  创建新的工单模板
// @Tags         ticket-templates
// @Accept       json
// @Produce      json
// @Param        request body models.CreateTicketTemplateRequest true "创建模板请求"
// @Success      200  {object}  models.CreateTicketTemplateResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /ticket-templates [post]
func (c *TicketTemplateController) CreateTicketTemplate(ctx *gin.Context) {
	var req models.CreateTicketTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用验证器校验请求
	if err := validator.ValidateCreateTicketTemplateRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.templateService.CreateTicketTemplate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetTicketTemplateByID 根据ID获取工单模板
// @Summary      获取工单模板详情
// @Description  根据ID获取工单模板详细信息
// @Tags         ticket-templates
// @Accept       json
// @Produce      json
// @Param        id path int true "模板ID"
// @Success      200  {object}  models.TicketTemplateResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      404  {object}  map[string]interface{} "模板不存在"
// @Router       /ticket-templates/{id} [get]
func (c *TicketTemplateController) GetTicketTemplateByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	template, err := c.templateService.GetTicketTemplateByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, template)
}

// UpdateTicketTemplate 更新工单模板
// @Summary      更新工单模板
// @Description  更新工单模板信息
// @Tags         ticket-templates
// @Accept       json
// @Produce      json
// @Param        id path int true "模板ID"
// @Param        request body models.UpdateTicketTemplateRequest true "更新模板请求"
// @Success      200  {object}  models.UpdateTicketTemplateResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      404  {object}  map[string]interface{} "模板不存在"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /ticket-templates/{id} [put]
func (c *TicketTemplateController) UpdateTicketTemplate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req models.UpdateTicketTemplateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id

	// 使用验证器校验请求
	if err := validator.ValidateUpdateTicketTemplateRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.templateService.UpdateTicketTemplate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteTicketTemplate 删除工单模板
// @Summary      删除工单模板
// @Description  根据ID删除工单模板
// @Tags         ticket-templates
// @Accept       json
// @Produce      json
// @Param        id path int true "模板ID"
// @Success      200  {object}  models.DeleteTicketTemplateResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /ticket-templates/{id} [delete]
func (c *TicketTemplateController) DeleteTicketTemplate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	req := models.DeleteTicketTemplateRequest{ID: id}
	response, err := c.templateService.DeleteTicketTemplate(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ListTicketTemplates 查询工单模板列表
// @Summary      查询工单模板列表
// @Description  分页查询工单模板列表，支持筛选和排序
// @Tags         ticket-templates
// @Accept       json
// @Produce      json
// @Param        page query int false "页码" default(1)
// @Param        size query int false "每页数量" default(10)
// @Param        name query string false "模板名称"
// @Param        creator query string false "创建者"
// @Param        version query string false "版本号"
// @Param        builtin query bool false "是否内置"
// @Success      200  {object}  models.ListTicketTemplateResponse
// @Failure      400  {object}  map[string]interface{} "请求参数错误"
// @Failure      500  {object}  map[string]interface{} "服务器内部错误"
// @Router       /ticket-templates [get]
func (c *TicketTemplateController) ListTicketTemplates(ctx *gin.Context) {
	var req models.ListTicketTemplateRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用验证器校验请求
	if err := validator.ValidateListTicketTemplateRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.templateService.ListTicketTemplates(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}
