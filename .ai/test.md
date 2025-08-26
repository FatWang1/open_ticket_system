# AI Project Test Template

## 1. 项目基本信息
- 项目名称：open_ticket_system
- 模块路径：github.com/FatWang1/open_ticket_system
- 技术栈：Go, Gin, GORM, MySQL

## 2. 测试策略

### 2.1 测试目标
- 功能完整性：确保所有核心功能可用
- 响应准确性：确保输出与预期一致
- 易用性与集成性：确保易于集成和扩展
- 健壮性：确保异常输入下系统稳定

### 2.2 测试类型
- 单元测试：聚焦核心算法/逻辑
- 集成测试：端到端流程验证
- 配置驱动测试：用配置文件驱动多场景测试

## 3. 单元测试函数模板

```golang
// Test_<FunctionName> 是一个针对 <FunctionName> 的单元测试模板。
func Test_<FunctionName>(t *testing.T) {
	// 1. 定义函数参数和返回值的类型
	type args struct {
		// 参数1: 类型
		// 参数2: 类型
		// ...
	}

	// 2. 定义测试用例切片
	tests := []struct {
		name    string // 测试用例的简短描述
		args    args   // 待测试函数的输入参数
		want    // 预期返回的结果
		wantErr bool   // 预期是否返回错误
	}{
		// 3. 填充具体的测试用例
		{
			name: "should_return_expected_result_for_valid_input",
			args: args{
				// 填充参数
			},
			want:    // 填充预期的成功结果
			wantErr: false,
		},
		{
			name: "should_return_error_for_invalid_input",
			args: args{
				// 填充参数以触发错误
			},
			want:    nil,
			wantErr: true,
		},
	}

	// 4. 遍历并执行测试
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := <FunctionName>(/* 传入 tt.args 中的参数 */)

			// 5. 错误断言
			if (err != nil) != tt.wantErr {
				t.Errorf("%s() error = %v, wantErr %v", <FunctionName>, err, tt.wantErr)
				return
			}

			// 6. 结果断言
			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("%s() got unexpected diff (-want +got):\n%s", <FunctionName>, diff)
			}
		})
	}
}
```

## 4. 测试用例模板

# 工单系统生产环境测试用例

## 1. 工单管理测试用例

