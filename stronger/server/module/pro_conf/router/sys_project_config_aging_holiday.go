package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_conf/api"
)

func InitProjectConfigAgingHoliday(Router *gin.RouterGroup) {
	ProjectConfigAgingHolidayRouter := Router.Group("pro-config").
		//Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(2))
	{
		//ProjectConfigAgingHolidayRouter.POST("project-config-aging-holiday/createProjectConfigAgingHolidayTable", api.CreateProjectConfigAgingHolidayTable)
		//节假日时效--添加节假日时效
		ProjectConfigAgingHolidayRouter.POST("project-config-aging-holiday/add", api.InsertProjectConfigAgingHoliday)
		//节假日时效--获取节假日时效
		ProjectConfigAgingHolidayRouter.GET("project-config-aging-holiday/list", api.GetProjectConfigAgingHoliday)
		//节假日时效--更新节假日时效
		ProjectConfigAgingHolidayRouter.POST("project-config-aging-holiday/edit", api.UpdateProjectConfigAgingHoliday)
	}
}
