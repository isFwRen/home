package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_manager/api"
)

func InitQualityManagement(Router *gin.RouterGroup) {
	QualityManagementRouter := Router.Group("pro-manager/").
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord())
		Use(middleware.SysLogger(2))
	{
		//QualityManagementRouter.GET("quality/creat", api.Creatable)
		QualityManagementRouter.GET("quality/management/list", api.GetQualityManagement)     //质量管理-查询
		QualityManagementRouter.POST("quality/management/add", api.AddQualityManagement)     //质量管理-新增
		QualityManagementRouter.POST("quality/management/edit", api.EditQualityManagement)   //质量管理-编辑
		QualityManagementRouter.POST("quality/management/delete", api.DeleteQualityData)     //质量管理-删除
		QualityManagementRouter.POST("quality/management/push", api.UploadQualityManagement) //质量管理-导入数据
		QualityManagementRouter.GET("quality/management/export", api.ExportQualityData)      //质量管理-导出
	}
}
