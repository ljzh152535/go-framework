// Copyright 2014 Manu Martinez-Almeida. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	goadmin_config "github.com/ljzh152535/go-framework/goadmin-config"
	"github.com/ljzh152535/go-framework/goadmin-utils/timeutils"
	"github.com/rs/xid"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mattn/go-isatty"
)

type consoleColorModeValue int

const (
	autoColor consoleColorModeValue = iota
	disableColor
	forceColor
)

const logTimeTpl = "2006-01-02T15:04:05.000Z07:00"
const LogID = "X-Log-Id"

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

var consoleColorMode = autoColor

// LoggerConfig defines the config for Logger middleware.
type LoggerConfig struct {
	// Optional. Default value is gin.defaultLogFormatter
	Formatter LogFormatter

	// Output is a writer where logs are written.
	// Optional. Default value is gin.DefaultWriter.
	Output io.Writer

	// SkipPaths is an url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

// LogFormatter gives the signature of the formatter function passed to LoggerWithFormatter
type LogFormatter func(params LogFormatterParams) string

// LogFormatterParams is the structure any formatter will be handed when time to log comes
type LogFormatterParams struct {
	LocalIp  string
	Env      string
	Hostname string
	Format   string // 文件输出格式
	LogId    string

	Request     *http.Request
	RequestData interface{}

	StartTime string
	EndTime   string
	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[any]any
	
	ResponseData interface{}
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *LogFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *LogFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *LogFormatterParams) ResetColor() string {
	return reset
}

// IsOutputColor indicates whether can colors be outputted to the log.
func (p *LogFormatterParams) IsOutputColor() bool {
	return consoleColorMode == forceColor || (consoleColorMode == autoColor && p.isTerm)
}

// defaultLogFormatter is the default log format function Logger middleware uses.
var defaultLogFormatter = func(param LogFormatterParams) string {
	var logId = param.Request.Header.Get("X-Log-Id")
	if logId == "" {
		logId = xid.New().String()
	}

	// 支持text , json
	if param.Format == "json" {
		var rel = map[string]any{
			"startTime":   param.StartTime,
			"endTime":     param.EndTime,
			"ip":          param.ClientIP,
			"log_id":      param.LogId,
			"time":        param.TimeStamp.Format(logTimeTpl),
			"latency":     param.Latency.Milliseconds(),
			"method":      param.Method,
			"path":        param.Path,
			"query":       param.Request.URL.RawQuery,
			"status":      param.StatusCode,
			"error":       param.ErrorMessage,
			"size":        param.BodySize,
			"local_ip":    param.LocalIp,
			"env":         param.Env,
			"hostname":    param.Hostname,
			"user_agent":  param.Request.UserAgent(),
			"referer":     param.Request.Referer(),
			"responeData": param.ResponseData,
			"requestData": param.RequestData,
		}
		relJson, _ := json.Marshal(rel)
		return string(relJson) + "\n"

	}

	// text
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	//if param.Latency > time.Minute {
	//	param.Latency = param.Latency.Truncate(time.Second)
	//}
	return fmt.Sprintf("[%s]\t%s\t%s\t%d\t%s%s%s\t%s\t%s\t%s%d%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
		param.ClientIP,
		param.LogId,
		param.TimeStamp.Format(logTimeTpl),
		param.Latency.Milliseconds(),
		methodColor, param.Method, resetColor,
		param.Path,
		param.Request.URL.RawQuery,
		statusColor, param.StatusCode, resetColor,
		param.ErrorMessage,
		param.BodySize,
		param.Env,
		param.Hostname,
		param.Request.UserAgent(),
		param.Request.Referer(),
	)
}

// DisableConsoleColor disables color output in the console.
func DisableConsoleColor() {
	consoleColorMode = disableColor
}

// ForceConsoleColor force color output in the console.
func ForceConsoleColor() {
	consoleColorMode = forceColor
}

// ErrorLogger returns a HandlerFunc for any error type.
func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

// ErrorLoggerT returns a HandlerFunc for a given error type.
func ErrorLoggerT(typ gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		errors := c.Errors.ByType(typ)
		if len(errors) > 0 {
			c.JSON(-1, errors)
		}
	}
}

// Logger instances a Logger middleware that will write the logs to gin.DefaultWriter.
// By default, gin.DefaultWriter = os.Stdout.
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{})
}

