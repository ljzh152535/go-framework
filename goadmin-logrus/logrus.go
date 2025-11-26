package goadmin_logrus

import (
	"bufio"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"time"
)

func InitLogrus(logConf CoreLogrus) *logrus.Entry {
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

	// 设置log的配置

	logNew.AddHook(newRotateHook(logConf))

	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Open Src File err", err)
	}
	writer := bufio.NewWriter(src)
	logNew.SetOutput(writer)

	name, err := os.Hostname()
	if err != nil {
		panic("无法获取到主机名:" + err.Error())
	}
	l := logNew.WithFields(logrus.Fields{
		//"env": logConf.LogEnv,
		//"loccal_ip": env.LocalIP(),
		"hostname": name,
	})
	log = l
	return l
}

// 设置日志level
func setLogrusLevel(logConf CoreLogrus, logNew *logrus.Logger) {
	switch logConf.Level {
	case "debug":
		logNew.SetLevel(logrus.DebugLevel)
	case "error":
		logNew.SetLevel(logrus.ErrorLevel)
	case "warn":
		logNew.SetLevel(logrus.WarnLevel)
	default:
		logNew.SetLevel(logrus.InfoLevel)
	}
}

func loadLogFile(logConf CoreLogrus) (io.Writer, error) {
	logPath := "logs/app.log"
	if logConf.LogDir != "" {
		logPath = logConf.LogDir + string(os.PathSeparator) + logConf.LogName + ".log"
	}

	// 判断logPath是相对路径还是绝对路径
	//if !filepath.IsAbs(logPath) {
	//	logPath = homeDir + "/" + logPath
	//}

	// 检查文件是否存在，不存在创建文件
	f, e := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
	if e != nil {
		return nil, e
	}
	return f, nil
}

func setLogrusConf(logConf CoreLogrus, logNew *logrus.Logger) *logrus.Entry {

	// 设置日志level
	setLogrusLevel(logConf, logNew)

	// 显示文件和代码行数
	logNew.SetReportCaller(logConf.IsSetReportCaller) // 显示文件和代码行数

	// 设置日志输出方式
	setLogrusOutput(logConf, logNew)

	// 设置日志格式 json 或者 text
	//setLogrusFormat(logConf, logNew)

	// 基础字段预设,比如项目名、环境、env、local_ip、hostname、idc
	name, err := os.Hostname()
	if err != nil {
		panic("无法获取到主机名:" + err.Error())
	}
	l := logNew.WithFields(logrus.Fields{
		//"env": logConf.LogEnv,
		//"loccal_ip": env.LocalIP(),
		"hostname": name,
	})
	log = l
	return l
}

// 设置日志切割
func setRotatelogs(logConf CoreLogrus) *rotatelogs.RotateLogs {

	writer, err := rotatelogs.New(
		logConf.LogDir+string(os.PathSeparator)+logConf.LogName+"-%Y%m%d.log",
		//rotatelogs.WithLinkName(logConf.LogDir+string(os.PathSeparator)+logConf.LogName+".log"), //生成软链，指向最新日志文件
		rotatelogs.WithLinkName(logConf.LogDir+string(os.PathSeparator)+logConf.LogName+".log"), //生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(time.Hour*24*time.Duration(logConf.MaxAge)),                       // 最大的保留时间 单位: 天
		rotatelogs.WithRotationTime(24*time.Hour),                                               //最小为1分钟轮询。默认60s  低于1分钟就按1分钟来
		rotatelogs.WithRotationSize(int64(logConf.MaxCapacity)*1024*1024),                       // 设置分割文件的大小为  单位: MB
		rotatelogs.WithRotationCount(uint(logConf.MaxCount)),
	)

	if err != nil {
		log.Error(err)
	}
	return writer
}
