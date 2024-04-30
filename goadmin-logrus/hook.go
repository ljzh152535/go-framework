package goadmin_logrus

import (
	//"github.com/orandin/lumberjackrus"

	"github.com/ljzh152535/go-framework/goadmin-logrus/lumberjackrus"
	"github.com/sirupsen/logrus"
)

// 初始化Logrushook
func InitLogrushook(logConf CoreLogrus) *logrus.Entry {
	logLocker.RLock()
	if log != nil {
		logLocker.RUnlock()
		return log
	}
	logLocker.RUnlock()

	// A,B,C
	logLocker.Lock()
	defer logLocker.Unlock()

	// 二次判断
	if log != nil {
		return log
	}

	logNew := logrus.New()

	// 显示代码行数
	logNew.SetReportCaller(logConf.IsSetReportCaller)

	logNew.SetLevel(transportLevel(logConf.Level))

	// 配置是否是控制台输出
	setLogrusOutputwithHook(logConf, logNew)

	// 添加字段
	l := addFields(logConf, logNew)
	log = l
	return l
}

func getLogName(logDir string, typeName string, logName string) string {
	return logDir + typeName + "-" + logName + ".log"
}

func newRotateHook(logConf CoreLogrus) logrus.Hook {
	logDir := logConf.LogDir + "/"
	getLevel := transportLevel(logConf.Level)

	hook, _ := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{
			// 通用日志配置
			Filename:   getLogName(logDir, "output", logConf.LogName),
			MaxSize:    logConf.MaxCapacity, // 单个文件备份的大小 以 M 为单位
			MaxBackups: logConf.MaxCount,    // 文件最多备份的数量
			MaxAge:     logConf.MaxAge,      // 保留文件的最大天数
			Compress:   true,                // 是否压缩
			LocalTime:  true,                // 是否启动本地时间
		},
		getLevel,
		setLogrusFormat(logConf),
		&lumberjackrus.LogFileOpts{
			// 针对不同日志级别的配置
			logrus.TraceLevel: &lumberjackrus.LogFile{
				Filename:   getLogName(logDir, "trace", logConf.LogName),
				MaxSize:    logConf.MaxCapacity,
				MaxBackups: logConf.MaxCount,
				MaxAge:     logConf.MaxAge,
				Compress:   true,
				LocalTime:  true,
			},
			logrus.DebugLevel: &lumberjackrus.LogFile{
				Filename:   getLogName(logDir, "debug", logConf.LogName),
				MaxSize:    logConf.MaxCapacity,
				MaxBackups: logConf.MaxCount,
				MaxAge:     logConf.MaxAge,
				Compress:   true,
				LocalTime:  true,
			},
			logrus.InfoLevel: &lumberjackrus.LogFile{
				Filename:   getLogName(logDir, "info", logConf.LogName),
				MaxSize:    logConf.MaxCapacity,
				MaxBackups: logConf.MaxCount,
				MaxAge:     logConf.MaxAge,
				Compress:   true,
				LocalTime:  true,
			},
			logrus.WarnLevel: &lumberjackrus.LogFile{
				Filename:   getLogName(logDir, "warn", logConf.LogName),
				MaxSize:    logConf.MaxCapacity,
				MaxBackups: logConf.MaxCount,
				MaxAge:     logConf.MaxAge,
				Compress:   true,
				LocalTime:  true,
			},
			logrus.ErrorLevel: &lumberjackrus.LogFile{
				Filename:   getLogName(logDir, "error", logConf.LogName),
				MaxSize:    logConf.MaxCapacity,
				MaxBackups: logConf.MaxCount,
				MaxAge:     logConf.MaxAge,
				Compress:   true,
				LocalTime:  true,
			},
			logrus.FatalLevel: &lumberjackrus.LogFile{
				Filename:   getLogName(logDir, "fatal", logConf.LogName),
				MaxSize:    logConf.MaxCapacity,
				MaxBackups: logConf.MaxCount,
				MaxAge:     logConf.MaxAge,
				Compress:   true,
				LocalTime:  true,
			},
			logrus.PanicLevel: &lumberjackrus.LogFile{
				Filename:   getLogName(logDir, "panic", logConf.LogName),
				MaxSize:    logConf.MaxCapacity,
				MaxBackups: logConf.MaxCount,
				MaxAge:     logConf.MaxAge,
				Compress:   true,
				LocalTime:  true,
			},
		},
	)
	return hook
}
