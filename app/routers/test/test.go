package test

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/controllers/test"
)

func get(deploymentGroup *gin.RouterGroup) {
	deploymentGroup.GET("/get", test.Get)
}

func list(deploymentGroup *gin.RouterGroup) {
	deploymentGroup.GET("/list", test.List)
}

func add(deploymentGroup *gin.RouterGroup) {
	deploymentGroup.POST("/create", test.Create)
}
func update(deploymentGroup *gin.RouterGroup) {
	deploymentGroup.POST("/update", test.Update)
}
func delete(deploymentGroup *gin.RouterGroup) {
	deploymentGroup.POST("/delete", test.Delete)
}
func deleteList(deploymentGroup *gin.RouterGroup) {
	deploymentGroup.POST("/deleteList", test.DeleteList)
}

func RegisterSubRouter(g *gin.RouterGroup) {
	deploymentGroup := g.Group("/test")
	get(deploymentGroup)
	list(deploymentGroup)
	add(deploymentGroup)
	update(deploymentGroup)
	delete(deploymentGroup)
	deleteList(deploymentGroup)
}
