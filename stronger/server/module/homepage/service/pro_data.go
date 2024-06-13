/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/3 14:56
 */

package service

import (
	"github.com/shopspring/decimal"
	"server/global"
	"server/module/homepage/model/request"
	"server/module/homepage/model/response"
	"server/module/pro_manager/model"
	"time"
)

var dateFormat = map[int][]string{
	0: {"YYYY-MM-DD", "2006-01-02"}, //日
	1: {"YYYY-MM", "2006-01"},       //月
	2: {"YYYY", "2006"},             //年
}

//GetProData 获取项目数据待处理数据
func GetProData(queryDay request.QueryDayReq) (err error, list []response.ProData) {
	str := `SELECT 
			a.count as all_count,
			b.count as re_count,
			c.count as not_quality_count,
			d.count as quality_user_count,
			e.count as not_input_count,
			f.count as input_user_count
		FROM 
		(SELECT count(1) from project_bills WHERE to_char(scan_at, 'YYYY-MM-DD') = ?) a,
		(SELECT count(1) from project_bills WHERE to_char(scan_at, 'YYYY-MM-DD') = ? and "stage" = 5) b,
		(SELECT count(1) from project_bills WHERE to_char(scan_at, 'YYYY-MM-DD') = ? and (stage = 3 or stage = 4) ) c,
		(SELECT count(DISTINCT quality_user_code) FROM project_bills where to_char(scan_at, 'YYYY-MM-DD') = ? and quality_user_code<>'') d,
		(SELECT count(1) FROM project_bills where to_char(scan_at, 'YYYY-MM-DD') = ? and (stage = 1 or stage = 2)) e,
		(SELECT count(1) from (
			SELECT DISTINCT op1_code FROM project_blocks where to_char(created_at, 'YYYY-MM-DD') = ?
			UNION
			SELECT DISTINCT op0_code FROM project_blocks where to_char(created_at, 'YYYY-MM-DD') = ?
			UNION
			SELECT DISTINCT op2_code FROM project_blocks where to_char(created_at, 'YYYY-MM-DD') = ?
			UNION
			SELECT DISTINCT opq_code FROM project_blocks where to_char(created_at, 'YYYY-MM-DD') = ?
		)g)f
		`
	list = make([]response.ProData, 0)
	for proCode, _ := range global.GProConf {
		db := global.ProDbMap[proCode]
		if db != nil {
			var obj response.ProData
			err = db.Raw(str, queryDay.QueryDay, queryDay.QueryDay, queryDay.QueryDay, queryDay.QueryDay, queryDay.QueryDay, queryDay.QueryDay, queryDay.QueryDay, queryDay.QueryDay, queryDay.QueryDay).
				Scan(&obj).Error
			if err != nil {
				return err, list
			}
			obj.ProCode = proCode
			if obj.AllCount == 0 {
				obj.RePercent = decimal.Zero
			} else {
				obj.RePercent = decimal.NewFromInt(obj.ReCount).DivRound(decimal.NewFromInt(obj.AllCount), 4)
			}
			list = append(list, obj)
		}
	}
	return err, list
}

//GetBusinessRanking 获取项目业务量趋势
func GetBusinessRanking(queryBusinessReq request.QueryBusinessReq) (err error, list []response.BusinessRanking) {
	timeStr := time.Now().Format(dateFormat[queryBusinessReq.RankingType][1])
	list = make([]response.BusinessRanking, 0)
	for proCode, _ := range global.GProConf {
		db := global.ProDbMap[proCode]
		var proTotal int64
		if db != nil {
			err = db.Model(&model.ProjectBill{}).
				Where("to_char(scan_at, '"+dateFormat[queryBusinessReq.RankingType][0]+"') = ?", timeStr).
				Count(&proTotal).Error
			if err != nil {
				return err, list
			}
			businessRanking := response.BusinessRanking{
				ProCode:   proCode,
				BillCount: proTotal,
			}
			list = append(list, businessRanking)
		}
	}
	return err, list
}

//GetAgingTrend 获取项目时效趋势
func GetAgingTrend(queryBusinessReq request.QueryBusinessReq) (err error, list []response.AgingTrend) {
	timeStr := time.Now().Format(dateFormat[queryBusinessReq.RankingType][1])
	timeYesterdayStr := time.Now().Add(-24 * 60 * time.Minute).Format(dateFormat[queryBusinessReq.RankingType][1])
	if queryBusinessReq.RankingType == 1 {
		timeYesterdayStr = time.Now().Add(-24 * 60 * time.Duration(time.Now().Day()) * time.Minute).Format(dateFormat[queryBusinessReq.RankingType][1])
	}
	list = make([]response.AgingTrend, 0)
	for proCode, _ := range global.GProConf {
		db := global.ProDbMap[proCode]
		var proTotal, proTimeout, proYesterdayTotal, proYesterdayTimeout int64
		if db != nil {
			err = db.Model(&model.ProjectBill{}).
				Where("to_char(scan_at, '"+dateFormat[queryBusinessReq.RankingType][0]+"') = ?", timeStr).
				Count(&proTotal).Error
			if err != nil {
				return err, list
			}
			err = db.Model(&model.ProjectBill{}).
				Where("is_timeout = true and to_char(scan_at, '"+dateFormat[queryBusinessReq.RankingType][0]+"') = ?", timeStr).
				Count(&proTimeout).Error
			if err != nil {
				return err, list
			}

			err = db.Model(&model.ProjectBill{}).
				Where("to_char(scan_at, '"+dateFormat[queryBusinessReq.RankingType][0]+"') = ?", timeYesterdayStr).
				Count(&proYesterdayTotal).Error
			if err != nil {
				return err, list
			}
			err = db.Model(&model.ProjectBill{}).
				Where("is_timeout = true and to_char(scan_at, '"+dateFormat[queryBusinessReq.RankingType][0]+"') = ?", timeYesterdayStr).
				Count(&proYesterdayTimeout).Error
			if err != nil {
				return err, list
			}

			agingTrend := response.AgingTrend{
				ProCode:               proCode,
				BillCount:             proTotal,
				TimeoutCount:          proTimeout,
				YesterdayBillCount:    proYesterdayTotal,
				YesterdayTimeoutCount: proYesterdayTimeout,
			}

			if proTotal == 0 {
				agingTrend.BillPercent = decimal.Zero
			} else {
				agingTrend.BillPercent = decimal.NewFromInt(agingTrend.TimeoutCount).DivRound(decimal.NewFromInt(agingTrend.BillCount), 4)
			}

			if proYesterdayTotal == 0 {
				agingTrend.YesterdayBillPercent = decimal.Zero
			} else {
				agingTrend.YesterdayBillPercent = decimal.NewFromInt(agingTrend.YesterdayTimeoutCount).DivRound(decimal.NewFromInt(agingTrend.YesterdayBillCount), 4)
			}
			list = append(list, agingTrend)
		}
	}
	return err, list
}
