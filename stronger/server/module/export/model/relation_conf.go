/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/26 10:07 上午
 */

package model

import (
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
)

//{
//	"单的信息": "",
//	"发票对象数组": [
//		{
//			"发票属性": "",
//			"账单号": "",
//			"发票其他信息": "",
//			"qingdan": [],
//			"baoxiaodan": []
//		},
//		{
//			"发票属性": "",
//			"账单号": "",
//			"发票其他信息": "",
//			"字段对象": []
//		}
//	]
//}

type ResultDataBill struct {
	Bill    model.ProjectBill `json:"bill"`    //单的信息
	Invoice []InvoiceMap      `json:"invoice"` //发票对象数组
}

type InvoiceMap struct {
	Id                  string                  `json:"id" msg:"发票属性"`                  //发票属性
	Code                string                  `json:"code" msg:"账单号"`                 //账单号
	Money               float64                 `json:"money" msg:"账单金额"`               //账单金额
	InvoiceType         string                  `json:"invoiceType" msg:"发票类型"`         //发票类型
	Invoice             [][]model3.ProjectField `json:"invoice" msg:"发票"`               //发票
	QingDan             [][]model3.ProjectField `json:"qingDan" msg:"清单"`               //清单
	BaoXiaoDan          [][]model3.ProjectField `json:"baoXiaoDan" msg:"报销单"`           //报销单
	InvoiceDaXiang      [][]model3.ProjectField `json:"invoiceDaXiang" msg:"发票大项"`      //发票大项
	ThirdBaoXiaoDan1    [][]model3.ProjectField `json:"thirdBaoXiaoDan1" msg:"第三方报销单1"` //第三方报销单1
	ThirdBaoXiaoDan2    [][]model3.ProjectField `json:"thirdBaoXiaoDan2" msg:"第三方报销单2"` //第三方报销单2
	ThirdBaoXiaoDan3    [][]model3.ProjectField `json:"thirdBaoXiaoDan3" msg:"第三方报销单3"` //第三方报销单3
	Operation           [][]model3.ProjectField `json:"operation" msg:"手术"`             //手术
	HospitalizationDate [][]model3.ProjectField `json:"hospitalizationDate" msg:"住院日期"` //住院日期
	HospitalizationFee  [][]model3.ProjectField `json:"hospitalizationFee" msg:"住院诊查费"` //住院诊查费
	ZhenDuan            [][]model3.ProjectField `json:"zhenDuan" msg:"诊断"`              //诊断
}

////////发票所拥有的的关系////////

type RelationConf struct {
	InvoiceFieldCode             TypeCode          `json:"invoiceFieldCode" msg:"发票"`               //发票
	QingDanFieldCode             TypeCode          `json:"qingDanFieldCode" msg:"清单"`               //清单
	BaoXiaoDanFieldCode          TypeCode          `json:"baoXiaoDanFieldCode" msg:"报销单"`           //报销单
	InvoiceDaXiangFieldCode      TypeCode          `json:"invoiceDaXiang" msg:"发票大项"`               //发票大项
	ThirdBaoXiaoDan1FieldCode    TypeCode          `json:"thirdBaoXiaoDan1" msg:"第三方报销单1"`          //第三方报销单1
	ThirdBaoXiaoDan2FieldCode    TypeCode          `json:"thirdBaoXiaoDan2" msg:"第三方报销单2"`          //第三方报销单2
	ThirdBaoXiaoDan3FieldCode    TypeCode          `json:"thirdBaoXiaoDan3" msg:"第三方报销单3"`          //第三方报销单3
	OperationFieldCode           TypeCode          `json:"operationFieldCode" msg:"手术"`             //手术
	HospitalizationDateFieldCode TypeCode          `json:"hospitalizationDateFieldCode" msg:"住院日期"` //住院日期
	HospitalizationFeeFieldCode  TypeCode          `json:"hospitalizationFeeFieldCode" msg:"住院诊查费"` //住院诊查费
	ZhenDuanFieldCode            TypeCode          `json:"zhenDuanFieldCode" msg:"诊断"`              //诊断
	OtherTempType                map[string]string `json:"otherTempType" msg:"其他模板类型"`              //其他模板类型
	TempTypeField                string            `json:"tempTypeField" msg:"模板类型字段"`              //模板类型字段
	InvoiceNumField              []string          `json:"invoiceNumField" msg:"账单号"`               //账单号
	MoneyField                   []string          `json:"moneyField" msg:"总金额"`                    //总金额
	InvoiceTypeField             string            `json:"invoiceTypeField" msg:"发票类型字段"`           //发票类型字段
}

type TypeCode struct {
	FieldCode []string `json:"fieldCode"` //关系字段
	BlockCode []string `json:"blockCode"` //非初审分块(支持正则,正则必须包含|，否则不能包含|)
}
