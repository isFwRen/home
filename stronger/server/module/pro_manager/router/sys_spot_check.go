/**
 * @Author: xingqiyi
 * @Description:公告管理
 * @Date: 2022/7/7 09:48
 */

package router

import (
	"server/module/pro_manager/api"

	"github.com/gin-gonic/gin"
)

func InitSysSpotCheckRouter(router *gin.RouterGroup) {
	sysSpotCheckRouter := router.Group("pro-manager/")
	{
		sysSpotCheckRouter.GET("sys-spot-check/page", api.GetSysSpotCheckByPage) // 质检配置
		sysSpotCheckRouter.POST("sys-spot-check/add", api.AddSysSpotCheck)       // 新增
		sysSpotCheckRouter.POST("sys-spot-check/edit", api.UpdateSysSpotCheck)   // 编辑

		sysSpotCheckRouter.GET("sys-spot-check-data/page", api.GetSysSpotCheckDataByPage) // 抽检数据

		sysSpotCheckRouter.GET("sys-spot-check-wrong/page", api.GetSysSpotCheckWrongByPage) // 错误明细

		sysSpotCheckRouter.GET("sys-spot-check-statistic/find", api.GetSysSpotCheckStatisticByPage) // 统计
	}
}
