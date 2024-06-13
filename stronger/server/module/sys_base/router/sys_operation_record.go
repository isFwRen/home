package router

import (
	"github.com/gin-gonic/gin"
	"server/module/sys_base/api"
)

func InitSysOperationRecordRouter(Router *gin.RouterGroup) {
	SysOperationRecordRouter := Router.Group("sysOperationRecord")
	//Use(middleware.JWTAuth())
	//Use(middleware.CasbinHandler())
	{
		SysOperationRecordRouter.POST("createSysOperationRecord", api.CreateSysOperationRecord)             // 新建SysOperationRecord
		SysOperationRecordRouter.DELETE("deleteSysOperationRecord", api.DeleteSysOperationRecord)           // 删除SysOperationRecord
		SysOperationRecordRouter.DELETE("deleteSysOperationRecordByIds", api.DeleteSysOperationRecordByIds) // 批量删除SysOperationRecord
		SysOperationRecordRouter.PUT("updateSysOperationRecord", api.UpdateSysOperationRecord)              // 更新SysOperationRecord
		SysOperationRecordRouter.GET("findSysOperationRecord", api.FindSysOperationRecord)                  // 根据ID获取SysOperationRecord
		SysOperationRecordRouter.GET("getSysOperationRecordList", api.GetSysOperationRecordList)            // 获取SysOperationRecord列表

	}
}
