package router

import (
	"server/middleware"
	"server/module/report_management/api"

	"github.com/gin-gonic/gin"
)

func InitBusinessAnalysis(Router *gin.RouterGroup) {
	BusinessAnalysisRouter := Router.Group("report-management").
		// Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(3))
	{
		BusinessAnalysisRouter.GET("business-analysis/download/list", api.GetBusinessDownloadAnalysis)      //来量分析--查询来量分析
		BusinessAnalysisRouter.GET("business-analysis/upload/list", api.GetBusinessUploadAnalysis)          //回传分析--查询回传分析
		BusinessAnalysisRouter.GET("business-analysis/download/export", api.ExportBusinessDownloadAnalysis) //来量分析--导出来量分析
		BusinessAnalysisRouter.GET("business-analysis/upload/export", api.ExportBusinessUploadAnalysis)     //回传分析--导出回传分析
	}
}
