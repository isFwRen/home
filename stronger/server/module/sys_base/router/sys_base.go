package router

import (
	api2 "server/module/sys_base/api"

	_ "github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	BaseRouter := Router.Group("sys-base/")
	{
		//BaseRouter.POST("sys-register/register", api2.Register)                    //账号注册
		//BaseRouter.POST("sys-login/login", api2.Login)                             //账号登录
		//BaseRouter.POST("sys-login/login-task", api2.TaskLogin)                    //账号登录
		//BaseRouter.POST("sys-user-operation-pw/resetPassword", api2.ResetPassword) //密码重置
		BaseRouter.POST("sys-login/get-user-qrCode", api2.GetUserQrCode)

		BaseRouter.POST("sys-login/user-info", api2.GetUserInfo) //根据登录状态获取token
	}
	return BaseRouter
}
