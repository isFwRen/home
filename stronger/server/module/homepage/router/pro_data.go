/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/3 14:48
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/homepage/api"
)

func InitProDataRouter(router *gin.RouterGroup) {
	homeRouter := router.Group("homepage/")
	{
		homeRouter.GET("pro-data/list", api.GetProData)                     // 获取项目数据待处理数据
		homeRouter.GET("pro-data/business-ranking", api.GetBusinessRanking) // 获取项目业务量趋势
		homeRouter.GET("pro-data/aging-trend", api.GetAgingTrend)           // 获取项目时效趋势
	}
}
