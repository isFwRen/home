/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 3:05 下午
 */

package model

import (
	"server/module/sys_base/model"
	"time"

	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	// "time"
)

type SysSpotCheckData struct {
	model.Model
	// CreatedCode string `json:"createdCode" from:"createdCode"`
	// CreatedName string `json:"createdName" from:"createdName"`
	ProCode   string          `json:"pro" from:"proCode"`
	Code      string          `json:"code" from:"code"`
	Name      string          `json:"name" from:"name"`
	Num       int             `json:"num" from:"num"`
	DoneNum   int             `json:"doneNum" from:"doneNum"`
	UndoneNum int             `json:"undoneNum" from:"undoneNum"` //图片加密秘钥
	WrongNum  int             `json:"wrongNum" from:"wrongNum"`
	Type      int             `json:"type" from:"type"`
	Ratio     decimal.Decimal `json:"ratio" from:"ratio"`
	SubmitDay time.Time       `json:"submitDay" from:"submitDay"` //图片加密秘钥
	BlockId   pq.StringArray  `json:"blockId" gorm:"type:varchar(32)[];" `
}
