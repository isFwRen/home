package api

import (
	"fmt"
	"server/global"
	"server/global/response"
	request2 "server/module/pro_manager/model/request"
	"server/module/pro_manager/service"
	responseSysBase "server/module/sys_base/model/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GetSysSpotCheckDataByPage(c *gin.Context) {
	var SysSpotCheckDataQuery request2.SysSpotCheckDataQuery
	err := c.ShouldBindQuery(&SysSpotCheckDataQuery)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, list := service.GetSysSpotCheckDataByPage(SysSpotCheckDataQuery)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     total,
			PageIndex: SysSpotCheckDataQuery.PageIndex,
			PageSize:  SysSpotCheckDataQuery.PageSize,
		}, c)
	}
}

func GetSysSpotCheckDataBlock(c *gin.Context) {
	var SysSpotCheckDataQuery request2.SysSpotCheckDataQuery
	err := c.ShouldBindQuery(&SysSpotCheckDataQuery)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	// err, total, list := service.GetSysSpotCheckDataByPage(SysSpotCheckDataQuery)
	// if err != nil {
	// 	global.GLog.Error("", zap.Error(err))
	// 	response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	// }
}
