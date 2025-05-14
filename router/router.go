package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/renoelis/returnall-go/controller"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 创建默认的gin引擎
	r := gin.Default()

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