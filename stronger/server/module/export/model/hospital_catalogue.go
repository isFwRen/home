/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/15 14:51
 */

package model

import (
	"server/module/sys_base/model"
	"time"
)

type HospitalCatalogue struct {
	model.Model
	BillId  string    `json:"billId" form:"billId"`
	BillNum string    `json:"billNum" form:"billNum"  excel:"账单号"`
	Type    int       `json:"type" form:"type"`
	Agency  string    `json:"agency" form:"agency"  excel:"机构号"`
	Name    string    `json:"name" form:"name"  excel:"名称"`
	Num     string    `json:"num" form:"num"  excel:"数量"`
	Date    time.Time `json:"date" form:"date"  excel:"日期"`
	Money   string    `json:"money" form:"money"  excel:"金额"`
}

type Agency struct {
	model.Model
	BillId         string `json:"billId" form:"billId"`
	BillNum        string `json:"billNum" form:"billNum" excel:"账单号"`
	MIcd10Code     string `json:"mIcd10Code" form:"mIcd10Code" excel:"疾病"`
	Agency         string `json:"agency" form:"agency" excel:"机构"`
	IsMatch        string `json:"isMatch" form:"isMatch" excel:"是否匹配"`
	ExpenMode      string `json:"expenMode" form:"expenMode" excel:"票据类型"`
	Hospital       string `json:"hospital" form:"hospital" excel:"医院名称"`
	CountMoney     string `json:"countMoney" form:"countMoney" excel:"总金额"`
	SocialPayMoney string `json:"socialPayMoney" form:"socialPayMoney" excel:"统筹金额"`
	OutMoney       string `json:"outMoney" form:"outMoney" excel:"范围外金额"`
	InnerMoney     string `json:"innerMoney" form:"innerMoney" excel:"范围内金额"`
}
