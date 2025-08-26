# AI Project Design Template

## 1. 项目基本信息
- 项目名称：open_ticket_system
- 模块路径：github.com/FatWang1/open_ticket_system
- 技术栈：Go, Gin, GORM, MySQL

### 1.1 业务逻辑核心
完全遵循  https://github.com/FatWang1/punched-tape 最新版本

## 2. 核心数据结构

### 2.1 基本数据结构
```go
// 基础实体结构 使用 GORM 作为 ORM

// User 是用户表结构体
type User struct {
	// gorm.Model 包含了常用的 ID, CreatedAt, UpdatedAt, DeletedAt 字段
	gorm.Model

	// 核心字段
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"` // `json:"-"` 防止密码被序列化到 JSON
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Nickname string `gorm:"type:varchar(50)" json:"nickname"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`

	// 账户状态字段
	Status      int8      `gorm:"type:tinyint;default:1;comment:1正常,2禁用,3注销" json:"status"`
	LastLoginAt time.Time `gorm:"type:datetime" json:"last_login_at"`

	// 角色与权限（如果需要）
	RoleID uint `gorm:"index" json:"role_id"`
	// 可选：一个简单的角色字段，而不是一个单独的表
	// Role string `gorm:"type:varchar(20);default:user" json:"role"`
}


// 工单主表
type Ticket struct {
    gorm.Model
    OrderNum      string `gorm:"column:order_num;type:varchar(255);not null;index" json:"order_num"`
    Status        string `gorm:"column:status;type:varchar(50);not null;index" json:"status"`
    Uid           string `gorm:"column:uid;type:varchar(255);uniqueIndex;not null" json:"uid"`
    Step          string `gorm:"column:step;type:varchar(255);not null" json:"step"`
    Memo          string `gorm:"column:memo;type:text" json:"memo"`
    TemplateID    uint   `gorm:"column:template_id;index" json:"template_id"` // 外键使用template_id
    Operators     []TicketOperator
    OperatedUsers []TicketOperatedUser
}

// 工单操作人表
type TicketOperator struct {
gorm.Model
TicketID uint   `gorm:"column:ticket_id;index;not null" json:"ticket_id"` // 外键使用ticket_id
Operator string `gorm:"column:operator;type:varchar(255);not null" json:"operator"`
Ticket   Ticket `gorm:"foreignKey:ID;references:TicketID"`
}

// 工单已操作用户表
type TicketOperatedUser struct {
gorm.Model
TicketID     uint   `gorm:"column:ticket_id;index;not null" json:"ticket_id"` // 外键使用ticket_id
OperatedUser string `gorm:"column:operated_user;type:varchar(255);not null" json:"operated_user"`
Ticket       Ticket `gorm:"foreignKey:ID;references:TicketID"`
}

// 工单模板表
type TicketTemplate struct {
gorm.Model
Uid         string `gorm:"column:uid;type:varchar(255);uniqueIndex;not null" json:"uid"`
StartStep   string `gorm:"column:start_step;type:varchar(255);not null" json:"start_step"`
Builtin     bool   `gorm:"column:builtin;default:false" json:"builtin"`
EndSteps    []TemplateEndStep
StepConfigs []StepConfig
}

// 模板结束节点表
type TemplateEndStep struct {
gorm.Model
TicketTemplateID uint   `gorm:"column:ticket_template_id;index;not null" json:"ticket_template_id"` // 外键使用ticket_template_id
EndStep          string `gorm:"column:end_step;type:varchar(255);not null" json:"end_step"`
TicketTemplate   TicketTemplate `gorm:"foreignKey:ID;references:TicketTemplateID"`
}

// 步骤配置表
type StepConfig struct {
gorm.Model
TicketTemplateID uint    `gorm:"column:ticket_template_id;index;not null" json:"ticket_template_id"` // 外键使用ticket_template_id
Step             string  `gorm:"column:step;type:varchar(255);not null;index" json:"step"`
State            string  `gorm:"column:state;type:varchar(255)" json:"state"`
SignType         string  `gorm:"column:sign_type;type:varchar(50);not null" json:"sign_type"`
JointSignRate    float32 `gorm:"column:joint_sign_rate;type:float" json:"joint_sign_rate"`
Operators        []StepOperator
NextSteps        []NextStep
}

