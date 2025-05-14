package config

import (
	"os"
)

// 默认配置
const (
	DefaultPort = "3006"
)

// GetPort 获取服务端口号
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return DefaultPort
	}
	return port
}

// GetLogPath 获取日志文件路径
func GetLogPath() string {
	logPath := os.Getenv("LOG_PATH")
	if logPath == "" {
		return "logs/api.log"
	}
	return logPath
} 