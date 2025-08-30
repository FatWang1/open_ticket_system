package utils

import (
	"context"
	"io"
	"log"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Logger 全局日志实例
	Logger *log.Logger
)

// InitLogger 初始化日志
func InitLogger() {
	// 配置日志轮转
	logWriter := &lumberjack.Logger{
		Filename:   "logs/open_ticket_system.log",
		MaxSize:    100,  // 每个日志文件最大100MB
		MaxBackups: 3,    // 保留3个备份文件
		MaxAge:     28,   // 保留28天
		Compress:   true, // 压缩旧文件
	}

	// 同时输出到文件和控制台
	multiWriter := io.MultiWriter(os.Stdout, logWriter)

	Logger = log.New(multiWriter, "", log.LstdFlags|log.Lshortfile)
}

// GetLogger 获取日志实例
func GetLogger() *log.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}

// LogInfo 记录信息日志
func LogInfo(ctx context.Context, msg string, fields ...interface{}) {
	GetLogger().Printf("[INFO] %s %v", msg, fields)
}

// LogError 记录错误日志
func LogError(ctx context.Context, msg string, err error, fields ...interface{}) {
	GetLogger().Printf("[ERROR] %s: %v %v", msg, err, fields)
}

// LogWarn 记录警告日志
func LogWarn(ctx context.Context, msg string, fields ...interface{}) {
	GetLogger().Printf("[WARN] %s %v", msg, fields)
}

// LogDebug 记录调试日志
func LogDebug(ctx context.Context, msg string, fields ...interface{}) {
	GetLogger().Printf("[DEBUG] %s %v", msg, fields)
}
