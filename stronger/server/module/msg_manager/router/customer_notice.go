/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/18 09:24
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/msg_manager/api"
)

func InitCustomerNotice(Router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(2, nil)
	customerNoticeRouter := Router.Group("msg-manager/")
	//Use(middleware.Limiter(limiter))
	//Use(middleware.CasbinHandler())
	{

		//客户通知
		customerNoticeRouter.GET("customer-notice/page", api.GetCustomerNoticeByPage) // 获取客户通知分页
		customerNoticeRouter.POST("customer-notice/reply", api.Reply)                 // 回复
	}

}
