package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_manager/api"
)

func InitTransactionRuleRouter(Router *gin.RouterGroup) {
	TransactionRouter := Router.Group("pro-manager/").
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord())
		Use(middleware.SysLogger(2))
	{
		TransactionRouter.GET("transactionRule/list", api.GetTransactionRule)
		TransactionRouter.GET("transactionRule/task/list", api.GetTaskTransactionRule)
		TransactionRouter.POST("transactionRule/add", api.AddTransactionRule)
		TransactionRouter.POST("transactionRule/edit", api.EditTransactionRule)
		TransactionRouter.POST("transactionRule/delete", api.DeleteTransactionRule)
	}
}
