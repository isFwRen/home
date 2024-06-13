package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/sys_base/api"
)

func InitUManagement(Router *gin.RouterGroup) {
	UManagementRouter := Router.Group("sys-base/").
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord()).
		Use(middleware.SysLogger(2))
	{
		UManagementRouter.GET("user/find", api.GetUsersInformation) //人员管理--查询员工
	}

}
