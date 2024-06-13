/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/3 14:57
 */

package response

import "github.com/shopspring/decimal"

type ProData struct {
	ProCode          string          `json:"proCode"`          //项目编码
	AllCount         int64           `json:"allCount"`         //总业务量
	ReCount          int64           `json:"reCount"`          //返回业务量
	RePercent        decimal.Decimal `json:"rePercent"`        //返回率
	NotQualityCount  int64           `json:"notQualityCount"`  //待质检业务量
	QualityUserCount int64           `json:"qualityUserCount"` //质检人数
	NotInputCount    int64           `json:"notInputCount"`    //待录入业务量
	InputUserCount   int64           `json:"inputUserCount"`   //录入人数
}

type BusinessRanking struct {
	ProCode   string `json:"proCode"`
	BillCount int64  `json:"billCount"`
}

type AgingTrend struct {
	ProCode               string          `json:"proCode"`
	BillCount             int64           `json:"billCount"`
	TimeoutCount          int64           `json:"timeoutCount"`
	BillPercent           decimal.Decimal `json:"billPercent"`
	YesterdayBillCount    int64           `json:"yesterdayBillCount"`
	YesterdayTimeoutCount int64           `json:"yesterdayTimeoutCount"`
	YesterdayBillPercent  decimal.Decimal `json:"yesterdayBillPercent"`
}
