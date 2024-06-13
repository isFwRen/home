/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/3 14:51
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
	responseSysBase "server/module/sys_base/model/response"
)

// GetProData
// @Tags homepage(首页--项目数据)
// @Summary 首页--获取项目数据待处理数据
// @Auth xingqiyi
// @Date 2022/8/3 14:51
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param queryDay      query   string   false   "格式2022-07-06"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /homepage/pro-data/list [get]
func GetProData(c *gin.Context) {
	var queryDay request.QueryDayReq
	err := c.ShouldBindQuery(&queryDay)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list := service.GetProData(queryDay)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List: list,
		}, c)
	}
}

// GetBusinessRanking
// @Tags homepage(首页--项目数据)
// @Summary 首页--获取项目业务量趋势
// @Auth xingqiyi
// @Date 2022/8/3 14:51
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param rankingType      query   int   false   "0:日排行，1：月排行，2：年排行"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /homepage/pro-data/business-ranking [get]
func GetBusinessRanking(c *gin.Context) {
	var queryBusinessReq request.QueryBusinessReq
	err := c.ShouldBindQuery(&queryBusinessReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list := service.GetBusinessRanking(queryBusinessReq)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List: list,
		}, c)
	}
}

// GetAgingTrend
// @Tags homepage(首页--项目数据)
// @Summary 首页--获取项目时效趋势
// @Auth xingqiyi
// @Date 2022/8/3 16:44
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param rankingType      query   int   false   "0:日，1：月"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /homepage/pro-data/aging-trend [get]
func GetAgingTrend(c *gin.Context) {
	var queryBusinessReq request.QueryBusinessReq
	err := c.ShouldBindQuery(&queryBusinessReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list := service.GetAgingTrend(queryBusinessReq)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List: list,
		}, c)
	}
}
