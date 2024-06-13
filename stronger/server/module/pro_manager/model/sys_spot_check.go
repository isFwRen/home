/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 3:05 下午
 */

package model

import (
	"server/module/sys_base/model"
	// "time"

	// "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type SysSpotCheck struct {
	model.Model
	CreatedCode string          `json:"createdCode" from:"createdCode"`
	CreatedName string          `json:"createdName" from:"createdName"`
	ProCode     string          `json:"proCode" from:"proCode"`
	Code        string          `json:"code" from:"code"`
	Name        string          `json:"name" from:"name"`
	Ratio       decimal.Decimal `json:"ratio" from:"ratio"`
	Type        int             `json:"type" from:"type"`
	Status      int             `json:"status" from:"status"` //图片加密秘钥
	//  DeadlineUploadTime time.Time      `json:"deadlineUploadTime" from:"deadlineUploadTime"`                           //图片加密秘钥
	//  Images             pq.StringArray `json:"images" form:"images" gorm:"type:varchar(100)[] comment:'一系列校验的数组'"`     // 原始图片

}
