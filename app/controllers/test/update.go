package test

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/global"
	"github.com/ljzh152535/go-framework/goadmin-gin/ReturnData"
)

func Update(r *gin.Context) {
	global.GVA_LOG.Debug("更新test")
	var a = "更新test1成功"
	b := ReturnData.GetResData(a, ReturnData.ITEMS)
	ReturnData.OKWithData(b, r)
}
