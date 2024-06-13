package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"server/global"
	sys_base2 "server/module/sys_base/model"
	"server/module/sys_base/service"
	"strings"
	"time"
	model3 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OperationRecord 操作记录
func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		if c.Request.URL.Path == "/socket.io/" {
			c.Next()
			return
		}
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				global.GLog.Error("read body from request error:", zap.Any("err", err))
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}
		userId := c.Request.Header.Get("x-user-id")
		record := sys_base2.SysOperationRecord{
			Ip:     c.Request.Header.Get("X-Real-IP"),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
			UserID: userId,
		}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		now := time.Now()

		c.Next()

		latency := time.Now().Sub(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency
		record.Resp = writer.body.String()
		record.ID = uint(time.Now().Unix())

		if err := service.CreateSysOperationRecord(record); err != nil {
			global.GLog.Error("create operation record error:", zap.Any("err", err))
		}
	}
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/socket.io/" {
			c.Next()
			return
		}
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		//fmt.Println(c.Request.Header.Get("X-Real-IP"))
		global.GLog.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("x-code", c.GetHeader(global.XCode)),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.Request.Header.Get("X-Real-IP")),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

type ProCode struct {
	ProCode string `json:"proCode" form:"proCode"`
}

// SysLogger 用户查看的日志
func SysLogger(logType int) gin.HandlerFunc {
	return func(c *gin.Context) {
		//if c.Request.Method == "GET" {
		//	c.Next()
		//	return
		//}
		uri := c.Request.RequestURI //uri
		if noLogger(uri) || global.GConfig.System.Env == "develop" {
			c.Next()
			return
		}
		claims, isExit := c.Get("claims")
		if !isExit {
			global.GLog.Error("未获取到登录人信息")
			c.Next()
			return
		}
		waitUse := claims.(*model3.CustomClaims)
		proCode := c.Request.Header.Get("pro-code")
		if proCode == "" {
			proCode = c.Query("proCode")
		}
		//if proCode == "" {
		//	var proCodeStruct ProCode
		//	err := c.ShouldBindJSON(&proCodeStruct)
		//	if err != nil {
		//		global.GLog.Error("获取参数:::" + err.Error())
		//	}
		//	proCode = proCodeStruct.ProCode
		//}
		title, ok := global.ApiPathTitle[c.Request.URL.Path]
		if strings.Index(title, "--") == -1 {
			global.GLog.Error(fmt.Sprintf("插入日志::: API标题%s有误,%s", title, c.Request.URL.Path))
			c.Next()
			return
		}
		if ok {
			err := service.AddSysLogger(sys_base2.SysLogger{
				ProCode:         proCode,
				FunctionModule:  strings.Split(title, "--")[0],
				ModuleOperation: strings.Split(title, "--")[1],
				OperationCode:   waitUse.Code,
				OperationName:   waitUse.Name,
				Content:         "",
				Api:             c.Request.URL.Path,
				LogType:         logType,
			})
			if err != nil {
				global.GLog.Error("插入日志:::" + err.Error())
			}
		}
		c.Next()
	}
}

// noLogger 不需要logger接口
func noLogger(str string) bool {
	List := []string{
		"/task/opNum",
		"/task/conf",
		"/task/releaseBill",
		"/task/op",
		"/task/submit",
		"/task/get-image-by-block-id",
		"/task/get-image-by-index",
		"/task/get-thumbnail",
		"/sys-socket-io-notice/customer-noticel",
	}
	for _, s := range List {
		if match, _ := regexp.MatchString(s, str); match {
			return true
		}
	}
	return false
}
