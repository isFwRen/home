/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 3:05 下午
 */

package model

import (
	"server/module/sys_base/model"
	"time"

	// "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type SysSpotCheckStatistic struct {
	model.Model
	ProCode    string          `json:"proCode" from:"proCode"`
	Sum        int64           `json:"sum" from:"sum"`
	CheckSum   int             `json:"checkSum" from:"checkSum"`
	WrongNum   int             `json:"wrongNum" from:"wrongNum"` //图片加密秘钥
	WrongRatio decimal.Decimal `json:"wrongRatio" from:"wrongRatio"`
	Ratio      decimal.Decimal `json:"ratio" from:"ratio"`
	SubmitDay  time.Time       `json:"submitDay" from:"submitDay"` //图片加密秘钥
	//  Images             pq.StringArray `json:"images" form:"images" gorm:"type:varchar(100)[] comment:'一系列校验的数组'"`     // 原始图片

}
