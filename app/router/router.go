package router

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	v1 "github.com/seanlan/xlvein/app/api/v1"
	"go.uber.org/zap"
	"time"
)

var Router *gin.Engine

// Setup 初始化Router
func Setup(debug bool) {
	//设置启动模式
	switch debug {
	case true:
		gin.SetMode(gin.DebugMode)
		break
	case false:
		gin.SetMode(gin.ReleaseMode)
	}
	Router = gin.New()
	//设置中间件
	Router.Use(
		ginzap.Ginzap(zap.L(), time.RFC3339, false),
		ginzap.RecoveryWithZap(zap.L(), true),
	)
	//设置路由
	// websocket 连接
	wsGroup := Router.Group("ws")
	{
		wsGroup.GET("connect", v1.SocketConnect)
	}
	// api 路由
	apiGroup := Router.Group("api/v1")
	imGroup := apiGroup.Group("im")
	{
		// 创建会话
		imGroup.POST("create", v1.CreateConversation)
	}
}

func Run(addr string) {
	err := Router.Run(addr)
	if err != nil {
		return
	}
}
