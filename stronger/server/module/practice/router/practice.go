/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 10:31 上午
 */

package router

import (
	// "server/middleware"
	"server/module/practice/api"

	// "github.com/didip/tollbooth"

	"github.com/gin-gonic/gin"
)

func InitPractice(Router *gin.RouterGroup) {
	// limiter := tollbooth.NewLimiter(1, nil)
	practiceRouter := Router.Group("practice/")
	// .
	//Use(middleware.GinRecovery(false)).
	// Use(middleware.Limiter(limiter))
	{
		practiceRouter.GET("/get", api.GetPracticeTask)
		practiceRouter.POST("/submit", api.SubmitTask)
		practiceRouter.POST("/exit", api.ExitPractice)

		practiceRouter.GET("/sum", api.GetPracticeSumByPage)
		practiceRouter.GET("/wrong", api.GetPracticeWrongByPage)

		practiceRouter.GET("/ask-list", api.GetPracticeAskList)

		// practiceRouter.GET("/sum", api.GetPracticeSumByPage)
		// practiceRouter.GET("/wrong", api.GetPracticeWrongByPage)
	}
}
