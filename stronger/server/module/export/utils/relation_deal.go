/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/24 11:58 上午
 */

package utils

import (
	"fmt"
	"go.uber.org/zap"
	"server/global"
	model2 "server/module/export/model"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
	"time"
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
//
//type ResultDataBill struct {
//	Bill    model.ProjectBill //单的信息
//	Invoice []InvoiceMap      //发票对象数组
//}
//
//type InvoiceMap struct {
//	Id               string                //发票属性
//	Code             string                //账单号
//	Money            float64               //账单金额
//	Invoice          [][]model3.ProjectField //发票
//	QingDan          [][]model3.ProjectField //清单
//	BaoXiaoDan       [][]model3.ProjectField //报销单
//	InvoiceDaXiang   [][]model3.ProjectField //发票大项
//	ThirdBaoXiaoDan1 [][]model3.ProjectField //第三方报销单1
//	ThirdBaoXiaoDan2 [][]model3.ProjectField //第三方报销单2
//	ThirdBaoXiaoDan3 [][]model3.ProjectField //第三方报销单3
//}

//RelationDeal 关系对应处理
func RelationDeal(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField, relationConf model2.RelationConf) model2.ResultDataBill {
	timeStart := time.Now()
	var resultDataBill model2.ResultDataBill
	resultDataBill.Bill = bill

	var invoiceMapList []model2.InvoiceMap
	//获取该单据的初审分块发票属性字段

	var otherMap model2.InvoiceMap
	otherMap.Invoice = make([][]model3.ProjectField, 1)
	otherMap.Id = "other"
	//for _, block := range blocks {
	//	//初审
	//	if block.Status == 1 {
	//		firstFields := fieldMap[block.ID]
	//		for _, field := range firstFields {
	//			//不属于发票(模板字段类型判断)Op0Value是数字5,6,7等，ResultValue已转bc003等
	//			if _, ok := relationConf.OtherTempType[field.Op0Value]; ok && field.Code == relationConf.TempTypeField {
	//				firstField := getOneFields(firstFields, field.BlockIndex)
	//				otherMap.Invoice[0] = append(otherMap.Invoice[0], firstField...)
	//				twoFields := getTwoFieldsBlockIndex(blocks, fieldMap, field.BlockIndex)
	//				otherMap.Invoice[0] = append(otherMap.Invoice[0], twoFields...)
	//			}
	//			//发票属性
	//			if field.Code == relationConf.InvoiceFieldCode.FieldCode {
	//				var invoiceMap model2.InvoiceMap
	//				invoiceMap.Id = field.ResultValue
	//				invoiceMap.Invoice = make([][]model3.ProjectField, 1)
	//				invoiceMap.Invoice[0] = packInvoiceMap(field.BlockIndex, firstFields, blocks, fieldMap, relationConf.InvoiceFieldCode.BlockCode)
	//				//获取账单号、总金额
	//				invoiceMap.Money, _ = strconv.ParseFloat(getFieldValue(relationConf.MoneyField, invoiceMap.Invoice[0]), 64)
	//				invoiceMap.Code = getFieldValue(relationConf.InvoiceNumField, invoiceMap.Invoice[0])
	//				invoiceMap.InvoiceType = getFieldValue(relationConf.InvoiceTypeField, invoiceMap.Invoice[0])
	//				for _, blockField := range firstFields {
	//					switch blockField.Code {
	//					case relationConf.QingDanFieldCode.FieldCode:
	//						//获取同一属性的清单所有字段
	//						if field.ResultValue == blockField.ResultValue {
	//							invoiceMap.QingDan = append(invoiceMap.QingDan, packInvoiceMap(blockField.BlockIndex, firstFields, blocks, fieldMap, relationConf.QingDanFieldCode.BlockCode))
	//							fmt.Println("11111")
	//						}
	//					case relationConf.BaoXiaoDanFieldCode.FieldCode:
	//						//获取同一属性的报销单所有字段
	//						if field.ResultValue == blockField.ResultValue {
	//							invoiceMap.BaoXiaoDan = append(invoiceMap.BaoXiaoDan, packInvoiceMap(blockField.BlockIndex, firstFields, blocks, fieldMap, relationConf.BaoXiaoDanFieldCode.BlockCode))
	//						}
	//					case relationConf.InvoiceDaXiang.FieldCode:
	//						//获取同一属性的发票大项所有字段
	//						if field.ResultValue == blockField.ResultValue {
	//							invoiceMap.InvoiceDaXiang = append(invoiceMap.InvoiceDaXiang, packInvoiceMap(blockField.BlockIndex, firstFields, blocks, fieldMap, relationConf.InvoiceDaXiang.BlockCode))
	//						}
	//					case relationConf.ThirdBaoXiaoDan1.FieldCode:
	//						//获取同一属性的第三报销单1所有字段
	//						if field.ResultValue == blockField.ResultValue {
	//							invoiceMap.ThirdBaoXiaoDan1 = append(invoiceMap.ThirdBaoXiaoDan1, packInvoiceMap(blockField.BlockIndex, firstFields, blocks, fieldMap, relationConf.ThirdBaoXiaoDan1.BlockCode))
	//						}
	//					case relationConf.ThirdBaoXiaoDan2.FieldCode:
	//						//获取同一属性的第三报销单2所有字段
	//						if field.ResultValue == blockField.ResultValue {
	//							invoiceMap.ThirdBaoXiaoDan2 = append(invoiceMap.ThirdBaoXiaoDan2, packInvoiceMap(blockField.BlockIndex, firstFields, blocks, fieldMap, relationConf.ThirdBaoXiaoDan2.BlockCode))
	//						}
	//					case relationConf.ThirdBaoXiaoDan3.FieldCode:
	//						//获取同一属性的第三报销单3所有字段
	//						if field.ResultValue == blockField.ResultValue {
	//							invoiceMap.ThirdBaoXiaoDan3 = append(invoiceMap.ThirdBaoXiaoDan3, packInvoiceMap(blockField.BlockIndex, firstFields, blocks, fieldMap, relationConf.ThirdBaoXiaoDan3.BlockCode))
	//						}
	//					}
	//
	//				}
	//				invoiceMapList = append(invoiceMapList, invoiceMap)
	//			}
	//		}
	//		break
	//	}
	//}

	invoiceMapList = append(invoiceMapList, otherMap)
	resultDataBill.Invoice = invoiceMapList
	global.GLog.Warn(bill.ProCode+" RelationDeal费时:::", zap.Any("毫秒", time.Since(timeStart).Milliseconds()))
	return resultDataBill
}

//packInvoiceMap 将同类型（发票||清单||报销单。。。）的所有字段等信息放到map对用的切片
func packInvoiceMap(blockIndex int, firstFields []model3.ProjectField, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField, blockCode string) []model3.ProjectField {
	var allFields = make([]model3.ProjectField, 0)
	//初审同一类型分块字段
	firstField := getOneFields(firstFields, blockIndex)
	allFields = append(allFields, firstField...)

	//非初审同一类型分块字段
	twoFields := getTwoFields(blocks, fieldMap, blockIndex, blockCode)
	allFields = append(allFields, twoFields...)

	fmt.Println(len(firstField), len(twoFields))
	return allFields
}

//getOneFields 获取初审同一属性的同一类型其他字段 （同一blockIndex）
func getOneFields(firstFields []model3.ProjectField, blockIndex int) []model3.ProjectField {
	var firstSameTypeFields = make([]model3.ProjectField, 0)
	for _, firstField := range firstFields {
		if firstField.BlockIndex == blockIndex {
			firstSameTypeFields = append(firstSameTypeFields, firstField)
		}
	}
	return firstSameTypeFields
}

//getTwoFields 获取非初审同一属性的同一类型其他字段（同一blockIndex 同一zero）
func getTwoFields(blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField, blockIndex int, blockCode string) []model3.ProjectField {
	var twoFields = make([]model3.ProjectField, 0)
	for _, block := range blocks {
		if block.Status == 2 && block.Zero == blockIndex &&
			((strings.Index(blockCode, "|") != -1 && utils.RegIsMatch(blockCode, block.Code)) ||
				strings.Index(blockCode, "|") == -1 && blockCode == block.Code) {
			twoFields = append(twoFields, fieldMap[block.ID]...)
		}
	}
	return twoFields
}

//getTwoFieldsBlockIndex 获取非初审同一类型其他字段（同一blockIndex 同一zero）
func getTwoFieldsBlockIndex(blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField, blockIndex int) []model3.ProjectField {
	var twoFields = make([]model3.ProjectField, 0)
	for _, block := range blocks {
		if block.Status == 2 && block.Zero == blockIndex {
			twoFields = append(twoFields, fieldMap[block.ID]...)
		}
	}
	return twoFields
}

//getFieldValue 获取字段的值
func getFieldValue(code string, fields []model3.ProjectField) (val string) {
	for _, field := range fields {
		if field.Code == code {
			return field.ResultValue
		}
	}
	return ""
}
