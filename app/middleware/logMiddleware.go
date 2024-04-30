package middlewares

//
//import (
//	"admin-backend/global"
//	"bytes"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"github.com/goccy/go-json"
//	"io"
//	"io/ioutil"
//	"net/http"
//	"strings"
//	"time"
//)
//
//// bodyLogWriter 定义一个存储响应内容的结构体
//type bodyLogWriter struct {
//	gin.ResponseWriter
//	body *bytes.Buffer
//}
//
//// Write 读取响应数据
//func (w bodyLogWriter) Write(b []byte) (int, error) {
//	w.body.Write(b)
//	return w.ResponseWriter.Write(b)
//}
//
//func getRequestBody(ctx *gin.Context) interface{} {
//	switch ctx.Request.Method {
//	case http.MethodGet:
//		return ctx.Request.URL.Query()
//	case http.MethodPost:
//		fallthrough
//	case http.MethodPut:
//		fallthrough
//	case http.MethodPatch:
//		var bodyBytes []byte
//		bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
//		if err != nil {
//			return nil
//		}
//		ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
//		return string(bodyBytes)
//	}
//	return nil
//}
//
//// RequestLog gin请求日志中间件
//func RequestLog(ctx *gin.Context) {
//	t := time.Now()
//
//	// 初始化bodyLogWriter
//	blw := &bodyLogWriter{
//		body:           bytes.NewBufferString(""),
//		ResponseWriter: ctx.Writer,
//	}
//	ctx.Writer = blw
//
//	// 获取请求信息
//	requestBody := getRequestBody(ctx)
//
//	ctx.Next()
//
//	// 记录响应信息
//	// 请求时间
//	costTime := time.Since(t)
//
//	// 响应内容
//	responseBody := blw.body.String()
//
//	// 日志格式
//	logContext := make(map[string]interface{})
//	logContext["request_uri"] = ctx.Request.RequestURI
//	logContext["request_method"] = ctx.Request.Method
//	logContext["refer_service_name"] = ctx.Request.Referer()
//	logContext["refer_request_host"] = ctx.ClientIP()
//	logContext["request_body"] = requestBody
//	logContext["request_time"] = t.String()
//	logContext["response_body"] = responseBody
//	logContext["time_used"] = fmt.Sprintf("%v", costTime)
//	logContext["header"] = ctx.Request.Header
//	global.GVA_LOG.WithFields(logContext).Info()
//	//global.GVA_LOG.Info(logContext)
//	//log.Println(logContext)
//}
//
//// LogLayout 日志layout
//type LogLayout struct {
//	Time      time.Time
//	Metadata  map[string]interface{} // 存储自定义原数据
//	Path      string                 // 访问路径
//	Query     string                 // 携带query
//	Body      string                 // 携带body数据
//	IP        string                 // ip地址
//	UserAgent string                 // 代理
//	Error     string                 // 错误
//	Cost      time.Duration          // 花费时间
//	Source    string                 // 来源
//}
//
//const ()
//
//type Logger struct {
//	// Filter 用户自定义过滤
//	Filter func(c *gin.Context) bool
//	// FilterKeyword 关键字过滤(key)
//	FilterKeyword func(layout *LogLayout) bool
//	// AuthProcess 鉴权处理
//	AuthProcess func(c *gin.Context, layout *LogLayout)
//	// 日志处理
//	Print func(LogLayout)
//	// Source 服务唯一标识
//	Source string
//}
//
//func (l Logger) SetLoggerMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		start := time.Now()
//		path := c.Request.URL.Path
//		query := c.Request.URL.RawQuery
//		var body []byte
//		if l.Filter != nil && !l.Filter(c) {
//			body, _ = c.GetRawData()
//			// 将原body塞回去
//			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
//		}
//		c.Next()
//		cost := time.Since(start)
//		layout := LogLayout{
//			Time:      time.Now(),
//			Path:      path,
//			Query:     query,
//			IP:        c.ClientIP(),
//			UserAgent: c.Request.UserAgent(),
//			Error:     strings.TrimRight(c.Errors.ByType(gin.ErrorTypePrivate).String(), "\n"),
//			Cost:      cost,
//			Source:    l.Source,
//		}
//		if l.Filter != nil && !l.Filter(c) {
//			layout.Body = string(body)
//		}
//		if l.AuthProcess != nil {
//			// 处理鉴权需要的信息
//			l.AuthProcess(c, &layout)
//		}
//		if l.FilterKeyword != nil {
//			// 自行判断key/value 脱敏等
//			l.FilterKeyword(&layout)
//		}
//		// 自行处理日志
//		l.Print(layout)
//	}
//}
//
//func DefaultLogger() gin.HandlerFunc {
//	return Logger{
//		Print: func(layout LogLayout) {
//			// 标准输出,k8s做收集
//			v, _ := json.Marshal(layout)
//			fmt.Println(string(v))
//		},
//		Source: "GVA",
//	}.SetLoggerMiddleware()
//}
