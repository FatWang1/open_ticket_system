package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"gorm.io/gorm"
)

// WorkflowMonitor 工单流程监控
type WorkflowMonitor struct {
	db *gorm.DB
}

// NewWorkflowMonitor 创建流程监控实例
func NewWorkflowMonitor(db *gorm.DB) *WorkflowMonitor {
	return &WorkflowMonitor{db: db}
}

// WorkflowMetrics 流程指标
type WorkflowMetrics struct {
	TotalTickets     int64         `json:"total_tickets"`
	RunningTickets   int64         `json:"running_tickets"`
	CompletedTickets int64         `json:"completed_tickets"`
	RejectedTickets  int64         `json:"rejected_tickets"`
	ClosedTickets    int64         `json:"closed_tickets"`
	AvgProcessTime   time.Duration `json:"avg_process_time"`
	OverdueTickets   int64         `json:"overdue_tickets"`
}

// GetWorkflowMetrics 获取流程指标
func (wm *WorkflowMonitor) GetWorkflowMetrics(ctx context.Context) (*WorkflowMetrics, error) {
	metrics := &WorkflowMetrics{}

	// 获取各状态工单数量
	if err := wm.db.Model(&models.Ticket{}).Count(&metrics.TotalTickets).Error; err != nil {
		return nil, fmt.Errorf("failed to count total tickets: %w", err)
	}

	if err := wm.db.Model(&models.Ticket{}).Where("status = ?", "running").Count(&metrics.RunningTickets).Error; err != nil {
		return nil, fmt.Errorf("failed to count running tickets: %w", err)
	}

	if err := wm.db.Model(&models.Ticket{}).Where("status = ?", "passed").Count(&metrics.CompletedTickets).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed tickets: %w", err)
	}

	if err := wm.db.Model(&models.Ticket{}).Where("status = ?", "rejected").Count(&metrics.RejectedTickets).Error; err != nil {
		return nil, fmt.Errorf("failed to count rejected tickets: %w", err)
	}

	if err := wm.db.Model(&models.Ticket{}).Where("status = ?", "closed").Count(&metrics.ClosedTickets).Error; err != nil {
		return nil, fmt.Errorf("failed to count closed tickets: %w", err)
	}

	// 计算平均处理时间
	var avgProcessTime float64
	if err := wm.db.Model(&models.Ticket{}).
		Select("AVG(TIMESTAMPDIFF(HOUR, created_at, updated_at))").
		Where("status IN (?)", []string{"passed", "rejected", "closed"}).
		Scan(&avgProcessTime).Error; err != nil {
		return nil, fmt.Errorf("failed to calculate avg process time: %w", err)
	}

	metrics.AvgProcessTime = time.Duration(avgProcessTime) * time.Hour

	// 计算逾期工单数量（超过72小时的运行中工单）
	overdueTime := time.Now().Add(-72 * time.Hour)
	if err := wm.db.Model(&models.Ticket{}).
		Where("status = ? AND created_at < ?", "running", overdueTime).
		Count(&metrics.OverdueTickets).Error; err != nil {
		return nil, fmt.Errorf("failed to count overdue tickets: %w", err)
	}

	return metrics, nil
}

// GetStepPerformance 获取步骤性能统计
func (wm *WorkflowMonitor) GetStepPerformance(ctx context.Context, templateID int) (map[string]interface{}, error) {
	var stepStats []struct {
		Step           string  `json:"step"`
		AvgProcessTime float64 `json:"avg_process_time"`
		TotalCount     int64   `json:"total_count"`
		SuccessCount   int64   `json:"success_count"`
		FailureCount   int64   `json:"failure_count"`
	}

	query := `
		SELECT 
			step,
			AVG(TIMESTAMPDIFF(HOUR, created_at, updated_at)) as avg_process_time,
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 'passed' THEN 1 ELSE 0 END) as success_count,
			SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as failure_count
		FROM tickets 
		WHERE template_id = ? 
		GROUP BY step
	`

	if err := wm.db.Raw(query, templateID).Scan(&stepStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get step performance: %w", err)
	}

	result := make(map[string]interface{})
	for _, stat := range stepStats {
		result[stat.Step] = map[string]interface{}{
			"avg_process_time": stat.AvgProcessTime,
			"total_count":      stat.TotalCount,
			"success_count":    stat.SuccessCount,
			"failure_count":    stat.FailureCount,
			"success_rate":     float64(stat.SuccessCount) / float64(stat.TotalCount) * 100,
		}
	}

	return result, nil
}

