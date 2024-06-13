/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/2/9 1:46 下午
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/sys_base/api"
)

func InitSysLoggerRouter(router *gin.RouterGroup) {
	sysLoggerRouter := router.Group("sys-logger/")
	//Use(middleware.JWTAuth()).
	//Use(middleware.CasbinHandler()).
	//Use(middleware.OperationRecord())
	{
		sysLoggerRouter.GET("list", api.GetPageByType) //分页获取系统日志
	}
}
