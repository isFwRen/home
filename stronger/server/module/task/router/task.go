package router

import (
	"server/middleware"
	"server/module/task/api"

	"github.com/gin-gonic/gin"
)

func InitTask(Router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(1, nil)
	taskRouter := Router.Group("task/").
		//Use(middleware.GinRecovery(false)).
		//Use(middleware.JWTAuth()).
		Use(middleware.SysLogger(4)).
		Use(middleware.PermHandler())
	//Use(middleware.CasbinHandler()).
	//Use(middleware.OperationRecord()) RefleshProConf
	{
		taskRouter.GET("op", api.GetOpTask)
		taskRouter.GET("modifyBlock", api.GetOpModifyBlock)
		taskRouter.GET("getImageSize", api.GetImageSzie)
		taskRouter.GET("opNum", api.GetOpNum)
		taskRouter.POST("submit", api.SubmitOpTask)
		taskRouter.POST("uploadImage", api.UploadImage)

		taskRouter.POST("keyLog", api.KeyLog)

		taskRouter.POST("releaseBill", api.ReleaseBill)
		taskRouter.POST("releaseBlock", api.ReleaseBlock)
		taskRouter.POST("releaseExitBlock", api.ReleaseExitBlock)

		taskRouter.GET("conf", api.GetProConf)
		taskRouter.GET("get-thumbnail", api.GetThumbnail)
		taskRouter.GET("get-image-by-index", api.GetImageByIndex)
		taskRouter.GET("get-image-by-block-id", api.GetImageByBlockId)

		//项目配置
		// sysProjectRouter.GET("sys-project/page", api.GetSysProjectByPage)   // 获取项目配置分页
		// sysProjectRouter.GET("sys-project/list", api.GetSysProjectList)     // 获取项目配置
		// sysProjectRouter.POST("sys-project/add", api.AddSysProject)         // 新增一个项目配置
		// sysProjectRouter.POST("sys-project/edit", api.UpdateSysProjectById) // 根据id更新
		//sysProjectRouter.POST("rmSysProjectByIds", api.RmSysProjectByIds)       // 根据id删除
	}
	confRouter := Router.Group("task/conf/").
		//Use(middleware.GinRecovery(false)).
		//Use(middleware.JWTAuth()).
		Use(middleware.SysLogger(3))
	{
		confRouter.POST("refresh-pro-conf", api.RefreshProConf)
	}

}
