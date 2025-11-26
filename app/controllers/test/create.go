package test

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/global"
	"github.com/ljzh152535/go-framework/goadmin-gin/ReturnData"
)

func Create(r *gin.Context) {
	global.GVA_LOG.Debug("创建test")

	var a = "创建test成功"
	ReturnData.OKWithMessage(a, r)
}
