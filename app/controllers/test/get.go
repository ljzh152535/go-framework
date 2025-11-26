package test

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/global"
	"github.com/ljzh152535/go-framework/goadmin-gin/ReturnData"
)

func Get(r *gin.Context) {
	global.GVA_LOG.Debug("获取test详情")

	var a = "获取test1成功"
	b := ReturnData.GetResData(a, ReturnData.ITEMS)
	ReturnData.OKWithData(b, r)
}
