package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger 自定义日志记录器结构体
type Logger struct {
	name      string
	logger    *log.Logger
	logFile   *os.File
	isConsole bool
}

// LogEntry 日志记录项结构体
type LogEntry struct {
	Time        string      `json:"time"`
	Name        string      `json:"name"`
	Level       string      `json:"level"`
	Message     string      `json:"message"`
	RequestInfo interface{} `json:"request_info,omitempty"`
	Exception   string      `json:"exception,omitempty"`
}

// NewLogger 创建新的日志记录器
func NewLogger(name string, logFile string, isConsole bool) (*Logger, error) {
	logger := &Logger{
		name:      name,
		isConsole: isConsole,
	}

	if logFile != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(logFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("无法创建日志目录: %v", err)
		}

		// 打开日志文件
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return nil, fmt.Errorf("无法打开日志文件: %v", err)
		}
		logger.logFile = file

		// 同时输出到控制台和文件
		var writer io.Writer
		if isConsole {
			writer = io.MultiWriter(os.Stdout, file)
		} else {
			writer = file
		}

		logger.logger = log.New(writer, "", 0)
	} else if isConsole {
		// 仅输出到控制台
		logger.logger = log.New(os.Stdout, "", 0)
	} else {
		return nil, fmt.Errorf("没有指定日志输出目标")
	}

	return logger, nil
}

// formatLogEntry 格式化日志条目为JSON
func (l *Logger) formatLogEntry(level, message string, requestInfo interface{}, exception string) string {
	entry := LogEntry{
		Time:    time.Now().Format(time.RFC3339),
		Name:    l.name,
		Level:   level,
		Message: message,
	}

	if requestInfo != nil {
		entry.RequestInfo = requestInfo
	}

	if exception != "" {
		entry.Exception = exception
	}

	// 转换为JSON格式
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// 如果JSON转换失败，使用基本格式
		return fmt.Sprintf("{\"time\":\"%s\",\"name\":\"%s\",\"level\":\"%s\",\"message\":\"%s\",\"error\":\"无法序列化日志\"}", 
			entry.Time, entry.Name, entry.Level, entry.Message)
	}

	return string(jsonData)
}

// Info 记录信息级别日志
func (l *Logger) Info(message string, requestInfo interface{}) {
	logEntry := l.formatLogEntry("INFO", message, requestInfo, "")
	l.logger.Println(logEntry)
}

// Error 记录错误级别日志
func (l *Logger) Error(message string, err error, requestInfo interface{}) {
	var exception string
	if err != nil {
		exception = err.Error()
	}
	logEntry := l.formatLogEntry("ERROR", message, requestInfo, exception)
	l.logger.Println(logEntry)
}

// Close 关闭日志文件
func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
} 