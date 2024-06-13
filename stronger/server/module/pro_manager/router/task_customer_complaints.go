package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_manager/api"
)

func InitTaskCustomerComplaints(Router *gin.RouterGroup) {
	TaskCustomerComplaintsRouter := Router.Group("task/").
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord())
		Use(middleware.SysLogger(2))
	{
		TaskCustomerComplaintsRouter.GET("customer_complaints/list", api.GetCustomerComplaints) //客户投诉-查询
	}
}
