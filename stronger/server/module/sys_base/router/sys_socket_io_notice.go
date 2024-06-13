/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/27 09:57
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/sys_base/api"
)

func InitSocketIoNotice(Router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(1, nil)
	socketIoNoticeRouter := Router.Group("sys-socket-io-notice/")
	//Use(middleware.GinRecovery(false)).
	//Use(middleware.Limiter(limiter)).
	//Use(middleware.CasbinHandler()).
	{
		socketIoNoticeRouter.POST("customer-notice", api.CustomerNotice) //客户消息2通知
		socketIoNoticeRouter.POST("business-push", api.BusinessPush)     //业务通知
	}
}
