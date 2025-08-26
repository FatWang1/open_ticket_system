# 工单模板验证器 (Ticket Template Validator)

## 概述

工单模板验证器使用 `github.com/go-playground/validator/v10` 库来验证工单模板相关API请求字段的必传性和类型。

## 功能特性

- **字段必传验证**: 使用 `required` 标签验证必填字段
- **类型验证**: 自动验证字段类型（string、int、uint等）
- **长度验证**: 使用 `min`、`max` 标签验证字符串长度
- **数值范围验证**: 使用 `gt`、`lt` 等标签验证数值范围
- **数组验证**: 使用 `min` 标签验证数组最小长度

## 验证标签说明

### 基础验证标签

- `required`: 字段必填
- `min`: 最小长度/最小值
- `max`: 最大长度/最大值
- `gt`: 大于指定值
- `omitempty`: 字段为空时跳过验证

## 使用方法

### 1. 导入验证器

```go
import "github.com/FatWang1/open_ticket_system/ticket_template/validator"
```

### 2. 验证请求

```go
// 验证创建工单模板请求
req := &models.CreateTicketTemplateRequest{
    Name:        "模板名称",
    Memo:        "模板描述",
    Version:     "v1.0",
    Creator:     "创建者",
    StartStep:   "开始步骤",
    EndStepList: []string{"结束步骤1", "结束步骤2"},
    ConfigList:  []*models.StepConfig{...},
}

err := validator.ValidateCreateTicketTemplateRequest(req)
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

### 工单模板验证

- `ValidateCreateTicketTemplateRequest(req *models.CreateTicketTemplateRequest) error`
- `ValidateUpdateTicketTemplateRequest(req *models.UpdateTicketTemplateRequest) error`
- `ValidateListTicketTemplateRequest(req *models.ListTicketTemplateRequest) error`

## 验证规则示例

### 创建工单模板请求

```go
type CreateTicketTemplateRequest struct {
    Name        string        `json:"name" validate:"required,min=1,max=100"`        // 必填，长度1-100
    Memo        string        `json:"memo" validate:"omitempty,max=1000"`            // 可选，最大长度1000
    Version     string        `json:"version" validate:"required,min=1,max=20"`     // 必填，长度1-20
    Creator     string        `json:"creator" validate:"required,min=1,max=50"`     // 必填，长度1-50
    StartStep   string        `json:"start_step" validate:"required,min=1,max=100"` // 必填，长度1-100
    EndStepList []string      `json:"end_step" validate:"required,min=1"`           // 必填，至少1个元素
    ConfigList  []*StepConfig `json:"config" validate:"required,min=1"`            // 必填，至少1个配置
}
```

### 更新工单模板请求

```go
type UpdateTicketTemplateRequest struct {
    ID          int           `json:"id" validate:"required,gt=0"`                    // 必填，大于0
    Name        *string       `json:"name,omitempty" validate:"omitempty,min=1,max=100"` // 可选，长度1-100
    Memo        *string       `json:"memo,omitempty" validate:"omitempty,max=1000"`      // 可选，最大长度1000
    Creator     *string       `json:"creator,omitempty" validate:"omitempty,min=1,max=50"` // 可选，长度1-50
    StartStep   *string       `json:"start_step,omitempty" validate:"omitempty,min=1,max=100"` // 可选，长度1-100
    EndStepList []string      `json:"end_step,omitempty" validate:"omitempty,min=1"`           // 可选，至少1个元素
    ConfigList  []*StepConfig `json:"config,omitempty" validate:"omitempty,min=1"`            // 可选，至少1个配置
}
```

### 查询工单模板列表请求

```go
type ListTicketTemplateRequest struct {
    Page    int     `json:"page" validate:"min=1,max=10000"`                    // 最小1，最大10000
    Size    int     `json:"size" validate:"min=1,max=1000"`                     // 最小1，最大1000
    Name    *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`  // 可选，长度1-100
    Creator *string `json:"creator,omitempty" validate:"omitempty,min=1,max=50"` // 可选，长度1-50
    Builtin *bool   `json:"builtin,omitempty"`                                   // 可选，无验证规则
}
```

## 测试

运行测试：

```bash
go test ./ticket_template/validator
```

## 注意事项

1. 验证器会自动设置分页参数的默认值（Page=1, Size=10）
2. 所有验证规则都在API模型的结构体标签中定义
3. 验证失败时会返回详细的错误信息，包含失败的字段和验证规则
4. 支持嵌套结构体的验证
5. 数组字段使用 `min=1` 确保至少有一个元素
6. 可选字段使用 `omitempty` 标签，当字段为空时跳过验证
