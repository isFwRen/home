package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_manager/api"
)

func InitTeachVideo(Router *gin.RouterGroup) {
	TeachVideoRouter := Router.Group("pro-manager/").
		Use(middleware.SysLogger(2))
	{
		TeachVideoRouter.GET("teachVideo/sync", api.BlockSync)
		TeachVideoRouter.GET("teachVideo/list", api.GetTeachVideoList)
		TeachVideoRouter.GET("teachVideo/task/list", api.GetTaskTeachVideoList)
		TeachVideoRouter.POST("teachVideo/edit", api.EditTeachVideo)
		TeachVideoRouter.POST("teachVideo/delete", api.DeleteTeachVideo)
		TeachVideoRouter.GET("teachVideo/export", api.ExportTeachVideo)
		TeachVideoRouter.POST("teachVideo/upload", api.UploadTeachVideo)
	}
}
