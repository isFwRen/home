/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/25 13:46
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/msg_manager/api"
)

func InitGroupNotice(router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(2, nil)
	groupNoticeRouter := router.Group("msg-manager/")
	//Use(middleware.Limiter(limiter))
	//Use(middleware.CasbinHandler())
	{
		groupNoticeRouter.GET("group-notice/get-by-group-id", api.GetGroupNoticeByGroupId) // 获取
		groupNoticeRouter.POST("group-notice/add", api.AddGroupNotice)                     // 新增
		groupNoticeRouter.POST("group-notice/re", api.ReGroupNotice)                       // 重置
		groupNoticeRouter.POST("group-notice/edit", api.EditGroupNotice)                   // 编辑
	}
}
