/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/28 15:27
 */

package request

import "server/module/sys_base/model"

type BusinessReportReq struct {
	model.BaseTimeRange
	Type    int    `json:"type" form:"type"`
	ProCode string `json:"proCode" form:"proCode"`
}
