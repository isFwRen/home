package B0121

import (
	"reflect"
	_func "server/global/func"
	model2 "server/module/export/model"
	"server/module/load/model"
	model3 "server/module/pro_manager/model"
)

type FormatObj struct {
	Bill    model3.ProjectBill   `json:"bill"`    //单的信息
	Fields  []model.ProjectField `json:"fields"`  //所有字段
	Invoice []FieldsMap          `json:"invoice"` //发票
}

type FieldsMap struct {
	Id      string               `json:"id"`    //发票属性
	Code    string               `json:"code"`  //账单号
	Money   float64              `json:"money"` //账单金额
	QingDan []FieldsMap          `json:"qingDan"`
	Fields  []model.ProjectField `json:"fields"` //发票
}

func FormatRenderObj(obj model2.ResultDataBill) (error, interface{}) {
	var formatObj FormatObj
	formatObj.Bill = obj.Bill

	for _, item := range obj.Invoice {
		//将所有字段丢到Fields
		eleLen := reflect.ValueOf(item).NumField()
		for j := 0; j < eleLen; j++ {
			if reflect.ValueOf(item).Field(j).Kind() != reflect.String && reflect.ValueOf(item).Field(j).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(item).Field(j).Interface().([][]model.ProjectField)
				for k := 0; k < len(fieldsArr); k++ {
					fields := fieldsArr[k]
					formatObj.Fields = append(formatObj.Fields, fields...)
				}
			}
		}

		if item.Id != "other" {
			var infoMap = FieldsMap{
				Id:     item.Id,
				Code:   item.Code,
				Money:  item.Money,
				Fields: append(item.Invoice[0]),
			}
			for _, qinDan := range item.QingDan {
				opFields := _func.FieldsFormat(qinDan)
				for _, fields := range opFields {
					infoMap.QingDan = append(infoMap.QingDan, FieldsMap{
						Fields: fields,
					})
				}

			}

			for _, fields := range item.BaoXiaoDan {
				infoMap.Fields = append(infoMap.Fields, fields...)
			}
			for _, fields := range item.ThirdBaoXiaoDan1 {
				infoMap.Fields = append(infoMap.Fields, fields...)
			}
			for _, fields := range item.ThirdBaoXiaoDan2 {
				infoMap.Fields = append(infoMap.Fields, fields...)
			}
			for _, fields := range item.ThirdBaoXiaoDan3 {
				infoMap.Fields = append(infoMap.Fields, fields...)
			}

			formatObj.Invoice = append(formatObj.Invoice, infoMap)
		}
	}

	if len(formatObj.Invoice) == 0 {
		formatObj.Invoice = append(formatObj.Invoice, FieldsMap{})
	}

	//global.GLog.Info("单号:::" + obj.Bill.BillNum + "	id:::" + obj.Bill.ID)
	return nil, formatObj
}
