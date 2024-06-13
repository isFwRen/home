package router

import (
	"server/module/training_guide/api"

	"github.com/gin-gonic/gin"
)

func InitTrainingStageRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	trainingStageRouter := Router.Group("/training-guide")
	{
		trainingStageRouter.POST("/training-stage/find", api.AddAndFindGuide)      //添加培训阶段
		trainingStageRouter.PUT("/training-stage/update", api.UpdateTrainingStage) //更新培训阶段
	}
	return trainingStageRouter
}
