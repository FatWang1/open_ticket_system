package models

import (
	"time"

	"gorm.io/gorm"
)

// Ticket 工单主表
// @Description 工单主表结构
type Ticket struct {
	gorm.Model
	OrderNum      string               `gorm:"column:order_num;type:varchar(255);not null;index" json:"order_num" example:"TICKET-20250127143000"` // 工单号
	Status        string               `gorm:"column:status;type:varchar(50);not null;index" json:"status" example:"running"`                      // 工单状态
	Uid           string               `gorm:"column:uid;type:varchar(255);uniqueIndex;not null" json:"uid" example:"uid-1234567890"`              // 工单唯一标识
	Step          string               `gorm:"column:step;type:varchar(255);not null" json:"step" example:"submit"`                                // 当前步骤
	Memo          string               `gorm:"column:memo;type:text" json:"memo" example:"需要请假3天"`                                                 // 工单备注
	TemplateID    uint                 `gorm:"column:template_id;index" json:"template_id" example:"1"`                                            // 模板ID
	Template      TicketTemplate       `gorm:"-"`                                                                                                  // 逻辑关联，不创建物理外键
	Operators     []TicketOperator     `gorm:"-"`
	OperatedUsers []TicketOperatedUser `gorm:"-"`
}

// TicketOperator 工单操作人表
type TicketOperator struct {
	gorm.Model
	TicketID uint   `gorm:"column:ticket_id;index;not null" json:"ticket_id"`
	Operator string `gorm:"column:operator;type:varchar(255);not null" json:"operator"`
	Ticket   Ticket `gorm:"-"` // 逻辑关联，不创建物理外键
}

// TicketOperatedUser 工单已操作用户表
type TicketOperatedUser struct {
	gorm.Model
	TicketID     uint   `gorm:"column:ticket_id;index;not null" json:"ticket_id"`
	OperatedUser string `gorm:"column:operated_user;type:varchar(255);not null" json:"operated_user"`
	Ticket       Ticket `gorm:"-"` // 逻辑关联，不创建物理外键
}

// TicketTemplate 工单模板表
// @Description 工单模板表结构
type TicketTemplate struct {
	gorm.Model
	Name      string `gorm:"column:name;type:varchar(100);not null" json:"name" example:"请假申请模板"`                 // 模板名称
	Memo      string `gorm:"column:memo;type:text" json:"memo" example:"请假申请流程"`                                  // 模板备注
	Version   string `gorm:"column:version;type:varchar(50);not null" json:"version" example:"1.0"`               // 版本号
	Creator   string `gorm:"column:creator;type:varchar(50);not null" json:"creator" example:"管理员"`               // 创建者
	Uid       string `gorm:"column:uid;type:varchar(255);uniqueIndex;not null" json:"uid" example:"template-123"` // 模板唯一标识
	StartStep string `gorm:"column:start_step;type:varchar(255);not null" json:"start_step" example:"submit"`     // 起始步骤
	Builtin   bool   `gorm:"column:builtin;default:false" json:"builtin" example:"false"`                         // 是否内置
	// EndSteps    []TemplateEndStep `gorm:"-"`
	// StepConfigs []StepConfigDB    `gorm:"-"`
}

// TemplateEndStep 模板结束节点表
type TemplateEndStep struct {
	gorm.Model
	TicketTemplateID uint           `gorm:"column:ticket_template_id;index;not null" json:"ticket_template_id"`
	EndStep          string         `gorm:"column:end_step;type:varchar(255);not null" json:"end_step"`
	TicketTemplate   TicketTemplate `gorm:"-"` // 逻辑关联，不创建物理外键
}

// StepConfigDB 步骤配置表（数据库模型）
type StepConfigDB struct {
	gorm.Model
	TicketTemplateID uint           `gorm:"column:ticket_template_id;index;not null" json:"ticket_template_id"`
	Step             string         `gorm:"column:step;type:varchar(255);not null;index" json:"step"`
	State            string         `gorm:"column:state;type:varchar(255)" json:"state"`
	SignType         string         `gorm:"column:sign_type;type:varchar(50);not null" json:"sign_type"`
	JointSignRate    float32        `gorm:"column:joint_sign_rate;type:float" json:"joint_sign_rate"`
	TicketTemplate   TicketTemplate `gorm:"-"` // 逻辑关联，不创建物理外键
	Operators        []StepOperator `gorm:"-"`
	NextSteps        []NextStep     `gorm:"-"`
}

// StepOperator 步骤操作人表
type StepOperator struct {
	gorm.Model
	StepConfigID uint         `gorm:"column:step_config_id;index;not null" json:"step_config_id"`
	Operator     string       `gorm:"column:operator;type:varchar(255);not null" json:"operator"`
	StepConfig   StepConfigDB `gorm:"-"` // 逻辑关联，不创建物理外键
}

// NextStep 下一步配置表
type NextStep struct {
	gorm.Model
	StepConfigID uint         `gorm:"column:step_config_id;index;not null" json:"step_config_id"`
	ToStep       string       `gorm:"column:to_step;type:varchar(255);not null" json:"to_step"`
	Operation    string       `gorm:"column:operation;type:varchar(255);not null" json:"operation"`
	StepConfig   StepConfigDB `gorm:"-"` // 逻辑关联，不创建物理外键
}

// User 用户表结构体
type User struct {
	ID          uint       `gorm:"primarykey" json:"id" example:"1"`
	CreatedAt   time.Time  `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2024-01-01T00:00:00Z"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Username    string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username" example:"admin"`
	Password    string     `gorm:"type:varchar(255);not null" json:"-"`
	Email       string     `gorm:"type:varchar(100);uniqueIndex" json:"email" example:"admin@example.com"`
	Nickname    string     `gorm:"type:varchar(50)" json:"nickname" example:"管理员"`
	Avatar      string     `gorm:"type:varchar(255)" json:"avatar" example:"https://example.com/avatar.jpg"`
	Status      int8       `gorm:"type:tinyint;default:1;comment:1正常,2禁用,3注销" json:"status" example:"1"`
	LastLoginAt *time.Time `gorm:"type:datetime" json:"last_login_at" example:"2024-01-01T00:00:00Z"`
	RoleID      uint       `gorm:"index" json:"role_id" example:"1"`
}
