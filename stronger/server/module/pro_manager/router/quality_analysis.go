package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_manager/api"
)

func InitQualityAnalysis(Router *gin.RouterGroup) {
	QualityAnalysisRouter := Router.Group("pro-manager/").
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord())
		Use(middleware.SysLogger(2))
	{
		//QualityManagementRouter.GET("quality", api.Creatable)
		QualityAnalysisRouter.GET("quality/analysis/list", api.GetQualityAnalysis)      //质量分析-按项目
		QualityAnalysisRouter.GET("quality/analysis/export", api.ExportQualityAnalysis) //质量分析-导出
	}
}
