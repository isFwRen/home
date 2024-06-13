/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 10:31 上午
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/task/api"
)

func InitCommonRouter(Router *gin.RouterGroup) {

	//limiter := tollbooth.NewLimiter(1, nil)
	billListRouter := Router.Group("task/common/")
	//Use(middleware.GinRecovery(false)).
	//Use(middleware.Limiter(limiter))
	{

		billListRouter.GET("bill-list/page", api.GetTaskBillByPage) //案件列表
	}
}
