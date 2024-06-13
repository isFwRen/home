/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/5 09:57
 */

package request

import (
	"server/module/homepage/model"
	"time"
)

type ProReportReq struct {
	ReportDay time.Time `json:"reportDay" form:"reportDay"`
}

type ProReportOtherInfoReq struct {
	ProReportOtherInfo model.ProReportOtherInfo
	ProReport          []model.ProReport
}
