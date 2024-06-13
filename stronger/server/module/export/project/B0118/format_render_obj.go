/**
 * @Author: xingqiyi
 * @Description:格式化渲染xml的对象
 * @Date: 2022/2/18 9:40 上午
 */

package B0118

import (
	"reflect"
	"server/global"
	model2 "server/module/export/model"
	"server/module/load/model"
	model3 "server/module/pro_manager/model"

	"go.uber.org/zap"
)

type FormatObj struct {
	Bill         model3.ProjectBill   `json:"bill"`         //单的信息
	Fields       []model.ProjectField `json:"fields"`       //所有字段
	HospitalList []FieldsMap          `json:"hospitalList"` //住院发票对象数组
	OpList       []FieldsMap          `json:"opList"`       //门诊发票对象数组
	SiList       []FieldsMap          `json:"siList"`       //报销单
	ThirdList    []FieldsMap          `json:"thirdList"`    //第三方报销单1
	ThirdList2   []FieldsMap          `json:"thirdList2"`   //第三方报销单2,3
}

type FieldsMap struct {
	Id      string               `json:"id"`      //发票属性
	Code    string               `json:"code"`    //账单号
	Money   float64              `json:"money"`   //账单金额
	Fields  []model.ProjectField `json:"fields"`  //发票,大项和诊断等
	QingDan []FieldsMap          `json:"qingDan"` //清单
}

