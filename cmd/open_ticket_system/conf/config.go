package conf

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/FatWang1/open_ticket_system/internal/client/database"
)

// Config 应用配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Log      LogConfig      `yaml:"log"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port    int    `yaml:"port" default:"8080"`
	Host    string `yaml:"host" default:"localhost"`
	Timeout int    `yaml:"timeout" default:"30"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `yaml:"driver" default:"mysql"`
	Host     string `yaml:"host" default:"localhost"`
	Port     int    `yaml:"port" default:"3306"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset" default:"utf8mb4"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host" default:"localhost"`
	Port     int    `yaml:"port" default:"6379"`
	Password string `yaml:"password"`
	Database int    `yaml:"database" default:"0"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level" default:"info"`
	Format     string `yaml:"format" default:"json"`
	Output     string `yaml:"output" default:"stdout"`
	MaxSize    int    `yaml:"max_size" default:"100"`
	MaxBackups int    `yaml:"max_backups" default:"3"`
	MaxAge     int    `yaml:"max_age" default:"28"`
	Compress   bool   `yaml:"compress" default:"true"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireTime int    `yaml:"expire_time" default:"24"` // 小时
}

func (j *JWTConfig) loadSecret() {
	envSecret := os.Getenv("JWT_SECRET")
	if envSecret == "" {
		log.Println("JWT_SECRET environment variable not set, using default value")
		return
	}
	j.Secret = envSecret
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {

	viper.AddConfigPath(configPath)
	viper.SetConfigName(fmt.Sprintf("config.%s.yaml", gin.Mode()))
	viper.SetConfigType("yaml")

	// 设置默认值
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	config.JWT.loadSecret()
	log.Println("Configuration loaded successfully")
	return &config, nil
}

// setDefaults 设置默认值
func setDefaults() {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.timeout", 30)

	viper.SetDefault("database.driver", "mysql")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.charset", "utf8mb4")

	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", 6379)
	viper.SetDefault("redis.database", 0)

	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "json")
	viper.SetDefault("log.output", "stdout")
	viper.SetDefault("log.max_size", 100)
	viper.SetDefault("log.max_backups", 3)
	viper.SetDefault("log.max_age", 28)
	viper.SetDefault("log.compress", true)

	viper.SetDefault("jwt.expire_time", 24)
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}

	if config.Database.Username == "" {
		return fmt.Errorf("database username is required")
	}

	if config.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}

	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	return nil
}

// GetMySQLConfig 获取MySQL配置
func (c *Config) GetMySQLConfig() *database.MySQLConfig {
	return &database.MySQLConfig{
		Host:     c.Database.Host,
		Port:     c.Database.Port,
		Username: c.Database.Username,
		Password: c.Database.Password,
		Database: c.Database.Database,
		Charset:  c.Database.Charset,
	}
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Database.Username,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Database,
		c.Database.Charset,
	)
}
