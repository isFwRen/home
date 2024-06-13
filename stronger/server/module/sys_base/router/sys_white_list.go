package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/sys_base/api"
)

func InitWhiteListRouter(Router *gin.RouterGroup) {
	WhiteListRouter := Router.Group("sys-base").
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord())
		Use(middleware.SysLogger(1))
	{
		WhiteListRouter.GET("white-list/list", api.GetWhileList)                   //人员管理--查询白名单
		WhiteListRouter.POST("white-list/edit", api.EditWhileList)                 //人员管理--修改白名单
		WhiteListRouter.POST("white-list/copy", api.CopyWhileList)                 //人员管理--复制白名单
		WhiteListRouter.POST("white-list/export", api.ExportWhileList)             //人员管理--导出白名单
		WhiteListRouter.GET("white-list/getBlockPeopleSum", api.GetBlockPeopleSum) //人员管理--导出白名单
	}
}
