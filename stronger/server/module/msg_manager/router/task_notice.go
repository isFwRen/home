/**
 * @Author: xingqiyi
 * @Description:录入通知
 * @Date: 2022/7/18 14:40
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/msg_manager/api"
)

func InitTaskNotice(Router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(1, nil)
	taskNoticeRouter := Router.Group("msg-manager/")
	//Use(middleware.Limiter(limiter))
	//Use(middleware.CasbinHandler())
	{
		taskNoticeRouter.POST("task-notice-msg/send", api.TaskNoticeSend)          // 发送录入通知
		taskNoticeRouter.GET("task-notice-msg/page", api.GetTaskNoticeSendMsgPage) // 录入通知信息记录
	}
}
