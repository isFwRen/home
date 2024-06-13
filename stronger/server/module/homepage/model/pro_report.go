/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/5 10:04
 */

package model

import (
	modelBase "server/module/sys_base/model"
	"time"
)

type ProReport struct {
	modelBase.Model
	ProCode           string    `json:"proCode" form:"proCode" excel:"项目编码"`                       //项目编码
	PredictValue      float64   `json:"predictValue" form:"predictValue" excel:"预估产值"`             //预估产值
	FinishValue       float64   `json:"finishValue" form:"finishValue" excel:"实际完成"`               //实际完成
	TimePercent       float64   `json:"timePercent" form:"timePercent" excel:"时间比例"`               //时间比例
	FinishPercent     float64   `json:"finishPercent" form:"finishPercent" excel:"完成比例"`           //完成比例
	MonthCount        int64     `json:"monthCount" form:"monthCount" excel:"月业务总量"`                //月业务总量
	DayCount          int64     `json:"dayCount" form:"dayCount" excel:"日业务量"`                     //日业务量
	MonthAgingPercent float64   `json:"monthAgingPercent" form:"monthAgingPercent" excel:"月时效保障率"` //月时效保障率
	MonthTimeoutCount int64     `json:"monthTimeoutCount" form:"monthTimeoutCount" excel:"月超时数量"`  //月超时数量
	DayTimeoutCount   int64     `json:"dayTimeoutCount" form:"dayTimeoutCount" excel:"日超时数量"`      //日超时数量
	MonthRightPercent float64   `json:"monthRightPercent" form:"monthRightPercent" excel:"月质量准确率"` //月质量准确率
	MonthErrorCount   int64     `json:"monthErrorCount" form:"monthErrorCount" excel:"月差错数量"`      //月差错数量
	DayErrorCount     int64     `json:"dayErrorCount" form:"dayErrorCount" excel:"日差错数量"`          //日差错数量
	ReportDate        time.Time `json:"reportDate" form:"reportDate" excel:"报表日期"`                 //报表日期
}

type ProReportOtherInfo struct {
	modelBase.Model
	ReportDate      time.Time `json:"reportDate" form:"reportDate"`           //报表日期
	UserCount       int64     `json:"userCount" form:"userCount"`             //编制人数
	ActiveUserCount int64     `json:"activeUserCount" form:"activeUserCount"` //实到人数
	ClosingTime     string    `json:"closingTime" form:"closingTime"`         //下班时间
	OtherMess       string    `json:"otherMess" form:"otherMess"`             //其他运行情况
}
