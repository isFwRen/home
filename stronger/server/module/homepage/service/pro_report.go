/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/5 10:02
 */

package service

import (
	"server/global"
	"server/module/homepage/model"
	"server/module/homepage/model/request"
)

//GetProReport 获取项目日报
func GetProReport(proReportReq request.ProReportReq) (err error, list []model.ProReport, otherInfo model.ProReportOtherInfo) {
	global.GDb.Model(&model.ProReportOtherInfo{}).
		Where("to_char(report_date,'YYYY-MM-DD') = ?",
			proReportReq.ReportDay.Format("2006-01-02")).First(&otherInfo)
	err = global.GDb.Model(&model.ProReport{}).
		Where("to_char(report_date,'YYYY-MM-DD') = ?",
			proReportReq.ReportDay.Format("2006-01-02")).Find(&list).Error
	return err, list, otherInfo
}
