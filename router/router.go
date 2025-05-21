package router

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/renoelis/returnall-go/controller"
)

// 自定义简化中间件 - 只显示请求成功/失败状态
func SimpleLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)
		method := c.Request.Method
		statusCode := c.Writer.Status()

		// 简洁输出，只显示状态
		statusText := "成功"
		if statusCode >= 400 {
			statusText = "失败"
		}

		fmt.Printf("请求: %s %s | 状态: %s (%d) | 耗时: %v\n",
			method, path, statusText, statusCode, latency)
	}
}

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 创建无日志的引擎
	r := gin.New()

	// 使用自定义简化日志和恢复中间件
	r.Use(SimpleLogger())
	r.Use(gin.Recovery())

	// 允许跨域请求
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 根路径
	r.GET("/", controller.RootHandler)

	// returnAll 基本路径 - 只保留这一个路由
	r.POST("/returnAll", controller.ReturnAllRequest)

	// 通用路径处理 - 支持任意路径参数
	r.POST("/returnAll/*path", controller.ReturnAllWithAnyPath)

	return r
}
