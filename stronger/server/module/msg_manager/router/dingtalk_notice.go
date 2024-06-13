/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/6 14:50
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/msg_manager/api"
)

func InitDingtalkNotice(Router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(1, nil)
	dingtalkNoticeRouter := Router.Group("msg-manager/")
	//Use(middleware.Limiter(limiter))
	//Use(middleware.CasbinHandler())
	{

		//群通知
		dingtalkNoticeRouter.POST("dingtalk-notice-msg/send", api.DingtalkNoticeSend)          // 发送群通知
		dingtalkNoticeRouter.GET("dingtalk-notice-msg/page", api.GetDingtalkNoticeSendMsgPage) // 群通知信息记录

	}

}
