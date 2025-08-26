package ticket

import (
	"github.com/gin-gonic/gin"
)

// Register 注册工单路由
func (c *TicketController) Register(router *gin.RouterGroup) {
	// 工单管理路由组
	ticketGroup := router.Group("/tickets")
	{
		// 创建工单
		ticketGroup.POST("", c.CreateTicket)

		// 查询工单列表
		ticketGroup.GET("", c.ListTickets)

		// 根据ID获取工单
		ticketGroup.GET("/:id", c.GetTicketByID)

		// 更新工单
		ticketGroup.PUT("/:id", c.UpdateTicket)

		// 删除工单
		ticketGroup.DELETE("/:id", c.DeleteTicket)

		// 工单审批
		ticketGroup.POST("/:id/approval", c.Approval)

		// 关闭工单
		ticketGroup.POST("/:id/close", c.CloseTicket)
	}
}
