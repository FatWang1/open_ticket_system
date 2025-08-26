package helper

import (
	"fmt"
	"time"
)

// GenerateOrderNum 生成工单号
func GenerateOrderNum() string {
	return fmt.Sprintf("TICKET-%s", time.Now().Format("20060102150405"))
}

// GenerateUID 生成唯一标识
func GenerateUID() string {
	return fmt.Sprintf("uid-%d", time.Now().UnixNano())
}

// IsEndStep 检查是否为结束步骤
func IsEndStep(step string) bool {
	// 默认的结束步骤列表，实际使用时应该从模板配置中获取
	endSteps := []string{"end", "complete", "finished", "approved", "rejected", "closed"}
	for _, endStep := range endSteps {
		if step == endStep {
			return true
		}
	}
	return false
}

// ValidateTicketStatus 验证工单状态
func ValidateTicketStatus(currentStatus, targetStatus string) error {
	// 定义状态转换规则
	validTransitions := map[string][]string{
		"running":  {"passed", "rejected", "closed"},
		"passed":   {},
		"rejected": {},
		"closed":   {},
	}

	if allowedStatuses, exists := validTransitions[currentStatus]; exists {
		for _, allowed := range allowedStatuses {
			if allowed == targetStatus {
				return nil
			}
		}
		return fmt.Errorf("invalid status transition from %s to %s", currentStatus, targetStatus)
	}

	return fmt.Errorf("unknown status: %s", currentStatus)
}

// BuildTicketFilters 构建工单查询过滤器
func BuildTicketFilters(name, creator, status, orderNum *string, templateID *int) map[string]interface{} {
	filters := make(map[string]interface{})

	if name != nil && *name != "" {
		filters["name"] = *name
	}
	if creator != nil && *creator != "" {
		filters["creator"] = *creator
	}
	if status != nil && *status != "" {
		filters["status"] = *status
	}
	if templateID != nil && *templateID > 0 {
		filters["template_id"] = *templateID
	}
	if orderNum != nil && *orderNum != "" {
		filters["order_num"] = *orderNum
	}

	return filters
}

// BuildTemplateFilters 构建模板查询过滤器
func BuildTemplateFilters(name, creator, version *string, builtin *bool) map[string]interface{} {
	filters := make(map[string]interface{})

	if name != nil && *name != "" {
		filters["name"] = *name
	}
	if creator != nil && *creator != "" {
		filters["creator"] = *creator
	}
	if version != nil && *version != "" {
		filters["version"] = *version
	}
	if builtin != nil {
		filters["builtin"] = *builtin
	}

	return filters
}

// FormatTimestamp 格式化时间戳
func FormatTimestamp(timestamp time.Time) string {
	return timestamp.Format("2006-01-02 15:04:05")
}

// ParseTimestamp 解析时间戳字符串
func ParseTimestamp(timestampStr string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", timestampStr)
}

// ValidateStepOperation 验证步骤操作
func ValidateStepOperation(operation string) error {
	validOperations := []string{"approve", "reject", "close", "reopen", "transfer"}
	for _, valid := range validOperations {
		if operation == valid {
			return nil
		}
	}
	return fmt.Errorf("invalid operation: %s", operation)
}

// ValidateTemplateConfig 验证模板配置
func ValidateTemplateConfig(startStep string, endSteps []string, stepConfigs []string) error {
	if startStep == "" {
		return fmt.Errorf("start step cannot be empty")
	}

	if len(endSteps) == 0 {
		return fmt.Errorf("end steps cannot be empty")
	}

	if len(stepConfigs) == 0 {
		return fmt.Errorf("step configs cannot be empty")
	}

	// 检查起始步骤是否在步骤配置中
	startStepFound := false
	for _, step := range stepConfigs {
		if step == startStep {
			startStepFound = true
			break
		}
	}

	if !startStepFound {
		return fmt.Errorf("start step '%s' not found in step configs", startStep)
	}

	return nil
}

// GenerateStepName 生成步骤名称
func GenerateStepName(prefix string, index int) string {
	return fmt.Sprintf("%s_%d", prefix, index)
}

// ValidateUserPermission 验证用户权限
func ValidateUserPermission(userID string, requiredRole string, currentRole string) error {
	if currentRole == "admin" {
		return nil // 管理员拥有所有权限
	}

	if currentRole == requiredRole {
		return nil
	}

	return fmt.Errorf("user %s with role %s does not have required role %s", userID, currentRole, requiredRole)
}

// CalculateApprovalProgress 计算审批进度
func CalculateApprovalProgress(currentStep string, totalSteps int, completedSteps int) float64 {
	if totalSteps == 0 {
		return 0.0
	}

	progress := float64(completedSteps) / float64(totalSteps) * 100.0
	if progress > 100.0 {
		progress = 100.0
	}

	return progress
}

// IsStepReversible 检查步骤是否可逆
func IsStepReversible(step string) bool {
	reversibleSteps := []string{"draft", "submitted", "reviewing"}
	for _, reversible := range reversibleSteps {
		if step == reversible {
			return true
		}
	}
	return false
}

// GetStepDeadline 获取步骤截止时间
func GetStepDeadline(createdAt time.Time, stepType string) time.Time {
	var deadline time.Duration

	switch stepType {
	case "urgent":
		deadline = 24 * time.Hour
	case "normal":
		deadline = 72 * time.Hour
	case "low":
		deadline = 168 * time.Hour // 7天
	default:
		deadline = 72 * time.Hour
	}

	return createdAt.Add(deadline)
}

// FormatDuration 格式化持续时间
func FormatDuration(duration time.Duration) string {
	if duration < time.Minute {
		return fmt.Sprintf("%.0fs", duration.Seconds())
	}
	if duration < time.Hour {
		return fmt.Sprintf("%.0fm", duration.Minutes())
	}
	if duration < 24*time.Hour {
		return fmt.Sprintf("%.0fh", duration.Hours())
	}
	return fmt.Sprintf("%.0fd", duration.Hours()/24)
}
