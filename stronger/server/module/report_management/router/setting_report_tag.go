package router

import (
	"github.com/gin-gonic/gin"
	"server/module/report_management/api"
)

func InitSettingReportTag(Router *gin.RouterGroup) {
	settingReportTagRouter := Router.Group("setting-report-tag/")
	{
		settingReportTagRouter.GET("get-tag-list", api.GetTagList)
		settingReportTagRouter.POST("get-user-tags", api.GetUserTags)
		settingReportTagRouter.POST("set-user-tags", api.SetUserTags)

	}
}
