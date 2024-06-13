/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/11 10:52
 */

package service

import (
	"gorm.io/gorm/clause"
	"server/global"
	"server/module/homepage/model"
	"server/module/homepage/model/request"
	"server/module/homepage/model/response"
	model2 "server/module/report_management/model"
	"time"
)

//SetTarget 设置个人产量目标
func SetTarget(target model.YieldTarget) error {
	return global.GDb.Model(&model.YieldTarget{}).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_code"}, {Name: "target_date"}},
		DoUpdates: clause.AssignmentColumns([]string{"target"}),
	}).Create(&target).Error
}

//GetRankingYield 获取排行榜
func GetRankingYield(yieldRankingReq request.YieldRankingReq) (err error, total int64, list []response.YieldRanking) {
	limit := yieldRankingReq.PageSize
	offset := yieldRankingReq.PageSize * (yieldRankingReq.PageIndex - 1)
	dateFormatStr := "'YYYY-MM-DD'"
	dateFormatLayoutStr := "2006-01-02"
	if yieldRankingReq.RankingType == 1 {
		dateFormatStr = "'YYYY-MM'"
		dateFormatLayoutStr = "2006-01"
	}
	nowDate := time.Now().Format(dateFormatLayoutStr)
	//global.GDb.Model(&model2.OutputStatisticsSummary{}).Where("to_char(submit_time,"+dateFormatStr+") = ?", nowDate).
	//	Group("to_char(submit_time," + dateFormatStr + ")").Count(&total)
	global.GDb.Raw("SELECT count(DISTINCT code) FROM output_statistics_summaries WHERE to_char(submit_time,"+dateFormatStr+") = ?", nowDate).Scan(&total)
	err = global.GDb.Model(&model2.OutputStatisticsSummary{}).Where("to_char(submit_time,"+dateFormatStr+") = ?", nowDate).
		Select("code as user_code,nick_name as user_name,to_char(submit_time," + dateFormatStr + ") as yield_date, sum(mary) as value").
		Group("user_code,user_name,to_char(submit_time," + dateFormatStr + ")").Order("value desc").Limit(limit).Offset(offset).Find(&list).Error
	return err, total, list
}

//GetUserYield 获取主页目标、产量和柱状图产量
func GetUserYield(code string) (err error, list []response.YieldRanking, target, yield float64) {
	nowDate := time.Now().Format("2006-01-02")
	d, _ := time.ParseDuration("-168h")
	startDate := time.Now().Add(d)
	d, _ = time.ParseDuration("-0h")
	endDate := time.Now().Add(d)
	dateFormatStr := "'YYYY-MM-DD'"
	global.GDb.Model(&model.YieldTarget{}).Select("target").Where("user_code = ? and to_char(target_date,"+dateFormatStr+") = ?", code, nowDate).First(&target)
	global.GDb.Model(&model2.OutputStatisticsSummary{}).Select("mary").Where("code = ? and to_char(submit_time,"+dateFormatStr+") = ?", code, nowDate).First(&yield)

	err = global.GDb.Model(&model2.OutputStatisticsSummary{}).Where("submit_time BETWEEN ? AND ?  and code = ?", startDate, endDate, code).
		Select("to_char(submit_time," + dateFormatStr + ") as yield_date, sum(mary) as value").
		Group("to_char(submit_time," + dateFormatStr + ")").Order("yield_date").Limit(7).Find(&list).Error
	return err, list, target, yield
}

//GetUserYieldRanking 获取个人产量排行
func GetUserYieldRanking(code string, rankingType int) (err error, userYieldRanking response.YieldRanking) {
	dateFormatStr := "'YYYY-MM-DD'"
	dateFormatLayoutStr := "2006-01-02"
	if rankingType == 1 {
		dateFormatStr = "'YYYY-MM'"
		dateFormatLayoutStr = "2006-01"
	}
	nowDate := time.Now().Format(dateFormatLayoutStr)
	//获取当前用户的产量
	err = global.GDb.Raw("SELECT  sum(mary) as value FROM output_statistics_summaries WHERE to_char(submit_time,"+dateFormatStr+") = ? and code = ? GROUP BY code,to_char(submit_time,"+dateFormatStr+")", nowDate, code).
		Scan(&userYieldRanking.Value).Error
	if err != nil {
		return err, userYieldRanking
	}
	//获取排名
	err = global.GDb.Raw("SELECT count(*) from (SELECT  sum(mary) as value FROM output_statistics_summaries WHERE to_char(submit_time,"+dateFormatStr+") = ?  GROUP BY code,to_char(submit_time,"+dateFormatStr+")) a where a.value  >= ?", nowDate, userYieldRanking.Value).
		Scan(&userYieldRanking.MyOrder).Error
	userYieldRanking.YieldDate = nowDate
	userYieldRanking.UserCode = code
	return err, userYieldRanking
}