func FormatRenderObj(obj model2.ResultDataBill) (error, interface{}) {

	var formatObj FormatObj
	formatObj.Bill = obj.Bill

	//获取不具有所属发票属性的字段，包括诊断
	var otherFields []model.ProjectField
	for _, item := range obj.Invoice {
		if item.Id == "other" {
			//otherFields = append(otherFields, item.Invoice[0]...)
			for _, fields := range item.Invoice {
				otherFields = append(otherFields, fields...)
			}
		}
	}
	//将所有字段丢到Fields
	formatObj.Fields = append(formatObj.Fields, otherFields...)

	global.GLog.Info("单号:::" + obj.Bill.BillNum + "	id:::" + obj.Bill.ID)
	//构造新渲染xml的对象
	for _, item := range obj.Invoice {
		fc205 := ""
		if item.Id != "other" {
			//将所有字段丢到Fields
			eleLen := reflect.ValueOf(item).NumField()
			for j := 0; j < eleLen; j++ {
				if reflect.ValueOf(item).Field(j).Kind() != reflect.String && reflect.ValueOf(item).Field(j).Kind() != reflect.Float64 {
					//每张发票每种类型的字段
					fieldsArr := reflect.ValueOf(item).Field(j).Interface().([][]model.ProjectField)
					for k := 0; k < len(fieldsArr); k++ {
						fields := fieldsArr[k]
						for l := 0; l < len(fields); l++ {
							formatObj.Fields = append(formatObj.Fields, fields[l])
							if fields[l].Code == "fc205" {
								fc205 = fields[l].FinalValue
							}
						}
					}
				}
			}

			//发票
			var invoiceDaXiangFields []model.ProjectField
			for _, fields := range item.InvoiceDaXiang {
				invoiceDaXiangFields = append(invoiceDaXiangFields, fields...)
			}

			var invoiceHospitalMap = FieldsMap{
				Id:     item.Id,
				Code:   item.Code,
				Money:  item.Money,
				Fields: append(append(item.Invoice[0], invoiceDaXiangFields...), otherFields...),
			}
			for _, ele := range item.BaoXiaoDan {
				invoiceHospitalMap.Fields = append(invoiceHospitalMap.Fields, ele...)
			}

			global.GLog.Info(item.Id, zap.Any("发票字段数量", len(item.Invoice)))
			global.GLog.Info(item.Id, zap.Any("发票大项字段数量", len(item.InvoiceDaXiang)))
			global.GLog.Info(item.Id, zap.Any("其他无关系字段数量", len(otherFields)))
			global.GLog.Info(item.Id, zap.Any("发票清单字段数量", len(item.QingDan)))
			if len(invoiceHospitalMap.Fields) > 0 {
				var qinDanFieldsMap []FieldsMap
				for _, qinDan := range item.QingDan {
					qinDanFieldsMap = append(qinDanFieldsMap, FieldsMap{
						Fields: qinDan,
					})
				}
				invoiceHospitalMap.QingDan = append(invoiceHospitalMap.QingDan, qinDanFieldsMap...)
			}
			global.GLog.Info(item.Id, zap.Any("发票类型fc003值只能是(1,2)", item.InvoiceType))
			//要具体需求发票里面的fc003区分
			if item.InvoiceType == "2" {
				//住院发票对象数组
				formatObj.HospitalList = append(formatObj.HospitalList, invoiceHospitalMap)
			} else if item.InvoiceType == "1" {
				//门诊发票对象数组
				formatObj.OpList = append(formatObj.OpList, invoiceHospitalMap)
			}

			//报销单
			if len(item.Invoice) > 0 {
				global.GLog.Info("----------SiList--fc205----------" + fc205 + "   " + item.InvoiceType)
				if !((fc205 != "0" && fc205 != "1" && item.InvoiceType != "") || fc205 == "") {
					var baoXiaoDanFields []model.ProjectField
					for _, fields := range item.Invoice {
						baoXiaoDanFields = append(baoXiaoDanFields, fields...)
					}
					baoXiaoDan := FieldsMap{
						Fields: baoXiaoDanFields,
					}
					global.GLog.Info(item.Id, zap.Any("报销单字段数量", len(baoXiaoDan.Fields)))
					//将属于发票分块的一些字段也放到这里
					baoXiaoDan.Fields = append(baoXiaoDan.Fields, item.Invoice[0]...)
					formatObj.SiList = append(formatObj.SiList, baoXiaoDan)
				}

			}

			//第三方报销单
			// 当fc205结果值为0、1且fc003不为空时，屏蔽thirdList节点； 当fc205结果为空时，屏蔽thirdList、siList节点
			global.GLog.Info("------ThirdList----fc205----------" + fc205 + "   " + item.InvoiceType)
			if !(((fc205 == "0" || fc205 == "1") && item.InvoiceType != "") || fc205 == "") {
				thirdBaoXiaoDan1 := FieldsMap{
					Fields: item.Invoice[0],
				}
				for _, fields := range item.BaoXiaoDan {
					thirdBaoXiaoDan1.Fields = append(thirdBaoXiaoDan1.Fields, fields...)
				}
				for _, fields := range item.ThirdBaoXiaoDan1 {
					thirdBaoXiaoDan1.Fields = append(thirdBaoXiaoDan1.Fields, fields...)
				}
				global.GLog.Info(item.Id, zap.Any("第三方报销单1字段数量", len(thirdBaoXiaoDan1.Fields)))
				//将属于发票分块的一些字段也放到这里
				//thirdBaoXiaoDan1.Fields = append(thirdBaoXiaoDan1.Fields, item.Invoice...)
				formatObj.ThirdList = append(formatObj.ThirdList, thirdBaoXiaoDan1)
			}

			//}

			//第三方报销单1，2,3
			if len(item.ThirdBaoXiaoDan1) > 0 {
				thirdBaoXiaoDan1 := FieldsMap{
					Fields: item.ThirdBaoXiaoDan1[0],
				}
				thirdBaoXiaoDan1.Fields = append(thirdBaoXiaoDan1.Fields, item.Invoice[0]...)
				if len(item.ThirdBaoXiaoDan2) > 0 {
					thirdBaoXiaoDan1.Fields = append(thirdBaoXiaoDan1.Fields, item.ThirdBaoXiaoDan2[0]...)
				}
				if len(item.ThirdBaoXiaoDan3) > 0 {
					thirdBaoXiaoDan1.Fields = append(thirdBaoXiaoDan1.Fields, item.ThirdBaoXiaoDan3[0]...)
				}
				global.GLog.Info(item.Id, zap.Any("第三方报销单2字段数量", len(thirdBaoXiaoDan1.Fields)))

				if len(thirdBaoXiaoDan1.Fields) > 0 {
					formatObj.ThirdList2 = append(formatObj.ThirdList2, thirdBaoXiaoDan1)
				}
			}
			//if len(item.ThirdBaoXiaoDan2) > 0 {
			//	thirdBaoXiaoDan2 := FieldsMap{
			//		Fields: item.ThirdBaoXiaoDan2[0],
			//	}
			//	thirdBaoXiaoDan2.Fields = append(thirdBaoXiaoDan2.Fields, item.Invoice[0]...)
			//	global.GLog.Info(item.Id, zap.Any("第三方报销单2字段数量", len(thirdBaoXiaoDan2.Fields)))
			//
			//	if len(thirdBaoXiaoDan2.Fields) > 0 {
			//		formatObj.ThirdList2 = append(formatObj.ThirdList2, thirdBaoXiaoDan2)
			//	}
			//}
			//if len(item.ThirdBaoXiaoDan3) > 0 {
			//	thirdBaoXiaoDan3 := FieldsMap{
			//		Fields: item.ThirdBaoXiaoDan3[0],
			//	}
			//	global.GLog.Info(item.Id, zap.Any("第三方报销单3字段数量", len(thirdBaoXiaoDan3.Fields)))
			//	if len(thirdBaoXiaoDan3.Fields) > 0 {
			//		formatObj.ThirdList2 = append(formatObj.ThirdList2, thirdBaoXiaoDan3)
			//	}
			//}
			global.GLog.Info(item.Id, zap.Any("-----------------------", "-----------------------"))
		}
	}
	return nil, formatObj
}
