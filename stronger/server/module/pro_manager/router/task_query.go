package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"
	"server/module/pro_manager/api"
)

func InitTaskQuery(Router *gin.RouterGroup) {
	TaskQueryRouter := Router.Group("pro-config").
		//Use(middleware.CasbinHandler()).
		//Use(middleware.OperationRecord())
		Use(middleware.SysLogger(2))
	{
		TaskQueryRouter.GET("task/detail/List", api.GetTaskListDetail)                                  //任务查询-查询待分配-已分配-缓存区-紧急件-优先件 详细信息
		TaskQueryRouter.GET("task/List", api.GetTaskList)                                               //任务查询-查询待分配-已分配-缓存区-紧急件-优先件
		TaskQueryRouter.POST("task/UrgencyBillOrPriorityBill/list", api.GetUrgencyBillOrPriorityBill)   //任务查询-查询紧急件/优先件
		TaskQueryRouter.POST("task/UrgencyBillOrPriorityBill/edit", api.SetUrgencyBillOrPriorityBill)   //任务查询-设置紧急件/优先件
		TaskQueryRouter.POST("task/PriorityOrganizationNumber/edit", api.SetPriorityOrganizationNumber) //任务查询-设置优先处理机构号
		TaskQueryRouter.GET("CaseDetails/list", api.GetCaseDetails)                                     //任务查询-案件明细
		TaskQueryRouter.GET("CaseDetails/block/list", api.GetCaseDetailsBlock)                          //任务查询-案件明细-分块明细
		TaskQueryRouter.GET("CaseDetails/field/list", api.GetCaseDetailsField)                          //任务查询-案件明细-字段明细
		TaskQueryRouter.POST("task/urgency-priority-bill/edit", api.SetUrgencyPriorityBill)             //任务查询-将单据紧急1,优先2,正常0状态
		TaskQueryRouter.GET("task/getProPermissionPeople", api.GetProPermissionPeople)                  //任务查询-获取拥有该项目各个工序的用户
	}
}
