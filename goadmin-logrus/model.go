package goadmin_logrus

import (
	"github.com/sirupsen/logrus"
	"sync"
)

// log实例
const logTimeTpl = "2006-01-02T15:04:05.000Z07:00"

var log *logrus.Entry
var logLocker sync.RWMutex

type CoreLogrus struct {
	Level             string `mapstructure:"level" json:"level" yaml:"level"`                                     // 级别
	LogFormat         string `mapstructure:"logFormat" json:"logFormat" yaml:"logFormat"`                         // 日志格式 json text
	IsSetReportCaller bool   `mapstructure:"isSetReportCaller" json:"isSetReportCaller" yaml:"isSetReportCaller"` // 显示文件和代码行数
	Output            string `mapstructure:"output" json:"output" yaml:"output"`                                  // 日志输出方式  file output
	LogDir            string `mapstructure:"logDir" json:"logDir" yaml:"logDir"`                                  // 日志目录
	LogName           string `mapstructure:"logName" json:"logName" yaml:"logName"`                               // 日志名
	MaxAge            int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`                               // 日志留存时间 单位:天
	MaxCapacity       int    `mapstructure:"max-capacity" json:"max-capacity" yaml:"max-capacity"`                // 单个日志的最大容量 单位:M
	MaxCount          int    `mapstructure:"max-count" json:"max-count" yaml:"max-count"`                         // 保存日志的最大数量
	TraceID           string `mapstructure:"traceID" json:"traceID" yaml:"traceID"`                               // 日志ID
	AppName           string `mapstructure:"appName" json:"appName" yaml:"appName"`                               // 应用名
}
