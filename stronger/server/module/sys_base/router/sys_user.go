package router

import (
	"github.com/gin-gonic/gin"
	"server/module/sys_base/api"
)

func InitUserRouter(Router *gin.RouterGroup) {
	UserOnlineRouter := Router.Group("sys-user/")
	//Use(middleware.JWTAuth()).
	//Use(middleware.CasbinHandler()).
	//Use(middleware.OperationRecord())
	{
		//UserRouter.POST("uploadHeaderImg", v1.UploadHeaderImg)   // 上传头像
		//UserOnlineRouter.POST("sys-base/getUserList", api.GetUserList) // 分页获取用户列表
		UserOnlineRouter.POST("uploadUserImage", api.UploadUserImage)     //图片上传
		UserOnlineRouter.GET("sync", api.Sync)                            //同步新增的数据
		UserOnlineRouter.POST("upload-user/avatar", api.UploadUserAvatar) //上传用户头像
	}

	UserChangePwd := Router.Group("sys-user/")
	//因为需要修改密码，暂时停止用户权限
	//Use(middleware.CasbinHandler()).
	//Use(middleware.OperationRecord())
	{
		UserChangePwd.POST("sys-user-operation-pw/changePassword", api.ChangePassword) // 修改密码
	}

	UserLeave := Router.Group("sys-user-leave/")
	//操作权限：本人或者管理员
	//Use(middleware.JWTAuth()).
	//Use(middleware.CasbinHandler()).
	//Use(middleware.OperationRecord())
	{
		UserLeave.POST("user-leave/resignation", api.Resignation) //提交离职
	}

	SysRootOperationRouter := Router.Group("sys-root-operation")
	//操作权限：管理员
	//Use(middleware.JWTAuth()).
	//Use(middleware.CasbinHandler()).
	//Use(middleware.OperationRecord())
	{
		//SysRootOperationRouter.DELETE("root-delete-user/deleteUser", api.DeleteUser)                  // 删除用户
		//SysRootOperationRouter.POST("root-set-user-authority/setUserAuthority", api.SetUserAuthority) // 设置用户权限
		SysRootOperationRouter.POST("root-query-user/queryUser", api.QueryUser) //工号查询
	}
}