// 步骤操作人表
type StepOperator struct {
gorm.Model
StepConfigID uint   `gorm:"column:step_config_id;index;not null" json:"step_config_id"` // 外键使用step_config_id
Operator     string `gorm:"column:operator;type:varchar(255);not null" json:"operator"`
StepConfig   StepConfig `gorm:"foreignKey:ID;references:StepConfigID"`
}

// 下一步配置表
type NextStep struct {
gorm.Model
StepConfigID uint   `gorm:"column:step_config_id;index;not null" json:"step_config_id"` // 外键使用step_config_id
ToStep       string `gorm:"column:to_step;type:varchar(255);not null" json:"to_step"`
Operation    string `gorm:"column:operation;type:varchar(255);not null" json:"operation"`
StepConfig   StepConfig `gorm:"foreignKey:ID;references:StepConfigID"`
}
```

### 2.3 api数据结构

```go
// 创建实体输入结构
type CreateTicketRequest struct {
    Name string `json:"name" validate:"required"` // 名称（必填）
	Creator string `json:"creator" validate:"required"`
	TemplateId  uint `json:"template_id" validate:"required"`
    Memo string `json:"memo"`
}

type CreateTicketTemplateResponse struct {
	Id   int `json:"id"`
}

// 更新实体输入结构
type UpdateTicketRequest struct {
    Id  int    `json:"id" validate:"required"`   // 实体ID（必填）
    Memo *string `json:"memo omitempty"`
}

type UpdateTicketTemplateResponse struct {
	Id   int `json:"id"`
}

type DeleteTicketRequest struct {
	Id   int `json:"id" validate:"required"`
}

type DeleteTicketResponse struct {
	Id   int `json:"id"`
}

type ApprovalRequest struct {
	Id   int `json:"id" validate:"required"`
	ApprovalUser string `json:"approval_user" validate:"required"`
	Memo *string `json:"memo omitempty"`
	Operation string `json:"operation" validate:"required"`
	NextStep string `json:"next_step" validate:"required"`
}

type CloseTicketRequest  struct {
	Id   int `json:"id" validate:"required"`
	Memo *string `json:"memo omitempty"`
	Operator string `json:"operator" validate:"required"`
}

type CloseTicketResponse struct {
	Id   int `json:"id"`
}

type ListTicketRequest struct {
	Page int `json:"page" default:"1"`
	Size int `json:"size" default:"10"`
	Name *string `json:"name omitempty"`
	Creator *string `json:"creator omitempty"`
	Status *string `json:"status omitempty"`
	TemplateId *int `json:"template_id omitempty"`
	OrderNum *string `json:"order_num omitempty"`
}

type CreateTicketTemplateRequest struct {
	Name string `json:"name" validate:"required"`
	Memo string `json:"memo"`
	Version string `json:"version" validate:"required"`
	Creator string `json:"creator" validate:"required"`
	StartStep string `json:"start_step" validate:"required"`
	EndStepList   []string `json:"end_step" validate:"required"`
	ConfigList []*StepConfig `json:"config"  validate:"required"`
}

type CreateTicketTemplateRequest struct {
	Id   int `json:"id" validate:"required"`
}

type UpdateTicketTemplateRequest struct {
	Id   int `json:"id" validate:"required"`
	Name *string `json:"name omitempty"`
	Memo *string `json:"memo omitempty"`
	Creator *string `json:"creator omitempty"`
	StartStep *string `json:"start_step omitempty"`
	EndStepList   []string `json:"end_step omitempty"`
	ConfigList []*StepConfig `json:"config omitempty"`
}

type UpdateTicketTemplateResponse struct {
	Id   int `json:"id"`
}

type ListTicketTemplateRequest struct {
	Page int `json:"page" default:"1"`
	Size int `json:"size" default:"10"`
	Name *string `json:"name omitempty"`
	Creator *string `json:"creator omitempty"`
	Builtin *bool `json:"builtin omitempty"`
}

type ListTicketTemplateResponse struct {
	Total int64 `json:"total"`
	List []*TicketTemplate `json:"list"`
}

