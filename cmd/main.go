package main

import (
	"fmt"
	"log"
	"os"

	"github.com/renoelis/returnall-go/config"
	"github.com/renoelis/returnall-go/controller"
	"github.com/renoelis/returnall-go/router"
)

func main() {
	// 确保日志目录存在
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatalf("无法创建日志目录: %v", err)
	}

	// 初始化日志记录器
	logPath := config.GetLogPath()
	if err := controller.InitLogger(logPath); err != nil {
		log.Fatalf("无法初始化日志记录器: %v", err)
	}

	// 设置路由
	r := router.SetupRouter()

	// 获取服务端口
	port := config.GetPort()
	
	// 启动服务
	fmt.Printf("服务启动，监听端口: %s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
} 