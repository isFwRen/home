package router

import (
	"github.com/gin-gonic/gin"
	"server/module/assessment_management/api"
)

func InitAssessLoggingRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	assessLoggingRouter := Router.Group("/assessment-exam")
	{
		assessLoggingRouter.POST("/test-procedure/start-exam", api.StartExam) //开始考核
		//assessLoggingRouter.POST("/test-procedure/stop-exam")                        //暂停考核
		assessLoggingRouter.POST("/test-procedure/end-exam", api.EndExam)            //考核结束计算分数// 结束考核
		assessLoggingRouter.GET("/test-procedure/get-project-list", api.GetExamList) //获取题目列表
	}
	return assessLoggingRouter
}
