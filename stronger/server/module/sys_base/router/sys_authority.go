package router

import (
	"github.com/gin-gonic/gin"
	"server/module/sys_base/api"
)

func InitAuthorityRouter(Router *gin.RouterGroup) {
	AuthorityRouter := Router.Group("authority")
	//Use(middleware.JWTAuth()).
	//Use(middleware.CasbinHandler2()).
	//Use(middleware.OperationRecord())
	{
		AuthorityRouter.POST("createAuthority", api.CreateNewAuthority)        // 创建角色
		AuthorityRouter.POST("deleteAuthority", api.DeleteAuthority)           // 删除角色
		AuthorityRouter.POST("updateAuthority", api.UpdateAuthority)           // 更新角色信息
		AuthorityRouter.POST("updateAuthorityPower", api.UpdateAuthorityPower) // 更新角色权限
		//AuthorityRouter.POST("copyAuthority", api.CopyAuthority)       // 更新角色
		//AuthorityRouter.POST("getAuthorityList", api.GetAuthorityList) // 获取角色列表
		//AuthorityRouter.POST("setDataAuthority", api.SetDataAuthority) // 设置角色资源权限
		AuthorityRouter.POST("getAuthority", api.GetAllAuthority) //获取全部角色
	}
}
