package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/sys_base/api"
)

func InitRoleManagement(Router *gin.RouterGroup) {
	RoleManagementRouter := Router.Group("sys-base/").
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord()).
		Use(middleware.SysLogger(2))
	{
		RoleManagementRouter.GET("role-management/list", api.GetRoleInformation)         //人员管理--查询角色
		RoleManagementRouter.POST("role-management/add", api.AddRoleInformation)         //人员管理--新增角色
		RoleManagementRouter.POST("role-management/edit", api.EditRoleInformation)       //人员管理--修改角色
		RoleManagementRouter.DELETE("role-management/delete", api.DeleteRoleInformation) //人员管理--删除角色

		RoleManagementRouter.GET("role-management/sys-menu/tree", api.GetRoleMenuTree)      //获取当前角色菜单树
		RoleManagementRouter.POST("role-management/sys-menu/relation-set", api.SetRelation) //设置当前角色该菜单的权限
	}

}
