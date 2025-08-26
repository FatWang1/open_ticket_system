package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLConfig MySQL配置
type MySQLConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	MaxRetries      int    `yaml:"max_retries"`       // 最大重试次数
	RetryDelay      int    `yaml:"retry_delay"`       // 重试延迟(秒)
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"` // 连接最大生命周期(小时)
	MaxIdleConns    int    `yaml:"max_idle_conns"`    // 最大空闲连接数
	MaxOpenConns    int    `yaml:"max_open_conns"`    // 最大打开连接数
}

// MySQLClient MySQL客户端
type MySQLClient struct {
	DB     *gorm.DB
	config *MySQLConfig
}

// NewMySQLClient 创建MySQL客户端
func NewMySQLClient(config *MySQLConfig) (*MySQLClient, error) {
	client := &MySQLClient{config: config}

	// 尝试连接数据库，支持重试
	var err error
	for i := 0; i <= config.MaxRetries; i++ {
		if i > 0 {
			log.Printf("Retrying database connection, attempt %d/%d", i+1, config.MaxRetries+1)
			time.Sleep(time.Duration(config.RetryDelay) * time.Second)
		}

		err = client.connect(config)
		if err == nil {
			break
		}

		log.Printf("Database connection attempt %d failed: %v", i+1, err)
	}

	// 即使连接失败，也返回客户端结构，但标记为未连接
	if err != nil {
		log.Printf("MySQL client created but connection failed: %v", err)
		return client, fmt.Errorf("failed to connect to MySQL after %d attempts: %w", config.MaxRetries+1, err)
	}

	// 自动迁移表结构
	if err := client.AutoMigrate(); err != nil {
		log.Printf("MySQL client connected but migration failed: %v", err)
		return client, fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("MySQL client initialized successfully")
	return client, nil
}

// connect 建立数据库连接
func (c *MySQLClient) connect(config *MySQLConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 设置连接池参数
	if config.MaxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	} else {
		sqlDB.SetMaxIdleConns(10)
	}

	if config.MaxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	} else {
		sqlDB.SetMaxOpenConns(100)
	}

	if config.ConnMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Hour)
	} else {
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	c.DB = db
	return nil
}

// Reconnect 重新连接数据库
func (c *MySQLClient) Reconnect() error {
	log.Println("Attempting to reconnect to database...")
	return c.connect(c.config)
}

// HealthCheck 健康检查
func (c *MySQLClient) HealthCheck() error {
	if c.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	sqlDB, err := c.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

// GetConnectionStats 获取连接池统计信息
func (c *MySQLClient) GetConnectionStats() map[string]interface{} {
	if c.DB == nil {
		return nil
	}

	sqlDB, err := c.DB.DB()
	if err != nil {
		return nil
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// AutoMigrate 自动迁移表结构
func (c *MySQLClient) AutoMigrate() error {
	return c.DB.AutoMigrate(
		&models.User{},
		&models.Ticket{},
		&models.TicketOperator{},
		&models.TicketOperatedUser{},
		&models.TicketTemplate{},
		&models.TemplateEndStep{},
		&models.StepConfig{},
		&models.StepOperator{},
		&models.NextStep{},
	)
}

// GetDB 获取数据库实例
func (c *MySQLClient) GetDB() *gorm.DB {
	return c.DB
}

// Close 关闭数据库连接
func (c *MySQLClient) Close() error {
	if c.DB == nil {
		return nil
	}

	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
