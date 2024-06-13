package router

import (
	"server/middleware"
	"server/module/report_management/api"

	"github.com/gin-gonic/gin"
)

func InitOutputStatistics(Router *gin.RouterGroup) {
	//人员产量统计
	OutputStatisticsRouter := Router.Group("report-management").
		Use(middleware.SysLogger(3))
	{
		//OutputStatisticsRouter.GET("output-Statistics/CreateTable", api.CreateTable)
		OutputStatisticsRouter.GET("output-Statistics/list", api.GetOutputStatistics)                   //产量统计--查询人员产量统计
		OutputStatisticsRouter.GET("output-Statistics/update", api.UpdateOutputStatistics)              //产量统计--更新产量统计
		OutputStatisticsRouter.GET("output-Statistics/export", api.ExportOutputStatistics)              //产量统计--导出人员(全部)产量统计
		OutputStatisticsRouter.GET("output-Statistics-detail/export", api.ExportOutputStatisticsDetail) //产量统计--导出人员(明细)产量统计
		OutputStatisticsRouter.POST("output-Statistics/delete", api.DeleteOutputStatisticsDetail)       //产量统计--清空人员(全部)产量统计
	}

	//人员产量统计--录入系统
	OutputStatisticsTaskRouter := Router.Group("report-management").
		Use(middleware.SysLogger(3))
	{
		OutputStatisticsTaskRouter.GET("output-Statistics-task/list", api.GetOutputStatisticsTask) //产量统计--查询人员产量统计
	}

	//OCR产量统计
	OcrOutputStatisticsRouter := Router.Group("report-management").
		Use(middleware.SysLogger(3))
	{
		OcrOutputStatisticsRouter.GET("ocr-output-Statistics/list", api.GetOCROutputStatistics) //产量统计--查询OCR产量统计
		OcrOutputStatisticsRouter.GET("ocr-output-Statistics/export", api.ExportOcrOutput)      //产量统计--导出OCR产量统计
	}

	//折算比例
	CorrectedRouter := Router.Group("report-management").
		Use(middleware.SysLogger(3))
	{
		CorrectedRouter.GET("corrected/list-log", api.GetEditCorrectedLog) //产量统计--获取更改项目折算比例日志
		CorrectedRouter.GET("corrected/list", api.GetCorrected)            //产量统计--获取项目折算比例
		CorrectedRouter.POST("corrected/edit", api.UpdateCorrected)        //产量统计--更新项目折算比例
		CorrectedRouter.POST("corrected/add", api.InsertCorrected)         //产量统计--增加项目折算比例
		CorrectedRouter.DELETE("corrected/delete", api.DeleteCorrected)    //产量统计--删除项目折算比例
	}
}
