/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/5 09:51
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/homepage/api"
)

func InitProReportRouter(router *gin.RouterGroup) {
	homeRouter := router.Group("homepage/")
	{
		homeRouter.GET("pro-report/list", api.GetProReport) // 获取项目日报
	}
}
