package goadmin_db

import (
	"fmt"
	goadim_dbModel "github.com/ljzh152535/go-framework/goadmin-db/model"
	"gorm.io/gorm"
)

//var (
//	Config config.Model
//)

func InitDB(dbArr map[string]goadim_dbModel.DBItem) *gorm.DB {
	// 初始化数据库
	db := DBW(dbArr)
	return db
}

func DBW(dbArr map[string]goadim_dbModel.DBItem, keys ...string) *gorm.DB {
	k := "default"
	if len(keys) > 0 {
		k = keys[0]
	}
	conf, ok := dbArr[k]
	if !ok {
		panic(fmt.Sprintf("db config %s not found", k))
	}

	cacheKey := fmt.Sprintf("%s_write", k)
	return loadDB(conf.Write, conf.Log, cacheKey)
}

func DBR(dbArr map[string]goadim_dbModel.DBItem, keys ...string) *gorm.DB {
	k := "default"
	if len(keys) > 0 {
		k = keys[0]
	}
	conf, ok := dbArr[k]
	if !ok {
		panic(fmt.Sprintf("db config %s not found"))
	}
	cacheKey := fmt.Sprintf("%s_read", k)
	return loadDB(conf.Read, conf.Log, cacheKey)
}

//func Log(confLog config.Log, env string) *logrus.Entry {
//	return log.Load(homeDir, confLog, env)
//}
