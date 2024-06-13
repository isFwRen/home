/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/24 11:58 上午
 */

package utils

import (
	"fmt"
	"reflect"
	"server/global"
	model2 "server/module/export/model"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
	"strconv"
	"time"

	"go.uber.org/zap"
)

//{
//	"单的信息": "",
//	"发票对象数组": [
//		{
//			"发票属性": "",
//			"账单号": "",
//			"发票类型": "",
//			"发票其他信息": "",
//			"qingdan": [],
//			"baoxiaodan": []
//		},
//		{
//			"发票属性": "",
//			"账单号": "",
//			"发票类型": "",
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
//	Type             string                //发票类型
//	Code             string                //账单号
//	Money            float64               //账单金额
//	Invoice          [][]model3.ProjectField //发票
//	QingDan          [][]model3.ProjectField //清单
//	BaoXiaoDan       [][]model3.ProjectField //报销单
//	InvoiceDaXiang   [][]model3.ProjectField //发票大项
//	ThirdBaoXiaoDan1 [][]model3.ProjectField //第三方报销单1
//	ThirdBaoXiaoDan2 [][]model3.ProjectField //第三方报销单2
//	ThirdBaoXiaoDan3 [][]model3.ProjectField //第三方报销单3
//Operation           [][]model3.ProjectField //手术
//HospitalizationDate [][]model3.ProjectField //住院日期
//HospitalizationFee  [][]model3.ProjectField //住院诊查费
//}

