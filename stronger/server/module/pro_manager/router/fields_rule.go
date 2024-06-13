package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_manager/api"
)

func InitFieldsRule(Router *gin.RouterGroup) {
	FieldsRuleRouter := Router.Group("pro-manager/").
		Use(middleware.SysLogger(2))
	{
		FieldsRuleRouter.GET("fieldsRule/sync", api.FieldsRuleSync)
		FieldsRuleRouter.GET("fieldsRule/list", api.GetFieldsRuleList)
		FieldsRuleRouter.POST("fieldsRule/edit", api.EditFieldsRule)
		FieldsRuleRouter.POST("fieldsRule/delete", api.DeleteFieldsRule)
		FieldsRuleRouter.GET("fieldsRule/export", api.ExportFieldsRule)
		FieldsRuleRouter.POST("fieldsRule/upload", api.UploadFieldsRule)
	}
}
