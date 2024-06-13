package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/sys_base/api"
)

func InitUserManagement(Router *gin.RouterGroup) {
	UserManagementRouter := Router.Group("sys-base").
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord()).
		Use(middleware.SysLogger(1))
	{
		UserManagementRouter.GET("user-management/list", api.GetUserInformation)                                          //人员管理--查询用户
		UserManagementRouter.GET("user-management/user-pro-permission/list", api.GetUserPermissionInformation)            //人员管理--查询用户权限
		UserManagementRouter.POST("user-management/add", api.AddUserInformation)                                          //人员管理--新增用户
		UserManagementRouter.POST("user-management/update-user-pro-permission", api.UpdateProPermissionInUserInformation) //人员管理--更新用户项目权限
		UserManagementRouter.POST("user-management/sys-user/change-role", api.ChangeRole)                                 //人员管理--修改角色

		UserManagementRouter.GET("user-management/sys-pro-permission/export", api.SysProPermissionExport)  //项目管理--导出项目权限
		UserManagementRouter.POST("user-management/sys-pro-permission/import", api.SysProPermissionImport) //项目管理--导入项目权限
	}

}