| 用例ID | 测试类别 | 测试目标 | 前置条件 | 输入数据 | 预期行为/输出 | 验证点 | 优先级 |
|--------|----------|----------|----------|----------|---------------|--------|--------|
| TICKET-001 | 单元测试 | 创建工单成功 | 1. 数据库连接正常<br>2. 模板存在<br>3. 用户已认证 | {<br>  "name": "测试工单",<br>  "creator": "user1",<br>  "template_id": 1,<br>  "memo": "测试备注"<br>} | HTTP 200<br>{ "id": 1 } | 1. 工单记录成功创建<br>2. Uid生成符合规范<br>3. 状态初始化为"running"<br>4. 关联操作人正确设置 | 高 |
| TICKET-002 | 单元测试 | 创建工单-必填字段缺失 | 1. 数据库连接正常<br>2. 用户已认证 | {<br>  "creator": "user1",<br>  "template_id": 1<br>} | HTTP 400<br>{ "error": "name is required" } | 1. 正确返回验证错误<br>2. 无工单记录创建 | 高 |
| TICKET-003 | 单元测试 | 创建工单-模板不存在 | 1. 数据库连接正常<br>2. 用户已认证 | {<br>  "name": "测试工单",<br>  "creator": "user1",<br>  "template_id": 9999<br>} | HTTP 400<br>{ "error": "template not found" } | 1. 返回模板不存在错误<br>2. 无工单记录创建 | 高 |
| TICKET-004 | 单元测试 | 查询工单详情 | 1. 工单已存在<br>2. 用户有权限访问 | GET /tickets/1 | HTTP 200<br>返回完整的工单详情 | 1. 返回正确的工单信息<br>2. 包含所有关联数据（操作人、已操作用户等）<br>3. 敏感字段（如密码）不返回 | 高 |
| TICKET-005 | 单元测试 | 查询不存在的工单 | 1. 数据库连接正常<br>2. 用户已认证 | GET /tickets/9999 | HTTP 404<br>{ "error": "ticket not found" } | 1. 返回404错误<br>2. 错误信息明确 | 中 |
| TICKET-006 | 单元测试 | 更新工单备注 | 1. 工单已存在<br>2. 用户有权限更新 | PUT /tickets/1<br>{ "memo": "更新后的备注" } | HTTP 200<br>{ "id": 1 } | 1. 备注字段正确更新<br>2. 其他字段保持不变<br>3. UpdatedAt时间戳更新 | 高 |
| TICKET-007 | 单元测试 | 更新不存在的工单 | 1. 数据库连接正常<br>2. 用户已认证 | PUT /tickets/9999<br>{ "memo": "测试备注" } | HTTP 404<br>{ "error": "ticket not found" } | 1. 返回404错误<br>2. 无数据更新 | 中 |
| TICKET-008 | 单元测试 | 删除工单 | 1. 工单已存在<br>2. 用户有权限删除 | DELETE /tickets/1 | HTTP 200<br>{ "id": 1 } | 1. 工单标记为软删除<br>2. 关联数据正确处理<br>3. DeletedAt时间戳设置 | 高 |
| TICKET-009 | 单元测试 | 删除已删除的工单 | 1. 工单已软删除<br>2. 用户已认证 | DELETE /tickets/1 | HTTP 404<br>{ "error": "ticket not found" } | 1. 返回404错误<br>2. 无额外数据库操作 | 低 |
| TICKET-010 | 单元测试 | 工单审批-联合审批通过 | 1. 工单处于联合审批步骤<br>2. 审批人是操作人之一<br>3. 未达到通过率 | {<br>  "approval_user": "user1",<br>  "operation": "approve",<br>  "next_step": "step2"<br>} | HTTP 200 | 1. 已操作用户列表更新<br>2. 工单状态不变（因未达到通过率）<br>3. 无步骤跳转 | 高 |
| TICKET-011 | 单元测试 | 工单审批-联合审批通过率达标 | 1. 工单处于联合审批步骤<br>2. 审批人是最后一位达到通过率的操作人 | {<br>  "approval_user": "user3",<br>  "operation": "approve",<br>  "next_step": "step2"<br>} | HTTP 200 | 1. 已操作用户列表更新<br>2. 工单步骤跳转到next_step<br>3. 操作人列表重置为下一节点操作人 | 高 |
| TICKET-012 | 单元测试 | 工单审批-串行审批 | 1. 工单处于串行审批步骤<br>2. 审批人是当前步骤操作人 | {<br>  "approval_user": "user1",<br>  "operation": "approve",<br>  "next_step": "step2"<br>} | HTTP 200 | 1. 已操作用户列表更新<br>2. 工单步骤跳转到next_step<br>3. 操作人列表更新为下一节点操作人 | 高 |
| TICKET-013 | 单元测试 | 工单审批-任意一人审批 | 1. 工单处于任意一人审批步骤<br>2. 审批人是操作人之一 | {<br>  "approval_user": "user1",<br>  "operation": "approve",<br>  "next_step": "step2"<br>} | HTTP 200 | 1. 工单步骤立即跳转<br>2. 操作人列表更新为下一节点操作人 | 高 |
| TICKET-014 | 单元测试 | 工单审批-非操作人审批 | 1. 工单存在<br>2. 审批人不是当前步骤操作人 | {<br>  "approval_user": "user99",<br>  "operation": "approve",<br>  "next_step": "step2"<br>} | HTTP 403<br>{ "error": "not authorized to approve" } | 1. 返回权限错误<br>2. 工单状态不变 | 高 |
| TICKET-015 | 单元测试 | 关闭工单成功 | 1. 工单处于可关闭状态<br>2. 用户有权限关闭 | {<br>  "memo": "关闭备注",<br>  "operator": "admin"<br>} | HTTP 200<br>{ "id": 1 } | 1. 工单状态更新为"closed"<br>2. 备注更新<br>3. 已操作用户列表更新 | 高 |
| TICKET-016 | 单元测试 | 关闭已完成工单 | 1. 工单状态为"passed"<br>2. 用户已认证 | {<br>  "memo": "测试",<br>  "operator": "admin"<br>} | HTTP 400<br>{ "error": "ticket already completed" } | 1. 返回状态错误<br>2. 工单状态不变 | 中 |
| TICKET-017 | 集成测试 | 工单流程完整执行 | 1. 模板配置完整<br>2. 用户有权限 | 1. 创建工单<br>2. 多次审批<br>3. 最终关闭 | 所有步骤HTTP 200 | 1. 工单按模板流程执行<br>2. 每个步骤状态正确<br>3. 最终状态为"passed" | 高 |
| TICKET-018 | 集成测试 | 并发审批处理 | 1. 工单处于联合审批步骤<br>2. 多个审批人同时审批 | 多个审批请求同时发送 | 所有请求HTTP 200 | 1. 避免并发问题<br>2. 通过率计算准确<br>3. 仅一次状态变更 | 高 |
| TICKET-019 | 边界测试 | 工单号长度边界 | 1. 数据库连接正常<br>2. 用户已认证 | {<br>  "name": "A",<br>  "creator": "user1",<br>  "template_id": 1,<br>  "memo": "测试"<br>} | HTTP 200 | 1. 工单成功创建<br>2. 工单号符合长度要求 | 中 |

