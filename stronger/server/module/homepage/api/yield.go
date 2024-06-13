/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/11 10:39
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"server/global"
	"server/global/response"
	"server/module/homepage/model"
	"server/module/homepage/model/request"
	response2 "server/module/homepage/model/response"
	"server/module/homepage/service"
	"time"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// SetTarget
// @Tags homepage(首页--我的主页)
// @Summary 主页--设置个人目标
// @Auth xingqiyi
// @Date 2022/7/11 10:39
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.Target true "目标实体"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /homepage/home/set-target [post]
func SetTarget(c *gin.Context) {
	var t request.Target
	err := c.ShouldBindJSON(&t)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	nowDate, err := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	var target = model.YieldTarget{
		UserCode:   customClaims.Code,
		UserName:   customClaims.NickName,
		Target:     t.Target,
		TargetDate: nowDate,
	}
	err = service.SetTarget(target)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("设置失败，%v", err), c)
	} else {
		response.OkWithData("设置成功", c)
	}
}

// GetRankingYield
// @Tags homepage(首页--我的主页)
// @Summary 主页--获取产量排行
// @Auth xingqiyi
// @Date 2022/7/11 16:17
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param rankingType    query   int   	  true    "排行类型0:日排行，1：月排行"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /homepage/home/ranking-yield [get]
func GetRankingYield(c *gin.Context) {
	var yieldRankingReq request.YieldRankingReq
	err := c.ShouldBindQuery(&yieldRankingReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	//获取登录者产值排行信息
	err, userYieldRanking := service.GetUserYieldRanking(customClaims.Code, yieldRankingReq.RankingType)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	userYieldRanking.UserName = customClaims.NickName

	err, total, list := service.GetRankingYield(yieldRankingReq)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		var r response2.RankingResult
		r.List = list
		r.Total = total
		r.UserYieldRanking = userYieldRanking
		r.PageSize = yieldRankingReq.PageIndex
		r.PageIndex = yieldRankingReq.PageSize
		response.OkWithData(r, c)
	}
}

// GetUserYield
// @Tags homepage(首页--我的主页)
// @Summary 主页--个人产量和目标
// @Auth xingqiyi
// @Date 2022/7/11 16:17
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /homepage/home/user-yield [get]
func GetUserYield(c *gin.Context) {
	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	//获取目标，当天，7天内的产值
	err, list, target, yield := service.GetUserYield(customClaims.Code)
	ad, _ := time.ParseDuration("-24h")

	//构造没有数据的日期
	var newRankings []response2.YieldRanking
	for i := 0; i < 7; i++ {
		d := time.Now().Add(ad * time.Duration(i)).Format("2006-01-02")
		newRanking := response2.YieldRanking{
			YieldDate: d,
			Value:     0,
		}
		for _, ranking := range list {
			if ranking.YieldDate == d {
				newRanking.Value = ranking.Value
				continue
			}
		}
		newRankings = append(newRankings, newRanking)
	}
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(map[string]interface{}{
			"list":   newRankings,
			"target": target,
			"yield":  yield,
		}, c)
	}
}
