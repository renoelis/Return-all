package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/renoelis/returnall-go/config"
	"github.com/renoelis/returnall-go/controller"
	"github.com/renoelis/returnall-go/router"
)

func main() {
	// 设置Gin为生产环境模式
	gin.SetMode(gin.ReleaseMode)

	// 禁用Gin默认的控制台颜色
	gin.DisableConsoleColor()

	// 自定义简化的请求日志中间件
	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		// 不输出路由信息
	}

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
