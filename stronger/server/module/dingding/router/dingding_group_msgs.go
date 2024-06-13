/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/10/31 23:22
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/dingding/api"
)

func InitDingdingGroupMsgsRouter(Router *gin.RouterGroup)  {
	DingdingGroupMsgsRouter := Router.Group("dingdingGroupMsgs")
		//Use(middleware.JWTAuth()).
		//Use(middleware.CasbinHandler())
	{
		//添加一条记录
		DingdingGroupMsgsRouter.POST("sendDingdingGroupMsg", api.SendDingdingGroupMsg)
		//分页获取钉钉消息
		DingdingGroupMsgsRouter.POST("getDingdingGroupMsgList", api.GetDingdingGroupMsgList)

	}


}
