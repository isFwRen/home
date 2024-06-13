package router

import (
	"github.com/gin-gonic/gin"
	"server/module/training_guide/api"
)

func InitRuleSchedule(r *gin.RouterGroup) gin.IRoutes {
	ruleScheduleRouter := r.Group("/training-guide")
	{
		ruleScheduleRouter.GET("/finish-read-doc", api.FinishDocReading)
		ruleScheduleRouter.POST("/finish-read", api.FinishStep) //完成文件阅读
	}
	return ruleScheduleRouter
}
