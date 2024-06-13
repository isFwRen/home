/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/5 09:52
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/global/response"
	"server/module/homepage/model/request"
	"server/module/homepage/service"
)

// GetProReport
// @Tags homepage(首页--项目日报)
// @Summary 首页--获取项目日报
// @Auth xingqiyi
// @Date 2022/8/5 09:52
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param reportDay       query   string   true    "项目日报日期格式2022-07-06T16:00:00.000Z"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /homepage/pro-report/list [get]
func GetProReport(c *gin.Context) {
	var proReportReq request.ProReportReq
	err := c.ShouldBindQuery(&proReportReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list, otherInfo := service.GetProReport(proReportReq)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithoutCode(err.Error(), fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(map[string]interface{}{
			"otherInfo": otherInfo,
			"list":      list,
		}, c)
	}
}
