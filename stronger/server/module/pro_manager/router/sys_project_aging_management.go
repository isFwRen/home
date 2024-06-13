package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_manager/api"
)

func InitProjectAgingManagement(Router *gin.RouterGroup) {
	ProjectAgingManagementRouter := Router.Group("pro-config").
		//Use(middleware.CasbinHandler())
		Use(middleware.SysLogger(2))
	{
		//时效管理--查询所有未回传的单子
		//ProjectAgingManagementRouter.GET("project-aging-management/list/Only", api.GetProjectAging)
		ProjectAgingManagementRouter.GET("project-aging-management/list", api.GetProjectAgingAll)
	}
}
