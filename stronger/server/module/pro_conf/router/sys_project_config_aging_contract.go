package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"

	//"server/middleware"
	"server/module/pro_conf/api"
)

func InitProjectConfigAgingContract(Router *gin.RouterGroup) {
	ProjectConfigAgingContractRouter := Router.Group("pro-config").
		//Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(2))
	{
		//时效配置--新增合同时效
		ProjectConfigAgingContractRouter.POST("project-config-aging-contract/add", api.InsertAgingContractConfig)
		////时效配置--删除合同时效配置
		ProjectConfigAgingContractRouter.DELETE("project-config-aging-contract/delete", api.DelAgingContractConfig)
		////时效配置--修改合同时效配置
		ProjectConfigAgingContractRouter.POST("project-config-aging-contract/edit", api.UpdateAgingContractConfig)
		////时效配置--查询合同时效配置
		ProjectConfigAgingContractRouter.GET("project-config-aging-contract/list", api.GetAgingContractConfig)
	}
}
