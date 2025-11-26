package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/global"
	"github.com/ljzh152535/go-framework/app/middleware"
	"github.com/ljzh152535/go-framework/app/routers/auth"
	"github.com/ljzh152535/go-framework/app/routers/test"
	gin_middlewares "github.com/ljzh152535/go-framework/goadmin-gin/middleware"
	gin_middle_logger "github.com/ljzh152535/go-framework/goadmin-gin/middleware/logger"
	"github.com/ljzh152535/go-framework/goadmin-os"
)

func InitRouter() *gin.Engine {
	//gin.SetMode(global.GVA_CONFIG.System.Env)
	r := gin.New()

	//router := gin.Default()

	// 允许跨域
	r.Use(gin_middlewares.CORSMiddleWare)

	if global.GVA_CONFIG.System.Env != "dev" {
		r.Use(middlewares.JWTAuth)
	}

	// 使用日志中间件
	if global.GVA_CONFIG.WebServerLog.Enable {
		r.Use(gin_middle_logger.New(global.GVA_CONFIG.WebServerLog, global.GVA_CONFIG.System.Env, goadmin_os.GetHomePath()))
	}

	r.Use(gin.Recovery())

	// 路由注册
	apiGroup := r.Group("/api")
	test.RegisterSubRouter(apiGroup)
	auth.RegisterSubRouters(apiGroup)

	return r
}
