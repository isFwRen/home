package router

import (
	"github.com/gin-gonic/gin"
	"server/middleware"

	//"server/middleware"
	"server/module/report_management/api"
)

func InitSalary(Router *gin.RouterGroup) {
	PTSalaryRouter := Router.Group("report-management").
		//Use(middleware.CasbinHandler())
		Use(middleware.SysLogger(3))
	{
		PTSalaryRouter.GET("/pt/salary/download", api.PtSalaryFileDownload) //工资数据--导出PT工资数据
		PTSalaryRouter.POST("/pt/salary/upload", api.PtSalaryFileUpload)    //工资数据--导入PT工资数据
		PTSalaryRouter.GET("/pt/salary/list", api.GetPtSalary)              //工资数据--查询工资数据
	}

	PTSalaryTaskRouter := Router.Group("report-management").
		//Use(middleware.CasbinHandler())
		Use(middleware.SysLogger(3))
	{
		PTSalaryTaskRouter.GET("/pt/salary-task/list", api.GetPtSalaryTask) //工资数据--查询工资数据
	}

	InternalSalaryRouter := Router.Group("report-management").
		//Use(middleware.CasbinHandler())
		Use(middleware.SysLogger(3))
	{
		InternalSalaryRouter.GET("/internal/salary/download", api.InternalSalaryFileDownload) //工资数据--导出内部工资数据
		InternalSalaryRouter.POST("/internal/salary/upload", api.InternalSalaryFileUpload)    //工资数据--导入内部工资数据
		InternalSalaryRouter.GET("/internal/salary/list", api.GetInternalSalary)              //工资数据--查询内部工资数据
	}
}
