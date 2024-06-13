/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/15 10:56 上午
 */

package model

import "github.com/shopspring/decimal"

type BillAmount struct {
	BillAmount       decimal.Decimal `json:"billAmount"`       //账单金额总和
	DeductionAmount  decimal.Decimal `json:"deductionAmount"`  //扣费金额总和
	FeeAmount        decimal.Decimal `json:"feeAmount"`        //报销金额总和
	AdjustmentAmount decimal.Decimal `json:"adjustmentAmount"` //调整金额总和
}
