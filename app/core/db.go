package core

import (
	"fmt"
	"github.com/ljzh152535/go-framework/app/global"
	"github.com/ljzh152535/go-framework/app/models/user"
	"os"
)

// RegisterTables 注册数据库表专用
func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(
		&user.User{},
	)
	if err != nil {
		//global.GVA_LOG.Error("register table failed", zap.Error(err))
		fmt.Println("register table failed", err.Error())
		os.Exit(0)
	}
	//global.GVA_LOG.Info("register table success")
	fmt.Println("register table success")
}
