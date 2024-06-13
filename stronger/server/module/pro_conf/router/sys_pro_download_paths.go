/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/30 15:10
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_conf/api"
)

func SysProDownloadPaths(Router *gin.RouterGroup) {
	sysProDownloadPathsRouter := Router.Group("pro-config/").
		//Use(middleware.JWTAuth()).
		//Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(1))
	{
		//根据项目id的下载配置
		sysProDownloadPathsRouter.GET("sys-pro-download-paths/list", api.SysProDownloadPaths)
		//根据勾选项更改下载路径
		sysProDownloadPathsRouter.POST("sys-pro-download-paths/set-available", api.SetDownloadPathAvailable)
		//获取程序进程开启情况
		sysProDownloadPathsRouter.GET("sys-pro-download-paths/process/list", api.GetProcessList)
	}
}
