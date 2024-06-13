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

func GetSysSpotCheckWrongByPage(c *gin.Context) {
	var SysSpotCheckWrongQuery request2.SysSpotCheckWrongQuery
	err := c.ShouldBindQuery(&SysSpotCheckWrongQuery)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, list := service.GetSysSpotCheckWrongByPage(SysSpotCheckWrongQuery)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     total,
			PageIndex: SysSpotCheckWrongQuery.PageIndex,
			PageSize:  SysSpotCheckWrongQuery.PageSize,
		}, c)
	}
}
