# 工单验证器 (Ticket Validator)

## 概述

工单验证器使用 `github.com/go-playground/validator/v10` 库来验证API请求字段的必传性和类型。

## 功能特性

- **字段必传验证**: 使用 `required` 标签验证必填字段
- **类型验证**: 自动验证字段类型（string、int、uint等）
- **长度验证**: 使用 `min`、`max` 标签验证字符串长度
- **数值范围验证**: 使用 `gt`、`lt` 等标签验证数值范围
- **枚举值验证**: 使用 `oneof` 标签验证字段值是否在指定范围内

## 验证标签说明

### 基础验证标签

- `required`: 字段必填
- `min`: 最小长度/最小值
- `max`: 最大长度/最大值
- `gt`: 大于指定值
- `oneof`: 值必须在指定范围内

### 工单相关验证标签

- `oneof=running passed rejected closed`: 工单状态必须是其中之一
- `oneof=approve reject`: 操作类型必须是其中之一

## 使用方法

### 1. 导入验证器

```go
import "github.com/FatWang1/open_ticket_system/ticket/validator"
```

### 2. 验证请求

```go
// 验证创建工单请求
req := &models.CreateTicketRequest{
    Name:       "工单名称",
    Creator:    "创建者",
    TemplateID: 1,
    Memo:       "备注信息",
}

err := validator.ValidateCreateTicketRequest(req)
if err != nil {
    // 处理验证错误
    log.Printf("验证失败: %v", err)
    return
}
```

### 3. 验证错误处理

验证失败时会返回 `validator.ValidationErrors` 类型的错误，包含详细的字段验证信息：

```go
if err != nil {
    if validationErrors, ok := err.(validator.ValidationErrors); ok {
        for _, fieldErr := range validationErrors {
            log.Printf("字段 %s 验证失败: %s", fieldErr.Field(), fieldErr.Tag())
        }
    }
    return
}
```

## 支持的验证函数

### 工单验证

- `ValidateCreateTicketRequest(req *models.CreateTicketRequest) error`
- `ValidateUpdateTicketRequest(req *models.UpdateTicketRequest) error`
- `ValidateApprovalRequest(req *models.ApprovalRequest) error`
- `ValidateCloseTicketRequest(req *models.CloseTicketRequest) error`
- `ValidateListTicketRequest(req *models.ListTicketRequest) error`

## 验证规则示例

### 创建工单请求

```go
type CreateTicketRequest struct {
    Name       string `json:"name" validate:"required,min=1,max=100"`        // 必填，长度1-100
    Creator    string `json:"creator" validate:"required,min=1,max=50"`      // 必填，长度1-50
    TemplateID uint   `json:"template_id" validate:"required,gt=0"`          // 必填，大于0
    Memo       string `json:"memo" validate:"omitempty,max=1000"`            // 可选，最大长度1000
}
```

### 工单审批请求

```go
type ApprovalRequest struct {
    ID           int     `json:"id" validate:"required,gt=0"`                    // 必填，大于0
    ApprovalUser string  `json:"approval_user" validate:"required,min=1,max=50"` // 必填，长度1-50
    Memo         *string `json:"memo,omitempty" validate:"omitempty,max=1000"`   // 可选，最大长度1000
    Operation    string  `json:"operation" validate:"required,oneof=approve reject"` // 必填，必须是approve或reject
    NextStep     string  `json:"next_step" validate:"required,min=1,max=100"`    // 必填，长度1-100
}
```

## 测试

运行测试：

```bash
go test ./ticket/validator
```

## 注意事项

1. 验证器会自动设置分页参数的默认值（Page=1, Size=10）
2. 所有验证规则都在API模型的结构体标签中定义
3. 验证失败时会返回详细的错误信息，包含失败的字段和验证规则
4. 支持嵌套结构体的验证
