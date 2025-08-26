# AI Project Workflow Template

## 1. AI工作流程概述
- 项目名称：open_ticket_system
- 项目类型：Web服务
- 项目目标：open_ticket_system

## 2. AI执行流程

### Step 1: 需求解析
- 解析 `guide.md` 中的项目结构和开发规范
- 解析 `design.md` 中的架构设计和接口定义
- 确定项目技术栈和依赖关系

### Step 2: 代码生成
- 根据 `guide.md` 生成项目目录结构
- 根据 `design.md` 生成核心代码文件
- 生成配置文件、依赖管理等基础文件

### Step 3: 测试验证
- 根据 `test.md` 生成测试用例
- 执行单元测试和集成测试
- 验证代码质量和功能完整性

### Step 4: 文档生成
- 生成API文档和用户手册
- 生成部署和运维文档
- 更新项目README和变更日志

## 3. AI协作规则

### 3.1 文件依赖关系
```
guide.md → 项目结构生成
design.md → 核心代码生成
test.md → 测试用例生成
workflow.md → 执行流程控制
```

### 3.2 生成优先级
1. 基础项目结构 (guide.md)
2. 核心业务代码 (design.md)
3. 测试验证代码 (test.md)
4. 文档和配置 (workflow.md)

### 3.3 质量检查
- 代码覆盖率 > 80%
- 所有测试用例通过
- 符合编码规范
- 文档完整性检查

## 4. 待AI完善部分

### 🔧 需要AI根据业务需求生成
```
TODO: 根据guide.md和design.md生成
- [ ] 核心业务模块代码
- [ ] 数据模型和接口定义
- [ ] 配置文件和依赖管理
- [ ] API文档和示例代码
```

### 📋 需要AI根据test.md验证
```
TODO: 根据test.md执行验证
- [ ] 单元测试用例生成和执行
- [ ] 集成测试场景验证
- [ ] 性能测试基准检查
- [ ] 安全测试漏洞扫描
```

## 5. 项目配置
- 模块路径：github.com/FatWang1/open_ticket_system
- GitHub仓库：FatWang1/open_ticket_system
- 技术栈：Go, Gin, GORM, MySQL
- Go版本：1.24

---