package goadmin_db

import (
	"fmt"
	goadmin_dbmodel "github.com/ljzh152535/go-framework/goadmin-db/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	dbInstance = make(map[string]*gorm.DB)
	dbLocker   sync.RWMutex
)

// GormMysql 初始化Mysql数据库
func loadDB(m goadmin_dbmodel.DBItemConf, confLog goadmin_dbmodel.DBLog, key string) *gorm.DB {

	if m.Database == "" {
		return nil
	}

	dbLocker.RLock()
	db, ok := dbInstance[key]
	if ok {
		dbLocker.RUnlock()
		return db
	}
	dbLocker.RUnlock()

	// 原子操作,代替锁

	dbLocker.Lock()
	defer dbLocker.Unlock()
	if _, exist := dbInstance[key]; exist {
		return dbInstance[key]
	}

	dbInstance[key] = getDBInstance(m, confLog)
	return dbInstance[key]
}

func getDBInstance(m goadmin_dbmodel.DBItemConf, confLog goadmin_dbmodel.DBLog) *gorm.DB {
	if m.Timeout <= 0 {
		m.Timeout = 5000
	}
	if m.WriteTimeOut <= 0 {
		m.WriteTimeOut = 5000
	}
	if m.ReadTimeOut <= 0 {
		m.ReadTimeOut = 5000
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s&timeout=%dms&writeTimeout=%dms&readTimeout=%dms",
		m.Username,
		m.Password,
		m.Host,
		m.Port,
		m.Database,
		m.Config,
		m.Timeout,
		m.WriteTimeOut,
		m.ReadTimeOut,
	)

	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		SkipInitializeWithVersion: false, // 根据版本自动配置

	}

	var l logger.LogLevel
	var dbLogger *log.Logger
	var gLogger logger.Interface

	if confLog.Enable {
		switch confLog.Level {
		case "silent":
			l = logger.Silent
		case "error":
			l = logger.Error
		case "info":
			l = logger.Info
		default:
			l = logger.Warn
		}

		if confLog.Type == "file" {
			logPath := confLog.Path
			if !filepath.IsAbs(confLog.Path) {
				//logPath = global.GVA_HOME_PATH + "/" + confLog.Path
				logPath = "./" + "logs" + "/" + confLog.Path
			}

			// 存文件
			f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}
			dbLogger = log.New(f, "\r\n", log.LstdFlags)
		} else if confLog.Type == "stdout" {
			dbLogger = log.New(os.Stdout, "\r\n", log.LstdFlags)
		}
		gLogger = New(dbLogger, logger.Config{
			LogLevel:                  l,                      // 日志级别
			SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
			IgnoreRecordNotFoundError: false,
			Colorful:                  true, // 着色打印
		}, "dev")
	}

	db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: gLogger,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	// 设置连接池的最大空闲连接数
	sqlDB.SetMaxIdleConns(m.MaxIdleConns)

	// 设置连接池的最大连接数
	sqlDB.SetMaxOpenConns(m.MaxOpenConns)

	// 设置连接最大存活时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

// RegisterTables 注册数据库表专用
//func RegisterTables(db *gorm.DB) {
//	err := db.AutoMigrate(
//		&test.User{},
//	)
//	if err != nil {
//		//global.GVA_LOG.Error("register table failed", zap.Error(err))
//		fmt.Println("register table failed", err.Error())
//		os.Exit(0)
//	}
//	//global.GVA_LOG.Info("register table success")
//	fmt.Println("register table success")
//}
