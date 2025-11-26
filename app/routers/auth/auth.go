package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/controllers/auth"
)

// 实现登录接口
func login(authGroup *gin.RouterGroup) {
	authGroup.POST("/login", auth.Login)
}

// 实现退出接口
func logout(authGroup *gin.RouterGroup) {
	authGroup.GET("/logout", auth.Logout)
}
func RegisterSubRouters(g *gin.RouterGroup) {
	// 配置登录功能的路由策略
	authGroup := g.Group("/auth")
	// 登录的功能
	login(authGroup)
	logout(authGroup)

}
