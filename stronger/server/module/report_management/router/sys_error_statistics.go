package router

import (
	"server/middleware"
	"server/module/report_management/api"

	"github.com/gin-gonic/gin"
)

func InitErrorStatistics(Router *gin.RouterGroup) {
	ErrorStatisticsRouter := Router.Group("report-management").
		// Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(3))
	{
		ErrorStatisticsRouter.GET("error-statistics/list", api.GetIncorrectList)                         //错误查询--查询错误明细
		ErrorStatisticsRouter.POST("error-statistics/export", api.ExportIncorrectList)                   //错误查询--导出错误明细
		ErrorStatisticsRouter.GET("error-statistics/task/list", api.GetIncorrectTaskList)                //错误查询(录入系统)--查询错误明细
		ErrorStatisticsRouter.POST("error-statistics/task/complain", api.ComplainTask)                   //错误查询(录入系统)--错误明细,是否申诉
		ErrorStatisticsRouter.POST("error-statistics/complain", api.Complain)                            //错误查询--错误明细,是否申诉
		ErrorStatisticsRouter.POST("error-statistics/wrong-confirm", api.WrongConfirm)                   //错误查询--错误明细,申诉是否通过
		ErrorStatisticsRouter.GET("error-statistics/wrong-analysis/list", api.IncorrectAnalysis)         //错误查询--查询错误分析
		ErrorStatisticsRouter.GET("error-statistics/wrong-analysis/export", api.ExportIncorrectAnalysis) //错误查询--导出获取错误分析
		ErrorStatisticsRouter.GET("error-statistics/ocr-analysis/list", api.OcrAnalysis)                 //错误查询--查询OCR错误明细

	}
}