## 2. 工单模板管理测试用例

| 用例ID | 测试类别 | 测试目标 | 前置条件 | 输入数据 | 预期行为/输出 | 验证点 | 优先级 |
|--------|----------|----------|----------|----------|---------------|--------|--------|
| TEMPLATE-001 | 单元测试 | 创建模板成功 | 1. 数据库连接正常<br>2. 用户有权限 | {<br>  "name": "测试模板",<br>  "version": "1.0",<br>  "creator": "admin",<br>  "start_step": "step1",<br>  "end_step": ["step3"],<br>  "config": [{<br>    "step": "step1",<br>    "sign_type": "anyone_sign",<br>    "operator": ["user1", "user2"],<br>    "next": [{"step": "step2", "operation": "approve"}]<br>  }]<br>} | HTTP 200<br>{ "id": 1 } | 1. 模板记录成功创建<br>2. Uid生成符合规范<br>3. 所有关联数据正确存储 | 高 |
| TEMPLATE-002 | 单元测试 | 创建模板-必填字段缺失 | 1. 数据库连接正常<br>2. 用户已认证 | {<br>  "name": "测试模板",<br>  "creator": "admin"<br>} | HTTP 400<br>{ "error": "version is required" } | 1. 返回验证错误<br>2. 无模板记录创建 | 高 |
| TEMPLATE-003 | 单元测试 | 创建模板-联合审批配置 | 1. 数据库连接正常<br>2. 用户有权限 | {<br>  "sign_type": "jointly_sign",<br>  "joint_sign_rate": 0.6<br>} | HTTP 200 | 1. joint_sign_rate正确存储<br>2. 仅在sign_type为jointly_sign时有效 | 高 |
| TEMPLATE-004 | 单元测试 | 创建模板-串行审批配置 | 1. 数据库连接正常<br>2. 用户有权限 | {<br>  "sign_type": "serial_sign"<br>} | HTTP 200 | 1. joint_sign_rate应为0或nil<br>2. 配置正确存储 | 高 |
| TEMPLATE-005 | 单元测试 | 查询模板详情 | 1. 模板已存在<br>2. 用户已认证 | GET /ticket_templates/1 | HTTP 200<br>返回完整的模板详情 | 1. 返回正确的模板信息<br>2. 包含所有关联配置<br>3. 内置模板标记正确 | 高 |
| TEMPLATE-006 | 单元测试 | 分页查询模板 | 1. 多个模板存在<br>2. 用户已认证 | GET /ticket_templates?page=1&size=10 | HTTP 200<br>{ "total": 15, "list": [...] } | 1. 正确分页<br>2. total计数准确<br>3. 排序正确 | 中 |
| TEMPLATE-007 | 单元测试 | 更新模板 | 1. 模板已存在<br>2. 用户有权限更新 | PUT /ticket_templates/1<br>{ "memo": "更新备注" } | HTTP 200<br>{ "id": 1 } | 1. 备注字段正确更新<br>2. 其他字段保持不变<br>3. UpdatedAt时间戳更新 | 高 |
| TEMPLATE-008 | 单元测试 | 更新内置模板 | 1. 内置模板存在<br>2. 用户已认证 | PUT /ticket_templates/1<br>{ "name": "修改名称" } | HTTP 403<br>{ "error": "cannot modify built-in template" } | 1. 返回权限错误<br>2. 模板数据不变 | 高 |
| TEMPLATE-009 | 单元测试 | 删除模板 | 1. 非内置模板存在<br>2. 用户有权限删除 | DELETE /ticket_templates/2 | HTTP 200<br>{ "id": 2 } | 1. 模板标记为软删除<br>2. 关联数据正确处理 | 高 |
| TEMPLATE-010 | 单元测试 | 删除内置模板 | 1. 内置模板存在<br>2. 用户已认证 | DELETE /ticket_templates/1 | HTTP 403<br>{ "error": "cannot delete built-in template" } | 1. 返回权限错误<br>2. 模板数据不变 | 高 |
| TEMPLATE-011 | 单元测试 | 模板筛选查询 | 1. 多个模板存在<br>2. 用户已认证 | GET /ticket_templates?name=测试&builtin=true | HTTP 200<br>返回匹配的模板列表 | 1. 正确应用筛选条件<br>2. 结果集符合预期 | 中 |
| TEMPLATE-012 | 集成测试 | 使用模板创建工单 | 1. 模板配置完整<br>2. 用户有权限 | 1. 创建工单使用该模板 | HTTP 200 | 1. 工单正确应用模板配置<br>2. 初始步骤和操作人正确 | 高 |
| TEMPLATE-013 | 边界测试 | 模板步骤数量上限 | 1. 数据库连接正常<br>2. 用户有权限 | 创建含20+步骤的模板 | HTTP 200 | 1. 模板成功创建<br>2. 所有步骤正确存储 | 低 |
| TEMPLATE-014 | 边界测试 | 模板名称超长 | 1. 数据库连接正常<br>2. 用户有权限 | { "name": "超长名称..."(100+字符) } | HTTP 400<br>{ "error": "name exceeds max length" } | 1. 返回验证错误<br>2. 无模板创建 | 中 |

