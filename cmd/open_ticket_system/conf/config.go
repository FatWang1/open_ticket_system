package conf

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/FatWang1/fatwang-go-utils/utils"
	"github.com/FatWang1/open_ticket_system/internal/client/database"
)

// Config 应用配置结构
type Config struct {
	Server   *ServerConfig    `yaml:"server"`
	Database *DatabaseConfig  `yaml:"database"`
	Log      *utils.LogConfig `yaml:"log"`
	JWT      *JWTConfig       `yaml:"jwt"`
	Redis    *RedisConfig     `yaml:"redis"`
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
	Port     int    `yaml:"port"`
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

// JWTConfig JWT配置
type JWTConfig struct {
	Secret             string `yaml:"secret"`
	AccessTokenExpire  int    `yaml:"access_token_expire" destructure:"access_token_expire"`   // access token过期时间(小时)
	RefreshTokenExpire int    `yaml:"refresh_token_expire" destructure:"refresh_token_expire"` // refresh token过期时间(小时，默认7天)
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
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(fmt.Sprintf("config.%s.yaml", gin.Mode()))
	v.SetConfigType("yaml")

	// 设置默认值
	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
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
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.timeout", 30)

	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.charset", "utf8mb4")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.database", 0)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")
	v.SetDefault("log.output", "stdout")
	v.SetDefault("log.max_size", 100)
	v.SetDefault("log.max_backups", 3)
	v.SetDefault("log.max_age", 28)
	v.SetDefault("log.compress", true)

	v.SetDefault("jwt.access_token_expire", 1)
	v.SetDefault("jwt.refresh_token_expire", 168)
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
