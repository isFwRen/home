package router

import (
	"github.com/gin-gonic/gin"
	"server/module/msg_manager/api"
)

func InitDingtalk(Router *gin.RouterGroup) {
	//limiter := tollbooth.NewLimiter(2, nil)
	dingtalkGroupRouter := Router.Group("msg-manager/")
	//Use(middleware.Limiter(limiter))
	//Use(middleware.CasbinHandler())
	{

		//PT群配置
		dingtalkGroupRouter.GET("dingtalk-group/page", api.GetDingtalkGroupByPage)     // 获取PT群配置分页
		dingtalkGroupRouter.POST("dingtalk-group/add", api.AddDingtalkGroup)           // 新增一个PT群配置
		dingtalkGroupRouter.POST("dingtalk-group/edit", api.UpdateDingtalkGroupById)   // 根据id更新PT群配置
		dingtalkGroupRouter.DELETE("dingtalk-group/delete", api.DelDingtalkGroupByIds) // 删除PT群配置
		dingtalkGroupRouter.GET("dingtalk-group/get-dict-const", api.GetDictConst)     // 获取常量

	}

}
