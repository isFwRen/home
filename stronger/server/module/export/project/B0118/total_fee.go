/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/15 10:50 上午
 */

package B0118

import (
	model4 "server/module/pro_conf/model"
	"server/module/pro_manager/model"
	"server/utils"

	"github.com/shopspring/decimal"
)

func TotalFee(formatQuality map[int][]model4.SysQuality, xml string) (err error, formatQuality2 map[int][]model4.SysQuality, obj model.BillAmount) {

	//账单金额总和thirdInSIAmount
	//扣费金额总和paymentFee
	//报销金额总和fee
	thirdInSIAmount := utils.GetNodeData(xml, "thirdInSIAmount")
	paymentFee := utils.GetNodeData(xml, "paymentFee")
	fee := utils.GetNodeData(xml, "fee")

	obj = model.BillAmount{
		BillAmount:       totalStr(thirdInSIAmount),
		FeeAmount:        totalStr(fee),
		DeductionAmount:  totalStr(paymentFee),
		AdjustmentAmount: decimal.NewFromFloat(0),
	}

	return nil, formatQuality, obj
}

func totalStr(str []string) decimal.Decimal {
	total := decimal.NewFromFloat(0)
	for _, s := range str {
		fromString, err := decimal.NewFromString(s)
		if err != nil {
			return total
		}
		total = decimal.Sum(total, fromString)
	}
	return total
}