// GetUserPerformance 获取用户性能统计
func (wm *WorkflowMonitor) GetUserPerformance(ctx context.Context, timeRange time.Duration) (map[string]interface{}, error) {
	var userStats []struct {
		Username        string  `json:"username"`
		TotalTickets    int64   `json:"total_tickets"`
		ApprovedCount   int64   `json:"approved_count"`
		RejectedCount   int64   `json:"rejected_count"`
		AvgResponseTime float64 `json:"avg_response_time"`
	}

	startTime := time.Now().Add(-timeRange)

	query := `
		SELECT 
			u.username,
			COUNT(DISTINCT t.id) as total_tickets,
			SUM(CASE WHEN t.status = 'passed' THEN 1 ELSE 0 END) as approved_count,
			SUM(CASE WHEN t.status = 'rejected' THEN 1 ELSE 0 END) as rejected_count,
			AVG(TIMESTAMPDIFF(HOUR, t.created_at, t.updated_at)) as avg_response_time
		FROM users u
		LEFT JOIN ticket_operators to ON u.username = to.operator
		LEFT JOIN tickets t ON to.ticket_id = t.id AND t.created_at >= ?
		GROUP BY u.username
		HAVING total_tickets > 0
	`

	if err := wm.db.Raw(query, startTime).Scan(&userStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get user performance: %w", err)
	}

	result := make(map[string]interface{})
	for _, stat := range userStats {
		result[stat.Username] = map[string]interface{}{
			"total_tickets":     stat.TotalTickets,
			"approved_count":    stat.ApprovedCount,
			"rejected_count":    stat.RejectedCount,
			"avg_response_time": stat.AvgResponseTime,
			"approval_rate":     float64(stat.ApprovedCount) / float64(stat.TotalTickets) * 100,
		}
	}

	return result, nil
}

// GetWorkflowTrends 获取流程趋势
func (wm *WorkflowMonitor) GetWorkflowTrends(ctx context.Context, days int) (map[string]interface{}, error) {
	var dailyStats []struct {
		Date           string `json:"date"`
		CreatedCount   int64  `json:"created_count"`
		CompletedCount int64  `json:"completed_count"`
		RejectedCount  int64  `json:"rejected_count"`
	}

	query := `
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as created_count,
			SUM(CASE WHEN status IN ('passed', 'rejected', 'closed') THEN 1 ELSE 0 END) as completed_count,
			SUM(CASE WHEN status = 'rejected' THEN 1 ELSE 0 END) as rejected_count
		FROM tickets 
		WHERE created_at >= DATE_SUB(CURDATE(), INTERVAL ? DAY)
		GROUP BY DATE(created_at)
		ORDER BY date
	`

	if err := wm.db.Raw(query, days).Scan(&dailyStats).Error; err != nil {
		return nil, fmt.Errorf("failed to get workflow trends: %w", err)
	}

	result := make(map[string]interface{})
	for _, stat := range dailyStats {
		result[stat.Date] = map[string]interface{}{
			"created_count":   stat.CreatedCount,
			"completed_count": stat.CompletedCount,
			"rejected_count":  stat.RejectedCount,
			"completion_rate": float64(stat.CompletedCount) / float64(stat.CreatedCount) * 100,
		}
	}

	return result, nil
}

// GetBottleneckSteps 获取瓶颈步骤
func (wm *WorkflowMonitor) GetBottleneckSteps(ctx context.Context) ([]map[string]interface{}, error) {
	var bottlenecks []struct {
		Step           string  `json:"step"`
		AvgProcessTime float64 `json:"avg_process_time"`
		TotalCount     int64   `json:"total_count"`
		PendingCount   int64   `json:"pending_count"`
	}

	query := `
		SELECT 
			step,
			AVG(TIMESTAMPDIFF(HOUR, created_at, updated_at)) as avg_process_time,
			COUNT(*) as total_count,
			SUM(CASE WHEN status = 'running' THEN 1 ELSE 0 END) as pending_count
		FROM tickets 
		WHERE status IN ('running', 'passed', 'rejected')
		GROUP BY step
		HAVING avg_process_time > 24 AND pending_count > 0
		ORDER BY avg_process_time DESC
	`

	if err := wm.db.Raw(query).Scan(&bottlenecks).Error; err != nil {
		return nil, fmt.Errorf("failed to get bottleneck steps: %w", err)
	}

	var result []map[string]interface{}
	for _, bottleneck := range bottlenecks {
		result = append(result, map[string]interface{}{
			"step":             bottleneck.Step,
			"avg_process_time": bottleneck.AvgProcessTime,
			"total_count":      bottleneck.TotalCount,
			"pending_count":    bottleneck.PendingCount,
			"bottleneck_score": bottleneck.AvgProcessTime * float64(bottleneck.PendingCount),
		})
	}

	return result, nil
}
