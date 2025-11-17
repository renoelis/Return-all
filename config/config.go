package config

import (
	"os"
	"strconv"
)

// 默认配置
const (
	DefaultPort         = "3006"
	DefaultMaxBodyBytes = 2 * 1024 * 1024
)

// GetPort 获取服务端口号
func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return DefaultPort
	}
	return port
}

// GetMaxBodyBytes 获取请求体最大字节数
// 优先从环境变量 MAX_BODY_BYTES 读取，读取失败或非法时回退到默认值
func GetMaxBodyBytes() int64 {
	val := os.Getenv("MAX_BODY_BYTES")
	if val == "" {
		return DefaultMaxBodyBytes
	}
	n, err := strconv.ParseInt(val, 10, 64)
	if err != nil || n <= 0 {
		return DefaultMaxBodyBytes
	}
	return n
}
