package pkg

// 审批类型
const (
	JointlySign = "jointly_sign" // 联合审批
	SerialSign  = "serial_sign"  // 串行审批
	AnyoneSign  = "anyone_sign"  // 任意人审批
)

// 工单状态
const (
	Running  = "running"  // 运行中
	Passed   = "passed"   // 已通过
	Rejected = "rejected" // 已拒绝
	Closed   = "closed"   // 已关闭
)

// 操作类型
const (
	Reject  = "reject"  // 拒绝
	Approve = "approve" // 通过
)

// 用户状态
const (
	UserStatusNormal   = 1 // 正常
	UserStatusDisabled = 2 // 禁用
	UserStatusDeleted  = 3 // 注销
)

// 分页默认值
const (
	DefaultPage = 1
	DefaultSize = 10
	MaxPage     = 1000
	MaxSize     = 100
)

// 字段长度限制
const (
	MaxNameLength     = 100
	MaxCreatorLength  = 50
	MaxMemoLength     = 1000
	MaxVersionLength  = 50
	MaxStepLength     = 100
	MaxOrderNumLength = 255
)
