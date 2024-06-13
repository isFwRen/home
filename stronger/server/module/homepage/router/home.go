/**
 * @Author: xingqiyi
 * @Description: 我的主页
 * @Date: 2022/7/11 10:10
 */

package router

import (
	"github.com/gin-gonic/gin"
	api2 "server/module/homepage/api"
	"server/module/pro_manager/api"
)

func InitHomeRouter(router *gin.RouterGroup) {
	homeRouter := router.Group("homepage/")
	{
		homeRouter.GET("home/announcement", api.GetHomePageAnnouncement) // 主页获取公告
		homeRouter.POST("home/set-target", api2.SetTarget)               // 主页设置个人目标
		homeRouter.GET("home/ranking-yield", api2.GetRankingYield)       // 获取产量排行
		homeRouter.GET("home/user-yield", api2.GetUserYield)             // 个人产量,个人目标

		homeRouter.GET("home/announcement-view", api.AnnouncementView) // 公告访问量
	}
}
