package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/report_management/api"
)

func InitAbnormalBill(Router *gin.RouterGroup) {
	AbnormalBillRouter := Router.Group("report-management").
		//Use(middleware.CasbinHandler()).
		Use(middleware.SysLogger(3))
	{
		AbnormalBillRouter.GET("abnormal-bill/list", api.GetAbnormalBill)      //特殊报表--查询异常件数据
		AbnormalBillRouter.GET("abnormal-bill/export", api.ExportAbnormalBill) //特殊报表--导出异常件数据
	}
}