// LoggerWithFormatter instance a Logger middleware with the specified log format function.
func LoggerWithFormatter(f LogFormatter) gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Formatter: f,
	})
}

// LoggerWithWriter instance a Logger middleware with the specified writer buffer.
// Example: os.Stdout, a file opened in write mode, a socket...
func LoggerWithWriter(out io.Writer, notlogged ...string) gin.HandlerFunc {
	return LoggerWithConfig(LoggerConfig{
		Output:    out,
		SkipPaths: notlogged,
	})
}

var output io.Writer = os.Stdout

type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

func getRequestBody(ctx *gin.Context) interface{} {
	switch ctx.Request.Method {
	case http.MethodGet:
		return ctx.Request.URL.Query()
	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		var bodyBytes []byte
		bodyBytes, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			return nil
		}
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return string(bodyBytes)
	}
	return nil
}

func New(conf goadmin_config.WebServerLog, prjEnv string, homeDir string) gin.HandlerFunc {
	switch conf.Output {
	case "file":
		var err error
		output, err = loadLogFile(conf, homeDir)
		if err != nil {
			panic(err)
		}
	default:
		output = os.Stdout
	}

	localIP := conf.LocalIP
	hostname := conf.HostName

	formatter := defaultLogFormatter

	return func(c *gin.Context) {
		var logId = c.GetHeader(LogID)
		if logId == "" {
			logId = xid.New().String()
			//c.Request.Header.Add(LogID, logId)
			//c.Set("claims", "Set")
			//c.Header("claims", "Header")
			c.Header(LogID, logId)
			//c.Request.Header.Set(LogID, logId)
			//c.Set(LogID, logId)
		}
		c.Set(LogID, logId) // 设置上下文
		strBody := ""
		var blw bodyLogWriter
		blw = bodyLogWriter{bodyBuf: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		startTime := time.Now()
		path := c.Request.URL.Path

		c.Next()

		strBody = strings.Trim(blw.bodyBuf.String(), "\n")
		var skip map[string]struct{}
		if length := len(conf.SkipPaths); length > 0 {
			skip = make(map[string]struct{}, length)
			for _, p := range conf.SkipPaths {
				skip[p] = struct{}{}

			}
		}

		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request: c.Request,
				Keys:    c.Keys,
			}

			endTime := time.Now()
			// 获取请求信息
			param.Latency = endTime.Sub(startTime)
			requestBody := getRequestBody(c)

			param.RequestData = requestBody
			// stop timer
			param.Format = conf.LogFormat
			param.LocalIp = localIP
			param.Hostname = hostname
			param.TimeStamp = time.Now()
			param.LogId = logId
			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
			param.BodySize = c.Writer.Size()
			param.Env = prjEnv
			param.Path = path
			param.ResponseData = strBody
			param.StartTime = timeutils.GetCurrentTime(startTime)
			param.EndTime = timeutils.GetCurrentTime(endTime)
			_, _ = fmt.Fprint(output, formatter(param))
		}

	}
}

func loadLogFile(conf goadmin_config.WebServerLog, homeDir string) (io.Writer, error) {
	logPath := "logs/access.log"
	if conf.LogPath != "" {
		logPath = conf.LogPath
	}

	// 判断logPath是相对路径还是绝对路径
	//if !filepath.IsAbs(logPath) {
	//	logPath = homeDir + "/" + logPath
	//}

	// 检查文件是否存在，不存在创建文件
	f, e := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if e != nil {
		return nil, e
	}
	return f, nil
}

// LoggerWithConfig instance a Logger middleware with config.
func LoggerWithConfig(conf LoggerConfig) gin.HandlerFunc {
	formatter := conf.Formatter
	if formatter == nil {
		formatter = defaultLogFormatter
	}

	out := conf.Output
	if out == nil {
		out = gin.DefaultWriter
	}

	notlogged := conf.SkipPaths

	isTerm := true

	if w, ok := out.(*os.File); !ok || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(w.Fd()) && !isatty.IsCygwinTerminal(w.Fd())) {
		isTerm = false
	}

	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := LogFormatterParams{
				Request: c.Request,
				isTerm:  isTerm,
				Keys:    c.Keys,
			}

			// Stop timer
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)

			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

			param.BodySize = c.Writer.Size()

			if raw != "" {
				path = path + "?" + raw
			}

			param.Path = path

			fmt.Fprint(out, formatter(param))
		}
	}
}
