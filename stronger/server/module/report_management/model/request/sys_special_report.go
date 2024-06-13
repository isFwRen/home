/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/16 14:18
 */

package request

import (
	"server/module/sys_base/model"
)

type NewHospitalAndCatalogueSearch struct {
	model.BaseTimePageCode
	Type int `json:"type" form:"type,default=-1"` //状态
}

type NewHospitalAndCatalogueExportSearch struct {
	model.BaseTimeRangeWithCode
	Type int `json:"type" form:"type,default=-1"` //状态
}
