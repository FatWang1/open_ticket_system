# AI Project Guide Template

## 1. 项目基本信息
- 项目名称：open_ticket_system
- 模块路径：github.com/FatWang1/open_ticket_system
- 技术栈：Go, Gin, GORM, MySQL

## 2. 推荐项目结构

```
/{{PROJECT_ROOT}}
├── .ai
│     ├── design.md
│     ├── guide.md
│     ├── test.md
│     └── workflow.md
├── .deploy
│     ├── k8s
│     │     ├── deployment.yaml
│     │     └── docker
│     │         ├── .dockerignore
│     │         └── Dokcerfile
│     └── scripts
│         ├── build.sh
│         └── deploy.sh
├── .github
│     └── workflows
│         └── go.yml
├── .gitignore
├── LICENSE
├── README.md
├── cmd
│     └── open_ticket_system
│         ├── conf
│         │     └── config.go
│         └── main.go
├── go.mod
├── go.sum
├── internal
│     ├── client
│     │     ├── database
│     │     └── third_party
│     ├── helper                # 数据库操作包装
│     ├── manager               # 数据库操作
│     ├── middleware
│     ├── models
│     └── utils
│         └── utils.go
├── pkg
│     ├── constant.go
│     └── version
│         └── version.go
├── ticket
│     ├── controllor.go
│     ├── route.go
│     ├── service.go
│     └── validator
│         └── validator.go
└── ticket_template
    ├── controllor.go
    ├── route.go
    ├── service.go
    └── validator
        └── validator.go
```



## 3. 开发规范

### 3.1 编码规范
- 遵循Go语言官方编码规范
- 使用统一的错误处理机制
- 实现完整的日志记录
- 保持代码文档的及时更新

### 3.2 命名规范
- 包名：小写字母，简短描述性
- 文件名：小写字母，下划线分隔
- 结构体：大驼峰命名
- 方法：大驼峰命名
- 变量：小驼峰命名
- 常量：大写下划线分隔

### 3.3 导入规范
```go
import (
    // 标准库
    "fmt"
    "net/http"
    
    // 第三方库
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    
    // 内部包
    "github.com/FatWang1/open_ticket_system/internal/config"
    "github.com/FatWang1/open_ticket_system/pkg/types"
)
```

## 4. 配置管理

### 4.1 配置结构
```go
type Config struct {
    MySql     *Mysql             `yaml:"mysql"`
    Log       *utils.LogConfig  `yaml:"log"`
    AccessLog *lumberjack.Logger `yaml:"accessLog"`
}

```

### 4.2 环境配置
- 开发环境：`config.dev.yaml`

## 5. 项目配置信息
- GitHub仓库：FatWang1/open_ticket_system
- 默认分支：master
- Go版本：1.24
- 许可证：AGPL-3.0