## 3. 用户权限与认证测试用例

| 用例ID | 测试类别 | 测试目标 | 前置条件 | 输入数据 | 预期行为/输出 | 验证点 | 优先级 |
|--------|----------|----------|----------|----------|---------------|--------|--------|
| AUTH-001 | 单元测试 | 有效JWT访问 | 1. 用户已登录<br>2. 有效token | Authorization: Bearer <valid_token> | HTTP 200 | 1. 请求成功处理<br>2. 用户身份正确识别 | 高 |
| AUTH-002 | 单元测试 | 无效JWT访问 | 1. 无有效会话 | Authorization: Bearer <invalid_token> | HTTP 401 | 1. 返回401错误<br>2. 无业务逻辑执行 | 高 |
| AUTH-003 | 单元测试 | 无token访问 | 1. 无有效会话 | 无Authorization头 | HTTP 401 | 1. 返回401错误<br>2. 无业务逻辑执行 | 高 |
| AUTH-004 | 单元测试 | 过期token访问 | 1. token已过期 | Authorization: Bearer <expired_token> | HTTP 401 | 1. 返回token过期错误<br>2. 无业务逻辑执行 | 高 |
| AUTH-005 | 单元测试 | 创建工单-权限不足 | 1. 普通用户<br>2. 模板存在 | 创建工单请求 | HTTP 403 | 1. 返回权限错误<br>2. 无工单创建 | 高 |
| AUTH-006 | 单元测试 | 审批工单-权限不足 | 1. 工单存在<br>2. 用户不是当前步骤操作人 | 审批请求 | HTTP 403 | 1. 返回权限错误<br>2. 工单状态不变 | 高 |
| AUTH-007 | 集成测试 | RBAC权限验证 | 1. 多角色用户<br>2. 不同权限配置 | 各类操作请求 | 按角色返回相应结果 | 1. 权限控制准确<br>2. 无越权操作 | 高 |

## 4. 数据库与异常处理测试用例

