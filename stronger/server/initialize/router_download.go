/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/4/12 14:45
 */

package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"server/global"
	router2 "server/module/download/router"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/middleware"
)

// RoutersDownload 初始化下载总路由
func RoutersDownload() (engine *gin.Engine) {
	router := gin.New()
	//跨域
	global.GLog.Debug("use middleware cors")
	router.Use(middleware.Cors())
	//为用户头像和文件提供静态地址
	router.Static(global.GConfig.LocalUpload.FilePath, global.GConfig.LocalUpload.FilePath)
	global.GLog.Debug("use middleware logger")
	//router.Use(middleware.GinLogger())
	//router.Use(middleware.JWTAuth())
	//Router.Use(middleware.GinRecovery(false))
	//swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GLog.Debug("register swagger handler")
	//方便统一添加路由组前缀 多服务器上线使用
	apiGroup := router.Group("/v1/")
	router2.InitDownload(apiGroup)
	return router
}
