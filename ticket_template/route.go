package ticket_template

import (
	"github.com/gin-gonic/gin"
)

// Register 注册工单模板路由
func (c *TicketTemplateController) Register(router *gin.RouterGroup) {
	// 工单模板管理路由组
	templateGroup := router.Group("/ticket_templates")
	{
		// 创建工单模板
		templateGroup.POST("", c.CreateTicketTemplate)

		// 查询工单模板列表
		templateGroup.GET("", c.ListTicketTemplates)

		// 根据ID获取工单模板
		templateGroup.GET("/:id", c.GetTicketTemplateByID)

		// 更新工单模板
		templateGroup.PUT("/:id", c.UpdateTicketTemplate)

		// 删除工单模板
		templateGroup.DELETE("/:id", c.DeleteTicketTemplate)
	}
}
