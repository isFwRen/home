/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/10/31 23:22
 */

package router

import (
	"github.com/gin-gonic/gin"
	"server/module/dingding/api"
)

func InitDingdingGroupRouter(Router *gin.RouterGroup)  {
	DingdingGroupsRouter := Router.Group("dingdingGroups")
		//Use(middleware.JWTAuth()).
		//Use(middleware.CasbinHandler())
	{
		//添加一条记录
		DingdingGroupsRouter.POST("addGroup", api.AddGroup)
		//分页获取所有钉钉群
		DingdingGroupsRouter.POST("getGroupList", api.GetGroupList)
		//根据id删除钉钉群
		DingdingGroupsRouter.POST("delGroup", api.DelGroup)
		DingdingGroupsRouter.POST("updateGroup", api.UpdateGroup)

	}


}
