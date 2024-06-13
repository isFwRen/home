/**
 * @Author: xingqiyi
 * @Description:公告管理
 * @Date: 2022/7/7 09:48
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/pro_manager/api"
)

func InitAnnouncementRouter(router *gin.RouterGroup) {
	announcementsRouter := router.Group("pro-manager/")
	{
		announcementsRouter.GET("announcement-manager/page", api.GetAnnouncementByPage)                  // 获取公告管理分页
		announcementsRouter.POST("announcement-manager/add", api.AddAnnouncement)                        // 新增一个公告
		announcementsRouter.POST("announcement-manager/change-status", api.ChangeStatusAnnouncementById) // 发布、取消发布、删除、恢复
		announcementsRouter.POST("announcement-manager/edit", api.UpdateAnnouncementById)                // 编辑
	}
}
