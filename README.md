[![GoDoc](https://pkg.go.dev/badge/github.com/FatWang1/open_ticket_system?utm_source=godoc)](https://pkg.go.dev/github.com/FatWang1/open_ticket_system)
[![Go Report Card](https://goreportcard.com/badge/github.com/fatwang1/open_ticket_system)](https://goreportcard.com/report/github.com/FatWang1/open_ticket_system)
[![codecov](https://codecov.io/github/FatWang1/open_ticket_system/branch/master/graph/badge.svg?token=2XWEF1Z3ZI)](https://codecov.io/github/FatWang1/open_ticket_system)
![GitHub License](https://img.shields.io/github/license/fatwang1/open_ticket_system)

# Open Ticket System

一个基于Go语言的企业级工单审批系统，支持多种审批模式（串行、并行、联合审批）。

## 🚀 特性

- **多种审批模式**: 支持串行审批、并行审批、联合审批
- **灵活模板配置**: 可自定义工单流程和步骤
- **完整权限控制**: 基于角色的权限管理系统
- **RESTful API**: 标准的REST API接口
- **数据库支持**: 使用MySQL作为主数据库，支持自动迁移
- **高性能**: 基于Gin框架，性能优异

## 🏗️ 技术架构

- **后端框架**: Gin
- **数据库**: MySQL + GORM
- **核心引擎**: [punched-tape](https://github.com/FatWang1/punched-tape) - 工单流程引擎
- **配置管理**: Viper
- **验证**: go-playground/validator
- **跨域**: gin-contrib/cors

## 📁 项目结构

```
open_ticket_system/
├── .ai/                          # AI指导文档
├── cmd/open_ticket_system/       # 主程序入口
│   ├── conf/                     # 配置文件
│   └── main.go                   # 主程序
├── internal/                     # 内部包
│   ├── client/                   # 客户端
│   │   └── database/             # 数据库连接
│   └── models/                   # 数据模型
├── ticket/                       # 工单管理模块
│   ├── controllor.go             # 控制器
│   ├── route.go                  # 路由
│   └── service.go                # 服务层
├── ticket_template/              # 工单模板管理模块
│   ├── controllor.go             # 控制器
│   ├── route.go                  # 路由
│   └── service.go                # 服务层
└── go.mod                        # Go模块文件
```

## 🚀 快速开始

### 1. 环境要求

- Go 1.24+
- MySQL 5.7+

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 配置数据库

编辑 `cmd/open_ticket_system/conf/config.dev.yaml` 文件：

```yaml
database:
  host: localhost
  port: 3306
  username: your_username
  password: your_password
  database: open_ticket_system
```

### 4. 运行项目

```bash
go run ./cmd/open_ticket_system/main.go
```

或者构建后运行：

```bash
go build ./cmd/open_ticket_system
./open_ticket_system
```

### 5. 访问API

- 健康检查: `GET /health`
- API文档: `GET /`

## 📚 API接口

### 工单管理

- `POST /tickets` - 创建工单
- `GET /tickets` - 查询工单列表
- `GET /tickets/:id` - 获取工单详情
- `PUT /tickets/:id` - 更新工单
- `DELETE /tickets/:id` - 删除工单
- `POST /tickets/:id/approval` - 工单审批
- `POST /tickets/:id/close` - 关闭工单

### 工单模板管理

- `POST /ticket_templates` - 创建模板
- `GET /ticket_templates` - 查询模板列表
- `GET /ticket_templates/:id` - 获取模板详情
- `PUT /ticket_templates/:id` - 更新模板
- `DELETE /ticket_templates/:id` - 删除模板

## 🔧 审批模式

### 1. 串行审批 (Serial Sign)
- 按顺序逐个审批
- 前一个审批人通过后，下一个才能审批

### 2. 并行审批 (Jointly Sign)
- 多个审批人同时审批
- 需要达到指定通过率才能进入下一步

### 3. 任意人审批 (Anyone Sign)
- 任意一个审批人通过即可
- 适用于简单审批场景

## 🗄️ 数据库设计

系统包含以下核心表：

- `tickets` - 工单主表
- `ticket_operators` - 工单操作人表
- `ticket_operated_users` - 工单已操作用户表
- `ticket_templates` - 工单模板表
- `template_end_steps` - 模板结束节点表
- `step_configs` - 步骤配置表
- `step_operators` - 步骤操作人表
- `next_steps` - 下一步配置表

## 🧪 测试

运行测试：

```bash
go test ./...
```

## 📝 开发说明

### 核心业务逻辑

本系统完全遵循 [punched-tape](https://github.com/FatWang1/punched-tape) 库的实现，该库提供了：

- 工单模板验证
- 工单审批流程管理
- 多种审批模式支持
- 状态机管理

### 扩展开发

1. 在 `internal/models/` 中添加新的数据模型
2. 在对应的模块中添加服务层逻辑
3. 在控制器中实现API接口
4. 在路由中注册新的端点

## 🤝 贡献

欢迎提交Issue和Pull Request！

## 📄 许可证

本项目采用 AGPL-3.0 许可证。

## 🔗 相关链接

- [punched-tape](https://github.com/FatWang1/punched-tape) - 工单流程引擎
- [Gin](https://github.com/gin-gonic/gin) - Web框架
- [GORM](https://gorm.io/) - ORM库
