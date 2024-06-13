package initialize

import (
	"bytes"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	routerDataEntry "server/module/data_entry_management/router"
	routerDingDing "server/module/dingding/router"
	router3 "server/module/homepage/router"
	router2 "server/module/msg_manager/router"
	"server/module/task/router"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	xingqiyiMiddleware "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/middleware"
	xingqiyiRouter "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/router"

	routerAssessmentManagement "server/module/assessment_management/router"
	routerPtTrainStage "server/module/training_guide/router"

	_ "server/docs"

	routerProConf "server/module/pro_conf/router"
	routerProManager "server/module/pro_manager/router"
	routerReport "server/module/report_management/router"
	"server/module/sys_base/model"
	routerSysBase "server/module/sys_base/router"
	"strings"
)

// Routers 初始化总路由
func Routers() (engine *gin.Engine) {
	Router := gin.New()

	// get global Monitor object
	m := ginmetrics.GetMonitor()
	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(Router)

	//跨域
	//global.GLog.Debug("use middleware cors")
	//Router.Use(middleware.Cors())
	//global.GLog.Debug("use middleware logger")
	//Router.Use(middleware.GinLogger())
	//Router.Use(middleware.JWTAuth())
	//Router.Use(middleware.GinRecovery(false))

	global.GLog.Debug("use middleware Cors")
	Router.Use(xingqiyiMiddleware.Cors())
	global.GLog.Debug("use middleware 404")
	Router.NoRoute(xingqiyiMiddleware.Gin404())
	global.GLog.Debug("use middleware logger")
	Router.Use(xingqiyiMiddleware.GinLogger())
	global.GLog.Debug("use middleware JWTAuth")
	Router.Use(xingqiyiMiddleware.JWTAuth())
	global.GLog.Debug("use middleware CasbinHandler")
	Router.Use(xingqiyiMiddleware.CasbinHandler())
	global.GLog.Debug("use middleware GinRecovery")
	Router.Use(xingqiyiMiddleware.GinRecovery(false))
	//global.GLog.Debug("use middleware Limiter")
	//limiter := tollbooth.NewLimiter(2, nil)
	//Router.Use(xingqiyiMiddleware.Limiter(limiter))

	//重定向到ssl地址
	//Router.Use(middleware.LoadTls())
	//为用户头像和文件提供静态地址
	Router.Static(global.GConfig.LocalUpload.FilePath, global.GConfig.LocalUpload.FilePath)
	//swagger
	if global.GConfig.System.Env != "prod" {
		Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		global.GLog.Debug("register swagger handler")
	}
	//方便统一添加路由组前缀 多服务器上线使用
	ApiGroup := Router.Group("")
	//登录、注册等基础接口
	xingqiyiRouter.InitBaseRouter(ApiGroup)
	//注册用户路由
	routerSysBase.InitUserRouter(ApiGroup)
	//注册基础功能路由 不做鉴权
	routerSysBase.InitBaseRouter(ApiGroup)
	//注册menu路由
	//router.InitMenuRouter(ApiGroup)
	//注册角色路由
	//routerSysBase.InitAuthorityRouter(ApiGroup)
	//权限相关路由
	//routerSysBase.InitCasbinRouter(ApiGroup)
	//jwt相关路由
	//router.InitJwtRouter(ApiGroup)
	//system相关路由
	//routerSysBase.InitSystemRouter(ApiGroup)
	//操作记录
	routerSysBase.InitSysOperationRecordRouter(ApiGroup)
	//项目配置
	routerProConf.InitSysProject(ApiGroup)
	//文件上传下载功能路由
	routerSysBase.InitFileUploadAndDownloadRouter(ApiGroup)
	//PT管理
	//router.InitPTManager(ApiGroup)
	//Tabs导航栏
	//routerSysBase.InitTabs(ApiGroup)
	//任务查询
	//routerProConf.InitTaskQuery(ApiGroup)
	//钉钉群管理
	//routerDingDing.InitDingdingGroupRouter(ApiGroup)
	//钉钉群消息
	//routerDingDing.InitDingdingGroupMsgsRouter(ApiGroup)
	//钉钉验证码
	routerDingDing.InitDingingBaseRouter(ApiGroup)
	//常量管理
	routerProConf.InitProjectConstManagement(ApiGroup)
	//路径配置管理
	routerProConf.SysProDownloadPaths(ApiGroup)
	//时效配置管理
	routerProConf.InitProjectConfigAging(ApiGroup)
	//节假日时效配置管理
	routerProConf.InitProjectConfigAgingHoliday(ApiGroup)
	//合同时效管理
	routerProConf.InitProjectConfigAgingContract(ApiGroup)
	//发送实时录入通知
	//router.InitEnteringMsgsRouter(ApiGroup)
	//报表管理
	routerReport.InitOutputStatistics(ApiGroup)
	//案件列表
	routerProManager.InitBillList(ApiGroup)
	//错误查询
	routerReport.InitErrorStatistics(ApiGroup)
	//来量分析
	routerReport.InitBusinessAnalysis(ApiGroup)
	//工资数据
	routerReport.InitSalary(ApiGroup)
	//时效管理
	routerProManager.InitProjectAgingManagement(ApiGroup)
	//异常件数据
	routerReport.InitAbnormalBill(ApiGroup)
	//业务明细表
	routerReport.InitProjectReport(ApiGroup)
	//项目报表配置表
	routerReport.InitSettingReportTag(ApiGroup)
	//角色管理
	routerSysBase.InitRoleManagement(ApiGroup)
	//用户管理
	routerSysBase.InitUserManagement(ApiGroup)
	//菜单管理
	routerSysBase.InitMenuRouter(ApiGroup)
	//项目权限统计
	routerSysBase.InitProPermissionSum(ApiGroup)
	//白名单
	routerSysBase.InitWhiteListRouter(ApiGroup)
	//案件明细
	routerProManager.InitTaskQuery(ApiGroup)
	//系统日志
	routerSysBase.InitSysLoggerRouter(ApiGroup)
	//录入管理
	routerDataEntry.InitDataEntryManagement(ApiGroup)
	//录入的保存api和日志要用而已
	router.InitTask(ApiGroup)
	//查询员工
	routerSysBase.InitUManagement(ApiGroup)
	//质量管理
	routerProManager.InitQualityManagement(ApiGroup)
	routerProManager.InitQualityAnalysis(ApiGroup)
	//客户投诉
	routerProManager.InitTaskCustomerComplaints(ApiGroup)
	//业务规则
	routerProManager.InitTransactionRuleRouter(ApiGroup)
	//报销单模板
	routerProManager.InitReimbursementFormTemplate(ApiGroup)
	//字段规则
	routerProManager.InitFieldsRule(ApiGroup)
	//教学视频
	routerProManager.InitTeachVideo(ApiGroup)
	//二期
	//钉钉群管理
	router2.InitDingtalk(ApiGroup)
	//钉钉群通知
	router2.InitDingtalkNotice(ApiGroup)
	//公告管理
	routerProManager.InitAnnouncementRouter(ApiGroup)
	//主页管理
	router3.InitHomeRouter(ApiGroup)
	//固定通知
	router2.InitGroupNotice(ApiGroup)
	//录入通知
	router2.InitTaskNotice(ApiGroup)
	//项目数据
	router3.InitProDataRouter(ApiGroup)
	//项目日报
	router3.InitProReportRouter(ApiGroup)

	routerProManager.InitSysSpotCheckRouter(ApiGroup)

	//客户通知
	router2.InitCustomerNotice(ApiGroup)
	routerSysBase.InitSocketIoNotice(ApiGroup)
	routerReport.InitSpecialReport(ApiGroup)

	//业务通知
	router2.InitBusinessPush(ApiGroup)

	//PT培训流程指引 -----*测试 正式需转移到 router_task.go
	routerPtTrainStage.InitTrainingStageRouter(ApiGroup)
	routerPtTrainStage.InitRuleSchedule(ApiGroup)
	//考核管理
	routerAssessmentManagement.InitAssessLoggingRouter(ApiGroup)

	global.GLog.Info("router register success")

	//gp := ginprometheus.New(Router)
	//Router.Use(gp.Middleware())
	//// metrics采样
	//Router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	//存储所有api
	//service.SaveApi(Router)
	return Router
}

