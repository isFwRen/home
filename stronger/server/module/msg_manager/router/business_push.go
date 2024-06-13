package router

import (
	"github.com/gin-gonic/gin"
	"server/module/msg_manager/api"
)

func InitBusinessPush(Router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(2, nil)
	customerNoticeRouter := Router.Group("msg-manager/")
	//Use(middleware.Limiter(limiter))
	{
		customerNoticeRouter.GET("business-push/page", api.GetBusinessPushByPage) // 获取业务通知分页
		customerNoticeRouter.POST("business-push/read", api.Read)                 // 已读业务通知
	}

}
