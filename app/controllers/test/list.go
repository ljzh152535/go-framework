package test

import (
	"github.com/gin-gonic/gin"
	"github.com/ljzh152535/go-framework/app/global"
	"github.com/ljzh152535/go-framework/goadmin-gin/ReturnData"
)

func List(r *gin.Context) {
	global.GVA_LOG.Debug("查询test列表")

	var a = []string{"test列表1", "test列表2"}
	b := ReturnData.GetResData(a, ReturnData.ITEMS)
	ReturnData.OKWithData(b, r)
}
