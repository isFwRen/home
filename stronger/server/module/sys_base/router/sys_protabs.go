package router
import (
	"github.com/gin-gonic/gin"
	"server/module/sys_base/api"
)

func InitTabs(Router *gin.RouterGroup) {
	TabsRouter := Router.Group("tabs")
	{
		TabsRouter.GET("getTabsList", api.GetTabsList) //获取Tabs列表
		TabsRouter.GET("getTabsLast", api.GetTabsLast) //获取Tabs最后一条记录
		TabsRouter.POST("addTabs", api.AddTabs)        //新增Tabs
		TabsRouter.POST("removeTabs", api.RemoveTabs)  //根据ID删除Tabs
		TabsRouter.POST("updateTabs", api.UpdateTabs)  //根据name(类似ID)更新tabs
	}
}