type DeleteTicketTemplateRequest struct {
	Id   int `json:"id" validate:"required"`
}

type DeleteTicketTemplateResponse struct {
	Id   int `json:"id"`
}


```

## 3. 主要接口/方法

```go
// 实体服务接口
type TicketService interface {
	CreateTicket(ctx context.Context, input *CreateTicketRequest) (*CreateTicketResponse, error)
	GetTicketByID(ctx context.Context, id int) (*Ticket, error)
	UpdateTicket(ctx context.Context, input *UpdateTicketRequest) (*UpdateTicketResponse, error)
	DeleteTicket(ctx context.Context, input *DeleteTicketRequest) (*DeleteTicketResponse, error)
	Approval(ctx context.Context, input *ApprovalRequest) error
	CloseTicket(ctx context.Context, input *CloseTicketRequest) (*CloseTicketResponse, error)
}

type TicketTemplateService interface {
	CreateTicketTemplate(ctx context.Context, input *CreateTicketTemplateRequest) (*CreateTicketTemplateResponse, error)
	GetTicketTemplateByID(ctx context.Context, id int) (*TicketTemplate, error)
	UpdateTicketTemplate(ctx context.Context, input *UpdateTicketTemplateRequest) (*UpdateTicketTemplateResponse, error)
	DeleteTicketTemplate(ctx context.Context, input *DeleteTicketTemplateRequest) (*DeleteTicketTemplateResponse, error)
	ListTicketTemplates(ctx context.Context, input *ListTicketTemplateRequest) (*ListTicketTemplateResponse, error)
}

// REST API路由
func (s *TicketController) Register(engine *gin.Engine) {
	engine.POST("/tickets", s.CreateTicket)
	engine.GET("/tickets/:id", s.GetTicketByID)
	engine.PUT("/tickets/:id", s.UpdateTicket)
	engine.DELETE("/tickets/:id", s.DeleteTicket)
	engine.POST("/tickets/:id/approval", s.Approval)
	engine.POST("/tickets/:id/close", s.CloseTicket)
}

func (s *TicketTemplateController) Register(engine *gin.Engine) {
	engine.POST("/ticket_templates", s.CreateTicketTemplate)
	engine.GET("/ticket_templates/:id", s.GetTicketTemplateByID)
	engine.PUT("/ticket_templates/:id", s.UpdateTicketTemplate)
	engine.DELETE("/ticket_templates/:id", s.DeleteTicketTemplate)
	engine.GET("/ticket_templates", s.ListTicketTemplates)
}
```

## 4. 配置与扩展

```go
// 应用配置结构
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    Redis    RedisConfig    `yaml:"redis"`
    Log      LogConfig      `yaml:"log"`
}

// 服务器配置
type ServerConfig struct {
    Port    int    `yaml:"port" default:"8080"`
    Host    string `yaml:"host" default:"localhost"`
    Timeout int    `yaml:"timeout" default:"30"`
}

// 数据库配置
type DatabaseConfig struct {
    Driver   string `yaml:"driver" default:"postgres"`
    Host     string `yaml:"host" default:"localhost"`
    Port     int    `yaml:"port" default:"5432"`
    Username string `yaml:"username"`
    Password string `yaml:"password"`
    Database string `yaml:"database"`
    SSLMode  string `yaml:"ssl_mode" default:"disable"`
}
```

## 5. 错误处理

使用 `github.com/pkg/errors` 进行错误处理

## 6. 中间件设计
使用 https://github.com/FatWang1/fatwang-go-utils/blob/master/utils/logger.go 进行日志记录
使用 https://github.com/golang-jwt/jwt/v5 进行JWT验证
使用 https://github.com/gin-contrib/cors 进行跨域请求
使用 https://github.com/FatWang1/punched-tape 进行核心业务处理 包括但不限于 进行工单/工单模版的创建、更新、删除、查询、审批、关闭
使用 https://github.com/go-playground/validator 进行参数验证
使用 https://github.com/spf13/viper 进行配置管理

## 7. 项目配置信息
- GitHub仓库：FatWang1/open_ticket_system
- 默认分支：master
- Go版本：1.24
- 许可证：AGPL-3.0