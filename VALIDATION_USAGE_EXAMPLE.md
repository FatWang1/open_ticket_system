# 验证器在控制器中的使用示例

## 概述

本文档展示了如何在控制器中使用验证器进行请求字段校验，确保API的健壮性和数据完整性。

## 修改内容

### 1. 工单控制器 (`ticket/controllor.go`)

在以下方法中添加了验证器校验：

- `CreateTicket()` - 创建工单
- `UpdateTicket()` - 更新工单  
- `Approval()` - 工单审批
- `CloseTicket()` - 关闭工单
- `ListTickets()` - 查询工单列表

### 2. 工单模板控制器 (`ticket_template/controllor.go`)

在以下方法中添加了验证器校验：

- `CreateTicketTemplate()` - 创建工单模板
- `UpdateTicketTemplate()` - 更新工单模板
- `ListTicketTemplates()` - 查询工单模板列表

## 使用示例

### 创建工单示例

```go
// CreateTicket 创建工单
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
```

### 工单审批示例

```go
// Approval 工单审批
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
```

## 验证流程

### 1. 请求解析
```go
var req models.CreateTicketRequest
if err := ctx.ShouldBindJSON(&req); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

### 2. 字段验证
```go
if err := validator.ValidateCreateTicketRequest(&req); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

### 3. 业务处理
```go
response, err := c.ticketService.CreateTicket(ctx, &req)
if err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
    return
}
```

### 4. 返回结果
```go
ctx.JSON(http.StatusOK, response)
```

## 验证失败示例

### 请求数据
```json
{
    "name": "",           // 空字符串，违反 required,min=1 规则
    "creator": "user",    // 有效
    "template_id": 0      // 0值，违反 gt=0 规则
}
```

### 验证失败响应
```json
{
    "error": "Key: 'CreateTicketRequest.Name' Error:Field validation for 'Name' failed on the 'required' tag\nKey: 'CreateTicketRequest.TemplateID' Error:Field validation for 'TemplateID' failed on the 'gt' tag"
}
```

### HTTP状态码
```
400 Bad Request
```

## 验证规则回顾

### 创建工单请求验证规则
```go
type CreateTicketRequest struct {
    Name       string `json:"name" validate:"required,min=1,max=100"`        // 必填，长度1-100
    Creator    string `json:"creator" validate:"required,min=1,max=50"`      // 必填，长度1-50
    TemplateID uint   `json:"template_id" validate:"required,gt=0"`          // 必填，大于0
    Memo       string `json:"memo" validate:"omitempty,max=1000"`            // 可选，最大长度1000
}
```

### 工单审批请求验证规则
```go
type ApprovalRequest struct {
    ID           int     `json:"id" validate:"required,gt=0"`                    // 必填，大于0
    ApprovalUser string  `json:"approval_user" validate:"required,min=1,max=50"` // 必填，长度1-50
    Memo         *string `json:"memo,omitempty" validate:"omitempty,max=1000"`   // 可选，最大长度1000
    Operation    string  `json:"operation" validate:"required,oneof=approve reject"` // 必填，必须是approve或reject
    NextStep     string  `json:"next_step" validate:"required,min=1,max=100"`    // 必填，长度1-100
}
```

## 优势

### 1. 数据完整性
- 确保必填字段不为空
- 验证字段长度和数值范围
- 检查枚举值有效性

### 2. 错误处理
- 验证失败时立即返回错误
- 提供详细的字段验证错误信息
- 统一的错误响应格式

### 3. 代码健壮性
- 在业务逻辑执行前进行验证
- 减少无效数据进入业务层
- 提高系统整体稳定性

### 4. 开发体验
- 验证规则集中在模型定义中
- 控制器代码简洁清晰
- 易于维护和扩展

## 注意事项

1. **验证顺序**: 先进行请求解析，再进行字段验证，最后执行业务逻辑
2. **错误处理**: 验证失败时直接返回400状态码和错误信息
3. **性能考虑**: 验证器性能优秀，对API响应时间影响很小
4. **扩展性**: 如需添加自定义验证规则，可以在验证器中注册自定义函数

## 总结

通过在控制器中集成验证器，我们实现了：

- ✅ 完整的请求字段验证
- ✅ 统一的错误处理机制
- ✅ 清晰的代码结构
- ✅ 优秀的开发体验

这为API的健壮性和可维护性提供了有力保障。
