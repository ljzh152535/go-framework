package main

import (
	"github.com/ljzh152535/go-framework/app/core"
	"github.com/ljzh152535/go-framework/app/global"
	"github.com/ljzh152535/go-framework/app/routers"
	goadmin_db "github.com/ljzh152535/go-framework/goadmin-db"
	goadmin_logrus "github.com/ljzh152535/go-framework/goadmin-logrus"
	goadmin_viper "github.com/ljzh152535/go-framework/goadmin-viper"
	"github.com/rs/xid"
)

func main() {

	global.GVA_VP = goadmin_viper.InitGinViper(&global.GVA_CONFIG)
	traceId := xid.New().String()
	global.GVA_CONFIG.LOG.TraceID = traceId
	//textTraceID := "[" + global.GVA_CONFIG.LOG.TraceID + "] "

	global.GVA_LOG = goadmin_logrus.InitLogrushook(global.GVA_CONFIG.LOG)

	global.GVA_DB = goadmin_db.InitDB(global.GVA_CONFIG.DB)
	//initialize.DBList()
	if global.GVA_DB != nil {
		core.RegisterTables() // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.GVA_DB.DB()
		defer db.Close()
	}

	// 路由初始化
	router := routers.InitRouter()
	addr := global.GVA_CONFIG.System.Addr()

	global.GVA_LOG.Infof("admin-backend 运行在: %s", addr)
	router.Run(addr)
}
