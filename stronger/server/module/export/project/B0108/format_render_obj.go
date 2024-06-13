package B0108

import (
	"reflect"
	_func "server/global/func"
	model2 "server/module/export/model"
	"server/module/load/model"
	model3 "server/module/pro_manager/model"
)

type FormatObj struct {
	Bill          model3.ProjectBill   `json:"bill"`          //单的信息
	Fields        []model.ProjectField `json:"fields"`        //所有字段
	ClinicInfo    []FieldsMap          `json:"clinicInfo"`    //门诊发票
	InpatientInfo []FieldsMap          `json:"inpatientInfo"` //住院发票
}

type FieldsMap struct {
	Id                  string               `json:"id"`                             //发票属性
	Code                string               `json:"code"`                           //账单号
	Money               float64              `json:"money"`                          //账单金额
	Fields              []model.ProjectField `json:"fields"`                         //发票
	QingDan             []FieldsMap          `json:"qingDan"`                        //清单
	BaoXiaoDan          []model.ProjectField `json:"BaoXiaoDan"`                     //报销单
	HospitalizationDate []model.ProjectField `json:"hospitalizationDate" msg:"住院日期"` //住院日期

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
			for _, fields := range item.BaoXiaoDan {
				infoMap.BaoXiaoDan = append(infoMap.BaoXiaoDan, fields...)
			}
			for _, fields := range item.HospitalizationDate {
				infoMap.HospitalizationDate = append(infoMap.HospitalizationDate, fields...)
			}
			for _, fields := range item.Operation {
				infoMap.Fields = append(infoMap.Fields, fields...)
			}
			for _, fields := range item.HospitalizationFee {
				infoMap.Fields = append(infoMap.Fields, fields...)
			}
			// var qinDanFieldsMap []FieldsMap
			for _, qinDan := range item.QingDan {
				opFields := _func.FieldsFormat(qinDan)
				for _, fields := range opFields {
					infoMap.QingDan = append(infoMap.QingDan, FieldsMap{
						Fields: fields,
					})
					// qinDanFieldsMap = append(qinDanFieldsMap, FieldsMap{
					// 	Fields: fields,
					// })
				}

			}
			// infoMap.QingDan = append(infoMap.QingDan, qinDanFieldsMap...)
			if item.InvoiceType == "门诊" {
				fc032 := ""
				for _, field := range infoMap.Fields {
					if field.Code == "fc032" {
						fc032 = field.FinalValue
						break
					}
				}
				if fc032 != "" {
					formatObj.ClinicInfo = append(formatObj.ClinicInfo, infoMap)
				}
			}
			if item.InvoiceType == "住院" {
				fc054 := ""
				for _, field := range infoMap.Fields {
					if field.Code == "fc054" {
						fc054 = field.FinalValue
						break
					}
				}
				if fc054 != "" {
					formatObj.InpatientInfo = append(formatObj.InpatientInfo, infoMap)
				}
			}
		}
	}

	if len(formatObj.ClinicInfo) == 0 {
		formatObj.ClinicInfo = append(formatObj.ClinicInfo, FieldsMap{})
	}
	for qq, ClinicInfo := range formatObj.ClinicInfo {
		if len(ClinicInfo.QingDan) == 0 {
			formatObj.ClinicInfo[qq].QingDan = append(formatObj.ClinicInfo[qq].QingDan, FieldsMap{})
		}
	}
	if len(formatObj.InpatientInfo) == 0 {
		formatObj.InpatientInfo = append(formatObj.InpatientInfo, FieldsMap{})
	}
	for qq, InpatientInfo := range formatObj.InpatientInfo {
		if len(InpatientInfo.QingDan) == 0 {
			formatObj.InpatientInfo[qq].QingDan = append(formatObj.InpatientInfo[qq].QingDan, FieldsMap{})
		}
	}
	//global.GLog.Info("单号:::" + obj.Bill.BillNum + "	id:::" + obj.Bill.ID)

	return nil, formatObj
}
