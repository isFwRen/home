package router

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http/httputil"
	"net/url"
	"reflect"
	"regexp"
	"server/global"
	"server/middleware"
	"server/module/pro_conf/api"
	"server/module/pro_conf/model"
	"server/module/pro_conf/service"
	"strings"
	api3 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

var listReg = regexp.MustCompile(`^/sys-const/info-list/B`)
var insertReg = regexp.MustCompile(`^/sys-const/insert$`)
var editReg = regexp.MustCompile(`^/sys-const/edit$`)
var delTablesReg = regexp.MustCompile(`^/sys-const/del-tables$`)
var delLinesReg = regexp.MustCompile(`^/sys-const/del-lines$`)
var importReg = regexp.MustCompile(`^/sys-const/import/B`)
var exportReg = regexp.MustCompile(`^/sys-const/export$`)
var releaseReg = regexp.MustCompile(`^/sys-const/release/B`)

func InitProjectConstManagement(Router *gin.RouterGroup) {
	projectCostManagementRouter := Router.Group("pro-config/").
		//Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(1))
	{
		//根据项目名称获取该项目下的常量表
		projectCostManagementRouter.GET("project-cost-management/list-const-name", api.GetConstTablesByProjectName)
		//直接保存传进来的数据
		projectCostManagementRouter.POST("project-cost-management/putConstTableByArr", api.PutConstTableByArr)
		//解析excel保存数据
		projectCostManagementRouter.POST("project-cost-management/putConstTableByExcel", api.PutConstTableWithExcel)
		//删除当前行
		projectCostManagementRouter.POST("project-cost-management/delete-line", api.DelConstTableLineById)
		//删除当前常量表
		projectCostManagementRouter.POST("project-cost-management/delete-table", api.DelConstTableById)
		//获取表头
		projectCostManagementRouter.GET("project-cost-management/list-table-top", api.GetTableTop)
		//获取导出常量表
		projectCostManagementRouter.GET("project-cost-management/export-const", api.ExportConstExcel)
		//获取常量操作日志
		projectCostManagementRouter.POST("const/operation-log/page", api.PageOperationLog)
	}

	constManagementRouter := Router.Group("").
		//Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(1)).
		//Use(middleware.Cors()).
		Use(constOperationRecord())
	{
		// 定义代理目标URL
		targetURL, _ := url.Parse(global.GConfig.System.ConstUrl)
		// 创建反向代理
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		// 定义代理路由
		constManagementRouter.Any("sys-const/*path", func(c *gin.Context) {
			global.GLog.Info("Proxying request to:", zap.Any("s", targetURL))
			c.Request.Header.Set("Authorization", "Bearer 5pyq6K6+572uYmFzZTY05qC85byP55qE546v5aKD5Y+Y6YePRFVSSUFOX0FQSV9UT0tFTuWtmOWcqOS4pemHjeeahOWuieWFqOmXrumimA==")
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	}
}

// constOperationRecord 常量操作记录
func constOperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		body := ""
		//f, _ := c.MultipartForm()
		if strings.Index(c.GetHeader("Content-Type"), "form-data; boundary") == -1 {
			bodyByte, err := io.ReadAll(c.Request.Body)
			body = string(bodyByte)
			if err != nil {
				global.GLog.Error("read body from request error:", zap.Any("err", err))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
			}
		}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		c.Next()

		bodyJson, _ := simplejson.NewJson([]byte(body))
		bodyRespJson, _ := simplejson.NewFromReader(writer.body)
		u, err := api3.GetUserByToken(c)
		if err != nil {
			return
		}
		constLog := model.ConstLog{UserId: c.GetHeader("x-user-id"), UserName: u.NickName, Status: c.Writer.Status()}
		switch {
		case listReg.MatchString(path) && c.Request.Method == "GET":
			constLog.Type = "查看"
			constLog.Content = "查看常量配置"
		case insertReg.MatchString(path) && c.Request.Method == "POST":
			//{"proCode":"B0106","name":"B0106_陕西国寿理赔_ICD10疾病编码","items":[{"疾病名称":"1","疾病代码":"2"},{"疾病名称":"3","疾病代码":"4"}]}
			constLog.Type = "编辑"
			items, _ := bodyJson.Get("items").Array()
			name, _ := bodyJson.Get("name").String()
			constLog.Content = fmt.Sprintf("将《%v》中", name)
			for _, item := range items {
				v := reflect.ValueOf(item)
				for _, key := range v.MapKeys() {
					value := v.MapIndex(key).Interface()
					constLog.Content += fmt.Sprintf("%v新增%v\n;", key.Interface(), value)
				}
			}
		case editReg.MatchString(path):
			constLog.Type = "编辑"
			//{"proCode":"B0106","name":"B0106_陕西国寿理赔_ICD10疾病编码","item":{"疾病名称":"绿脓杆菌感染","疾病代码":"A49.8081"},"itemBefore":{"疾病名称":"绿脓杆菌感染11","疾病代码":"A49.808111"},"id":"00F6xrgmGmBdMPnm"}
			name, _ := bodyJson.Get("name").String()
			item := reflect.ValueOf(bodyJson.Get("item").Interface())
			itemBefore := reflect.ValueOf(bodyJson.Get("itemBefore").Interface())
			msg := ""
			msgBefore := ""
			for _, key := range item.MapKeys() {
				value := item.MapIndex(key).Interface()
				valueBefore := itemBefore.MapIndex(key).Interface()
				msg += fmt.Sprintf("【%v】%v", key.Interface(), value)
				msgBefore += fmt.Sprintf("【%v】%v", key.Interface(), valueBefore)
			}
			constLog.Content = fmt.Sprintf("将《%v》中%v修改为%v", name, msgBefore, msg)
		case delTablesReg.MatchString(path) && c.Request.Method == "POST":
			//{"proCode":"B0106","name":["B0106_陕西国寿理赔_ICD10疾病编码"]}
			constLog.Type = "删除"
			names, _ := bodyJson.Get("name").StringArray()
			constLog.Content = fmt.Sprintf("删除《%v》常量表", strings.Join(names, "、"))
		case delLinesReg.MatchString(path) && c.Request.Method == "POST":
			//{"proCode":"B0106","ids":["00F6xrgmGmBdMPnm"],"items":[{"疾病名称":"绿脓杆菌感染","疾病代码":"A49.8081"}],"name":"B0106_陕西国寿理赔_ICD10疾病编码"}
			constLog.Type = "删除"
			name, _ := bodyJson.Get("name").String()
			items, _ := bodyJson.Get("items").Array()
			msg := ""
			for _, item := range items {
				v := reflect.ValueOf(item)
				for _, key := range v.MapKeys() {
					value := v.MapIndex(key).Interface()
					msg += fmt.Sprintf("【%v】%v,", key.Interface(), value)
				}
				msg += ";\n"
			}
			constLog.Content = fmt.Sprintf("将《%v》中\n%v进行删除", name, msg)
		case importReg.MatchString(path) && c.Request.Method == "POST":
			constLog.Type = "上传"
			var names []string
			tables, err := bodyRespJson.Get("tables").Array()
			if err != nil {
				global.GLog.Error("上传 bodyRespJson tables", zap.Error(err))
				return
			}
			for _, table := range tables {
				names = append(names, table.(map[string]interface{})["name"].(string))
			}
			constLog.Content = fmt.Sprintf("上传《%v》常量表", strings.Join(names, "、"))
		case exportReg.MatchString(path) && c.Request.Method == "POST":
			constLog.Type = "导出"
			name, _ := bodyJson.Get("name").String()
			constLog.Content = fmt.Sprintf("导出《%v》常量表", name)
		case releaseReg.MatchString(path) && c.Request.Method == "PATCH":
			//{"name":["B0106_陕西国寿理赔_医疗机构61"],"proCode":"B0106"}
			constLog.Type = "发布"
			names, _ := bodyJson.Get("name").StringArray()
			constLog.Content = fmt.Sprintf("发布《%v》常量表", strings.Join(names, "、"))
		default:
			global.GLog.Error(path, zap.Error(errors.New("该请求不做记录")))
			return
		}

		err = service.InsertConstLog(constLog)
		if err != nil {
			global.GLog.Error(path, zap.Error(errors.New("记录常量操作失败")))
			return
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
