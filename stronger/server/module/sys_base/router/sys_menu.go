/**
 * @Author: xingqiyi
 * @Description:菜单路由
 * @Date: 2022/1/5 10:30 上午
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/sys_base/api"
)

func InitMenuRouter(router *gin.RouterGroup) {
	menuRouter := router.Group("sys-menu/")
	//Use(middleware.JWTAuth()).
	//Use(middleware.CasbinHandler())
	//Use(middleware.OperationRecord())
	{
		menuRouter.POST("add", api.AddMenu)         //新增菜单或按钮
		menuRouter.DELETE("delete", api.DeleteMenu) //删除菜单或按钮
		menuRouter.GET("tree", api.GetMenuTreeList) //获取菜单或按钮
		menuRouter.POST("edit", api.EditMenu)       //修改菜单或按钮

		menuRouter.GET("role/get", api.GetMenuTreeAndProByToken) //获取当前登录角色的menu
		menuRouter.GET("api/get", api.GetApis)                   //获取所有api
	}
}
