package router

import (
	"github.com/gin-gonic/gin"
)

func TrainingManagementRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	trainingStageRouter := Router.Group("training-management")
	{
		//trainingStageRouter.GET("/page", api.PageTrainingManagement)          //分页查询培训
		//trainingStageRouter.GET("/info", api.InfoTrainingManagement)          //查询培训详情
		//trainingStageRouter.POST("/edit", api.EditTrainingManagement)         //审核 培训管理
		//trainingStageRouter.POST("/export", api.ExportTrainingManagementInfo) //导出培训管理
	}
	return trainingStageRouter
}
