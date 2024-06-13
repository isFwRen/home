package api

import (
	"server/global/response"
	request2 "server/module/pro_manager/model/request"
	"server/module/pro_manager/service"
	responseSysBase "server/module/sys_base/model/response"

	"github.com/gin-gonic/gin"
)

func GetSysSpotCheckStatisticByPage(c *gin.Context) {
	var SysSpotCheckStatisticQuery request2.SysSpotCheckStatisticQuery
	err := c.ShouldBindQuery(&SysSpotCheckStatisticQuery)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, list := service.GetStaticByType(SysSpotCheckStatisticQuery)
	response.OkWithData(responseSysBase.BasePageResult{
		List:  list,
		Total: total,
	}, c)
	// err, total, list := service.GetSysSpotCheckStatisticByPage(SysSpotCheckStatisticQuery)
	// if err != nil {
	// 	global.GLog.Error("", zap.Error(err))
	// 	response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	// } else {
	// 	response.OkWithData(responseSysBase.BasePageResult{
	// 		List:      list,
	// 		Total:     total,
	// 		PageIndex: SysSpotCheckStatisticQuery.PageIndex,
	// 		PageSize:  SysSpotCheckStatisticQuery.PageSize,
	// 	}, c)
	// }
}
