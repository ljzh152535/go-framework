package goadmin_logrus

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type LogFormatter struct{}

func (m *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
	var msg string
	//entry.Logger.SetReportCaller(true)
	//HasCaller()为true才会有调用信息
	if entry.HasCaller() {
		fName := filepath.Base(entry.Caller.File)
		//msg = fmt.Sprintf("[%s] [%s] [%s:%d %s] [%s] %s\n",
		//timestamp, entry.Level.String(), fName, entry.Caller.Line, entry.Caller.Function, entry.Message, entry.Message)
		msg = fmt.Sprintf("[%s] [%s] [%s:%d %s] %s\n",
			timestamp, entry.Level.String(), fName, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		msg = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)
	}

	b.WriteString(msg)
	return b.Bytes(), nil
}

// 设置日志输出格式json text
func setLogrusFormat(logConf CoreLogrus) logrus.Formatter {

	if logConf.LogFormat == "json" {
		return &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05.000",
			// runtime.Frame: 帧,可用于获取调用者返回的PC值的函数、文件或者是行信息
			CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
				fileName := path.Base(frame.File)
				return frame.Function, fmt.Sprintf("%s:%d", fileName, frame.Line)
			},
		}
	} else {
		return &LogFormatter{}
	}
}

// 设置日志输出方式
func setLogrusOutputwithHook(logConf CoreLogrus, logNew *logrus.Logger) {
	if logConf.Output == "file" {
		// 配置logrus的hook
		logNew.AddHook(newRotateHook(logConf))
		// 禁用控制台输出
		if logConf.Output != "output" {
			src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				fmt.Println("Open Src File err", err)
			}
			writer := bufio.NewWriter(src)
			logNew.SetOutput(writer)
		}
	} else if logConf.Output == "output" {
		logNew.SetOutput(os.Stdout)
		// 设置日志格式
		logNew.SetFormatter(setLogrusFormat(logConf))
	} else {
		panic("日志格式配置错误,请检查配置,值为:" + logConf.Output)
	}
}

// 设置日志输出方式
func setLogrusOutput(logConf CoreLogrus, logNew *logrus.Logger) {
	if logConf.Output == "file" {
		f, e := loadLogFile(logConf)
		if e != nil {
			panic(e)
		}
		logNew.SetOutput(f)
	} else {
		logNew.SetOutput(os.Stdout)
	}
	if logConf.Output != "file" {
		//&logrus.Logger{}.SetOutput(os.Stdout)
	}
}

// 自定义添加字段
func addFields(logConf CoreLogrus, logNew *logrus.Logger) *logrus.Entry {

	name, err := os.Hostname()
	if err != nil {
		panic("无法获取到主机名:" + err.Error())
	}

	ips, err := GetLocalIPs()
	if err != nil {
		log.Fatal(err)
	}

	l := logNew.WithFields(logrus.Fields{
		"loccal_ips": ips,
		"hostname":   name,
		"traceID":    logConf.TraceID,
		"appName":    logConf.AppName,
	})

	return l
}

// 获取ip

func GetLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addresses {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP)
			}
		}
	}
	return ips, nil
}

// TransportLevel 根据字符串转化为 zapcore.Level
func transportLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "warn":
		return logrus.WarnLevel
	case "error":
		return logrus.WarnLevel
	case "panic":
		return logrus.PanicLevel
	case "fatal":
		return logrus.FatalLevel
	default:
		return logrus.DebugLevel
	}
}