| 用例ID | 测试类别 | 测试目标 | 前置条件 | 输入数据 | 预期行为/输出 | 验证点 | 优先级 |
|--------|----------|----------|----------|----------|---------------|--------|--------|
| DB-001 | 单元测试 | 数据库连接中断恢复 | 1. 模拟数据库断开 | 业务操作请求 | 1. 短暂失败<br>2. 恢复后正常 | 1. 有重试机制<br>2. 优雅降级 | 高 |
| DB-002 | 单元测试 | 事务回滚测试 | 1. 模拟部分操作失败 | 创建工单(含多表操作) | HTTP 500<br>回滚所有操作 | 1. 无部分数据创建<br>2. 数据库一致性保持 | 高 |
| DB-003 | 单元测试 | 唯一约束冲突 | 1. 工单号已存在 | 创建同名工单 | HTTP 400<br>{ "error": "order_num already exists" } | 1. 返回明确错误<br>2. 无数据创建 | 中 |
| DB-004 | 单元测试 | 外键约束验证 | 1. 模板不存在 | 创建工单(引用不存在模板) | HTTP 400<br>{ "error": "template not found" } | 1. 返回模板不存在错误 | 高 |

## 5. 性能与负载测试用例

| 用例ID | 测试类别 | 测试目标 | 前置条件 | 输入数据 | 预期行为/输出 | 验证点 | 优先级 |
|--------|----------|----------|----------|----------|---------------|--------|--------|
| PERF-001 | 性能测试 | 单接口吞吐量 | 1. 系统空载 | 100并发创建工单 | RPS ≥ 50<br>错误率 < 0.1% | 1. 满足性能指标<br>2. 资源使用合理 | 高 |
| PERF-002 | 性能测试 | 复杂流程处理 | 1. 系统空载 | 50并发完整工单流程 | 平均响应时间 ≤ 200ms | 1. 复杂操作性能达标<br>2. 无超时 | 高 |
| PERF-003 | 负载测试 | 逐步加压测试 | 1. 系统空载 | 从100到1000并发逐步加压 | 1. 平滑性能下降<br>2. 无服务崩溃 | 1. 可预测的性能曲线<br>2. 有明确容量上限 | 高 |
| PERF-004 | 负载测试 | 持续负载稳定性 | 1. 系统运行中 | 500并发持续30分钟 | 错误率 < 0.5%<br>内存稳定 | 1. 无内存泄漏<br>2. 无连接泄漏 | 高 |
| PERF-005 | 压力测试 | 极限压力测试 | 1. 系统运行中 | 2000并发短时冲击 | 1. 服务降级<br>2. 快速恢复 | 1. 有熔断机制<br>2. 无数据损坏 | 中 |
| PERF-006 | 容量测试 | 数据库容量 | 1. 100万+工单记录 | 查询操作 | 响应时间 ≤ 500ms | 1. 大数据量查询性能<br>2. 索引有效性 | 中 |
| PERF-007 | 混合测试 | 混合场景测试 | 1. 系统运行中 | 70%查询+20%创建+10%更新 | 系统稳定响应 | 1. 混合负载下性能达标<br>2. 无死锁 | 高 |
| PERF-008 | 稳定性测试 | 7x24小时稳定性 | 1. 系统启动 | 模拟真实流量模式 | 1. 无服务中断<br>2. 无性能下降 | 1. 长期运行稳定性<br>2. 资源使用趋势 | 高 |
## 5. 测试工具和框架

### 5.1 推荐测试工具
```go
// 测试框架
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
    "github.com/google/go-cmp/cmp"
)

// 测试数据库
import (
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/golang-migrate/migrate/v4"
)

// HTTP测试
import (
    "net/http/httptest"
    "github.com/gin-gonic/gin"
)
```

### 5.2 测试配置
.cmd/open_ticket_system/conf/config.dev.yaml

## 6. 测试命令

```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 查看覆盖率报告
go tool cover -html=coverage.out

# 运行基准测试
go test -bench=. ./...

# 运行特定测试
go test -run TestFunctionName ./...
```

## 7. 质量检查标准

### 7.1 覆盖率要求
- 单元测试覆盖率 > 80%
- 集成测试覆盖率 > 60%
- 关键路径覆盖率 > 95%

### 7.2 性能要求
- 单元测试执行时间 < 1秒
- 集成测试执行时间 < 30秒

### 7.3 质量标准
- 所有测试用例通过
- 无测试代码重复
- 测试命名清晰明确
- 测试数据独立隔离

## 8. 项目配置信息
- GitHub仓库：FatWang1/open_ticket_system
- 默认分支：master
- Go版本：1.24
- 许可证：AGPL-3.0