// RelationDealManyInvoiceType 关系对应处理
func RelationDealManyInvoiceType(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField, relationConf model2.RelationConf) model2.ResultDataBill {
	timeStart := time.Now()
	var resultDataBill model2.ResultDataBill
	resultDataBill.Bill = bill

	var invoiceMapList []model2.InvoiceMap
	//获取该单据的初审分块发票属性字段

	var otherMap model2.InvoiceMap
	otherMap.Id = "other"
	for _, block := range blocks {
		//初审
		if block.Status == 1 {
			firstFields := fieldMap[block.ID]
			for _, field := range firstFields {
				//不属于发票(模板字段类型判断)
				if _, ok := relationConf.OtherTempType[field.Op0Value]; ok && field.Code == relationConf.TempTypeField {
					var allFields = make([]model3.ProjectField, 0)
					firstField := getOneFields(firstFields, field.BlockIndex)
					allFields = append(allFields, firstField...)
					twoFields := getTwoFieldsBlockIndex(blocks, fieldMap, field.BlockIndex)
					fmt.Println("没有关系的类型", field.Name, "\t", field.Code, "\t", field.ResultValue, "\t", len(firstField), len(twoFields))
					allFields = append(allFields, twoFields...)
					otherMap.Invoice = append(otherMap.Invoice, allFields)
				}
				//发票属性
				for i, invoiceFieldCode := range relationConf.InvoiceFieldCode.FieldCode {
					if invoiceFieldCode == field.Code {
						fmt.Println("发票类型\t", field.Name, "\t", field.Code, "\t", field.ResultValue)
						var invoiceMap model2.InvoiceMap
						invoiceMap.Id = field.ResultValue
						invoiceMap.Invoice = make([][]model3.ProjectField, 1)
						invoiceMap.Invoice[0] = packInvoiceMap(field.BlockIndex, firstFields, blocks, fieldMap, relationConf.InvoiceFieldCode.BlockCode[i])
						//获取账单号、总金额
						invoiceMap.Money, _ = strconv.ParseFloat(getFieldValue(relationConf.MoneyField[i], invoiceMap.Invoice[0]), 64)
						invoiceMap.Code = getFieldValue(relationConf.InvoiceNumField[i], invoiceMap.Invoice[0])
						if relationConf.InvoiceTypeField != "" {
							invoiceMap.InvoiceType = getFieldValue(relationConf.InvoiceTypeField, invoiceMap.Invoice[0])
						} else {
							invoiceMap.InvoiceType = "门诊"
							if i == 1 {
								invoiceMap.InvoiceType = "住院"
							}
						}
						for _, blockField := range firstFields {
							//所有类型都要在这里列举
							typeObj := map[string]*[][]model3.ProjectField{"QingDan": &invoiceMap.QingDan,
								"BaoXiaoDan": &invoiceMap.BaoXiaoDan, "InvoiceDaXiang": &invoiceMap.InvoiceDaXiang,
								"ThirdBaoXiaoDan1": &invoiceMap.ThirdBaoXiaoDan1, "ThirdBaoXiaoDan2": &invoiceMap.ThirdBaoXiaoDan2,
								"ThirdBaoXiaoDan3": &invoiceMap.ThirdBaoXiaoDan3, "Operation": &invoiceMap.Operation,
								"HospitalizationDate": &invoiceMap.HospitalizationDate, "HospitalizationFee": &invoiceMap.HospitalizationFee,
								"ZhenDuan": &invoiceMap.ZhenDuan}
							for fieldName, typeVal := range typeObj {
								fieldConf := reflect.ValueOf(relationConf).FieldByName(fieldName + "FieldCode")
								if fieldConf.Kind().String() == "invalid" {
									//global.GLog.Warn("模板类型 【" + fieldName + " 】没有值")
									continue
								}
								fieldConfObj := fieldConf.Interface().(model2.TypeCode)
								if len(fieldConfObj.FieldCode) == 0 || len(fieldConfObj.BlockCode) == 0 {
									//global.GLog.Warn("模板类型 【" + fieldName + " 】没有值")
									continue
								}
								//global.GLog.Warn("模板类型 【" + fieldName + " 】")
								//处理只有一个相应配置字段的 情况例子：发票属性有两个，清单属性只有一个
								num := i
								if num > len(fieldConfObj.FieldCode)-1 {
									num = 0
								}

								//模板关系 field.ResultValue(发票属性) == blockField.ResultValue(所属发票属性)
								if fieldConfObj.FieldCode[num] == blockField.Code &&
									field.ResultValue == blockField.ResultValue {
									fmt.Println("有关系的类型\t", blockField.Name, "\t", blockField.Code, "\t", field.ResultValue)
									*typeVal = append(*typeVal, packInvoiceMap(blockField.BlockIndex, firstFields, blocks, fieldMap, fieldConfObj.BlockCode[0]))
								}
							}
							//switch blockField.Code {
							//case relationConf.QingDanFieldCode.FieldCode:
							//	//获取同一属性的清单所有字段
							//	if field.ResultValue == blockField.ResultValue {
							//		invoiceMap.QingDan = append(invoiceMap.QingDan, packInvoiceMap(blockField.BlockIndex, firstFields, blocks, fieldMap, relationConf.QingDanFieldCode.BlockCode))
							//		fmt.Println("11111")
							//	}
						}
						invoiceMapList = append(invoiceMapList, invoiceMap)
					}
				}
			}
			break
		}
	}

	invoiceMapList = append(invoiceMapList, otherMap)
	resultDataBill.Invoice = invoiceMapList
	global.GLog.Warn(bill.ProCode+" RelationDeal费时:::", zap.Any("毫秒", time.Since(timeStart).Milliseconds()))
	return resultDataBill
}

// GetBillType 返回单据类型
func GetBillType(obj model2.ResultDataBill) int {
	var billTypeName []string
	outpatient := 0
	inHospital := 0
	Meager := 0
	for _, invoiceMap := range obj.Invoice {
		if invoiceMap.InvoiceType != "" && ((invoiceMap.InvoiceType == "门诊" || invoiceMap.InvoiceType == "住院") || (invoiceMap.InvoiceType == "1" || invoiceMap.InvoiceType == "2")) {
			billTypeName = append(billTypeName, invoiceMap.InvoiceType)
		}
	}
	for _, item := range billTypeName {
		if item == "门诊" || item == "1" {
			outpatient = 1

		}
		if item == "住院" || item == "2" {
			inHospital = 2
		}
		if outpatient != 0 && inHospital != 0 {
			Meager = 3
		}
	}
	if Meager != 0 {
		obj.Bill.BillType = Meager
		return Meager
	} else if outpatient != 0 {
		obj.Bill.BillType = outpatient
		return outpatient
	} else if inHospital != 0 {
		obj.Bill.BillType = inHospital
		return inHospital
	}
	return 0
}
