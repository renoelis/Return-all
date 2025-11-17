package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/renoelis/returnall-go/config"
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

	// 设置路由
	r := router.SetupRouter()

	// 获取服务端口
	port := config.GetPort()

	// 启动服务
	fmt.Printf("服务启动，监听端口: %s\n", port)
	if err := r.Run(":" + port); err != nil {
		panic(fmt.Sprintf("服务启动失败: %v", err))
	}
}
