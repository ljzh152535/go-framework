package test

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/global"
	"github.com/ljzh152535/go-framework/goadmin-gin/ReturnData"
)

func Delete(r *gin.Context) {
	global.GVA_LOG.Debug("删除test")
	var a = "删除test1成功"
	b := ReturnData.GetResData(a, ReturnData.ITEMS)
	ReturnData.OKWithData(b, r)
}

func DeleteList(r *gin.Context) {
	global.GVA_LOG.Debug("批量删除test")
	var a = []string{"删除test1成功", "删除test2成功"}
	b := ReturnData.GetResData(a, ReturnData.ITEMS)
	ReturnData.OKWithData(b, r)
}
