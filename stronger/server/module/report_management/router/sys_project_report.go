package router

import (
	"github.com/gin-gonic/gin"
	"server/module/report_management/api"
)

func InitProjectReport(Router *gin.RouterGroup) {
	InitProjectReportRouter := Router.Group("report-management/") //.
	//Use(middleware.CasbinHandler()).
	//Use(middleware.SysLogger(3))
	{
		InitProjectReportRouter.GET("project-report/business-details/list", api.GetBusinessDetails)      //项目报表--查看业务明细表
		InitProjectReportRouter.GET("project-report/business-details/export", api.ExportBusinessDetails) //项目报表--导出业务明细表

		InitProjectReportRouter.GET("project-report/business-report", api.GetBusinessReport)  //项目报表--业务量报表
		InitProjectReportRouter.GET("project-report/aging-report", api.GetAgingReport)        //项目报表--时效报表
		InitProjectReportRouter.GET("project-report/deal-time-report", api.GetDealTimeReport) //项目报表--处理时长报表

		InitProjectReportRouter.GET("project-report/report-export", api.ExportReport) //项目报表--导出报表

		InitProjectReportRouter.POST("project-report/set-report-info", api.SetOtherReportInfo) //项目报表--设置项目日报其他信息
		InitProjectReportRouter.GET("project-report/report-info-export", api.ReportInfoExport) //项目报表--导出日报

		InitProjectReportRouter.GET("project-report/get-char-sum", api.GetCharSum)       //项目报表--获取字符统计
		InitProjectReportRouter.GET("project-report/export-char-sum", api.ExportCharSum) //项目报表--导出字符统计
	}
}
