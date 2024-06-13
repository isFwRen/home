package initialize

import (
	"server/global"
	routerAssessmentManagement "server/module/assessment_management/router"
	routerDataEntry "server/module/data_entry_management/router"
	routerDingDing "server/module/dingding/router"
	router3 "server/module/homepage/router"
	practiceRouter "server/module/practice/router"
	router2 "server/module/pro_manager/router"
	routerReport "server/module/report_management/router"
	routerSysBase "server/module/sys_base/router"
	"server/module/task/router"
	routerPtTrainStage "server/module/training_guide/router"
	xingqiyiMiddleware "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/middleware"
	xingqiyiRouter "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/router"

	// _ "server/docs_task"

	"github.com/gin-gonic/gin"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// "github.com/swaggo/gin-swagger/swaggerFiles"
)

// 初始化总路由

func RoutersTask() *gin.Engine {
	//db_config.InitSqlite()
	var Router = gin.New()
	//Router.Use(middleware.GinLogger())
	//Router.Use(middleware.Cors())
	//Router.Use(middleware.JWTAuth())

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

	//为用户头像和文件提供静态地址
	Router.Static(global.GConfig.LocalUpload.FilePath, global.GConfig.LocalUpload.FilePath)
	//swagger
	// Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// global.GLog.Debug("register swagger handler")

	//方便统一添加路由组前缀 多服务器上线使用
	ApiGroup := Router.Group("")
	//登录、注册等基础接口
	xingqiyiRouter.InitBaseRouter(ApiGroup)
	//案件列表
	router2.InitBillList(ApiGroup)
	router.InitTask(ApiGroup)
	//钉钉验证码
	routerDingDing.InitDingingBaseRouter(ApiGroup)
	//注册用户路由
	routerSysBase.InitUserRouter(ApiGroup)
	//注册基础功能路由 不做鉴权
	routerSysBase.InitBaseRouter(ApiGroup)
	//错误查询
	routerReport.InitErrorStatistics(ApiGroup)
	routerReport.InitOutputStatistics(ApiGroup)
	routerReport.InitSalary(ApiGroup)
	routerDataEntry.InitDataEntryManagement(ApiGroup)
	//业务规则
	router2.InitTransactionRuleRouter(ApiGroup)
	//报销单模板
	router2.InitReimbursementFormTemplate(ApiGroup)
	//教学视频
	router2.InitTeachVideo(ApiGroup)
	//字段规则
	router2.InitFieldsRule(ApiGroup)
	//客户投诉
	router2.InitTaskCustomerComplaints(ApiGroup)
	//角色管理
	//routerSysBase.InitRoleManagement(ApiGroup)
	//用户管理
	//routerSysBase.InitUserManagement(ApiGroup)
	//菜单管理
	routerSysBase.InitMenuRouter(ApiGroup)

	//PT培训流程指引 -----*测试 正式需转移到 router_task.go
	routerPtTrainStage.InitTrainingStageRouter(ApiGroup)
	routerPtTrainStage.InitRuleSchedule(ApiGroup)
	//考核管理
	routerAssessmentManagement.InitAssessLoggingRouter(ApiGroup)
	//培训管理
	//routerTraining.TrainingManagementRouter(ApiGroup)
	//练习
	practiceRouter.InitPractice(ApiGroup)

	//二期
	//主页管理
	router3.InitHomeRouter(ApiGroup)

	global.GLog.Info("upload router register success")
	return Router
}
