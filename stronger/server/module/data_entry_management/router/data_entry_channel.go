package router

import (
	"github.com/gin-gonic/gin"
	"server/module/data_entry_management/api"
)

func InitDataEntryManagement(Router *gin.RouterGroup) {
	DataEntryManagementRouter := Router.Group("data-entry") //.
	//Use(middleware.JWTAuth()).
	//Use(middleware.CasbinHandler())
	//Use(middleware.JWTAuth()).
	//Use(middleware.SysLogger(1))
	{
		//录入管理-录入通道
		DataEntryManagementRouter.GET("channel/list", api.GetDataEntryChannelInformation)
	}
}
