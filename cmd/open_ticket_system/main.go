// @title           Open Ticket System API
// @version         1.0
// @description     这是一个完整的工单管理系统API，支持工单创建、审批、模板管理等功能。
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description 请输入JWT token，格式：Bearer {token}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FatWang1/open_ticket_system/cmd/open_ticket_system/conf"
	"github.com/FatWang1/open_ticket_system/internal/client/database"
	"github.com/FatWang1/open_ticket_system/internal/middleware"
	"github.com/FatWang1/open_ticket_system/internal/utils"
	"github.com/FatWang1/open_ticket_system/ticket"
	"github.com/FatWang1/open_ticket_system/ticket_template"

	_ "github.com/FatWang1/open_ticket_system/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// 初始化日志
	utils.InitLogger()
	logger := utils.GetLogger()

	// 加载配置
	configPath := "cmd/open_ticket_system/conf/"

	config, err := conf.LoadConfig(configPath)
	if err != nil {
		logger.Printf("Failed to load config: %v", err)
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库
	mysqlClient, err := database.NewMySQLClient(config.GetMySQLConfig())
	if err != nil {
		logger.Printf("Failed to initialize MySQL client: %v", err)
		log.Fatalf("Failed to initialize MySQL client: %v", err)
	}
	defer mysqlClient.Close()

	// 创建Gin引擎
	engine := gin.Default()

	// 添加中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middleware.CORSMiddleware())

	// 创建JWT配置
	jwtConfig := &middleware.JWTConfig{
		Secret: config.JWT.Secret,
	}

	// 注册路由
	registerRoutes(engine, mysqlClient, jwtConfig)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port),
		Handler: engine,
	}

	// 启动服务器
	go func() {
		log.Printf("Server starting on %s:%d", config.Server.Host, config.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Server.Timeout)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// registerRoutes 注册路由
func registerRoutes(engine *gin.Engine, mysqlClient *database.MySQLClient, jwtConfig *middleware.JWTConfig) {
	// 健康检查
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
	})

	// Swagger文档
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API版本组
	v1 := engine.Group("/api/v1")
	{
		// 工单相关路由
		ticketController := ticket.NewTicketController(mysqlClient.DB)
		ticketController.Register(v1)

		// 工单模板相关路由
		templateController := ticket_template.NewTicketTemplateController(mysqlClient.DB)
		templateController.Register(v1)
	}

	// 需要认证的路由组
	auth := engine.Group("/api/v1")
	auth.Use(middleware.JWTAuthMiddleware(jwtConfig))
	{
		// 这里可以添加需要认证的路由
		auth.GET("/profile", func(c *gin.Context) {
			userID := c.GetUint("user_id")
			username := c.GetString("username")
			role := c.GetString("role")
			c.JSON(http.StatusOK, gin.H{
				"user_id":  userID,
				"username": username,
				"role":     role,
			})
		})
	}
}