// 设置单文件访问,不能访问目录
func limitDir(w http.ResponseWriter, r *http.Request) {
	global.GLog.Error("不能访问目录:%s\n" + r.URL.Path)
	old := r.URL.Path
	name := path.Clean("D:/code/20160902/src" + strings.Replace(old, "/js", "/client", 1))
	info, err := os.Lstat(name)
	if err == nil {
		if !info.IsDir() {
			http.ServeFile(w, r, name)
		} else {
			http.NotFound(w, r)
		}
	} else {
		http.NotFound(w, r)
	}
}

func saveApi(Router *gin.Engine) {
	for _, v := range Router.Routes() {
		if !(v.Method == "HEAD" ||
			strings.Contains(v.Handler, "github.com") ||
			strings.Contains(v.Path, "/swagger/") ||
			strings.Contains(v.Path, "/static/") ||
			strings.Contains(v.Path, "/form-generator/") ||
			strings.Contains(v.Path, "/sys/tables")) {

			// 根据接口方法注释里的@Summary填充接口名称，适用于代码生成器
			// 可在此处增加配置路径前缀的if判断，只对代码生成的自建应用进行定向的接口名称填充
			jsonFile, _ := ioutil.ReadFile("docs/swagger.json")
			jsonData, _ := simplejson.NewFromReader(bytes.NewReader(jsonFile))
			urlPath := v.Path
			idPatten := "(.*)/:(\\w+)" // 正则替换，把:id换成{id}
			reg, _ := regexp.Compile(idPatten)
			if reg.MatchString(urlPath) {
				urlPath = reg.ReplaceAllString(v.Path, "${1}/{${2}}") // 把:id换成{id}
			}
			apiTitle, _ := jsonData.Get("paths").Get(urlPath).Get(strings.ToLower(v.Method)).Get("summary").String()

			err1 := global.GDb.Debug().Where(model.SysApi{Path: v.Path, Action: v.Method}).
				Attrs(model.SysApi{Handle: v.Handler, Title: apiTitle}).
				FirstOrCreate(&model.SysApi{}).
				Update("title", apiTitle).
				Error
			fmt.Println(err1)
		}
	}
}
