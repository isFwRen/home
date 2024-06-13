package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/sys_base/api"
)

func InitProPermissionSum(Router *gin.RouterGroup) {
	ProPermissionSumRouter := Router.Group("sys-base").
		//Use(middleware.JWTAuth()).
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord())
		//Use(middleware.JWTAuth()).
		Use(middleware.SysLogger(1))
	{
		ProPermissionSumRouter.GET("pro-permission-check/list", api.GetProPermissionSum) //人员管理--查询项目权限统计
	}
}
