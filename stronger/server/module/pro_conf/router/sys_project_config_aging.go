package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_conf/api"
)

func InitProjectConfigAging(Router *gin.RouterGroup) {
	ProjectConfigAgingRouter := Router.Group("pro-config").
		//Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(2))
	{
		//时效配置--新增时效配置
		ProjectConfigAgingRouter.POST("project-config-aging/add", api.InsertAgingConfig)
		//时效配置--删除时效配置
		ProjectConfigAgingRouter.DELETE("project-config-aging/delete", api.DelAgingConfig)
		//时效配置--修改时效配置
		ProjectConfigAgingRouter.POST("project-config-aging/edit", api.UpdateAgingConfig)
		//时效配置--查询时效配置
		ProjectConfigAgingRouter.GET("project-config-aging/list", api.GetAgingConfigByConfigType)
	}
}
