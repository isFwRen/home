package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_manager/api"
)

func InitReimbursementFormTemplate(Router *gin.RouterGroup) {
	ReimbursementFormTemplateRouter := Router.Group("pro-manager/").
		//Use(middleware.CasbinHandler())
		Use(middleware.SysLogger(2))
	{
		ReimbursementFormTemplateRouter.GET("reimbursementFormTemplate/list", api.GetReimbursementFormTemplate)
		ReimbursementFormTemplateRouter.GET("reimbursementFormTemplate/task/list", api.GetTaskReimbursementFormTemplate)
		ReimbursementFormTemplateRouter.POST("reimbursementFormTemplate/add", api.AddReimbursementFormTemplate)
		ReimbursementFormTemplateRouter.POST("reimbursementFormTemplate/delete", api.DeleteReimbursementFormTemplate)
		ReimbursementFormTemplateRouter.POST("reimbursementFormTemplate/rename", api.ReNameReimbursementFormTemplate)
	}
}
