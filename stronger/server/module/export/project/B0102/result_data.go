package B0102

import (
	"fmt"
	"reflect"
	model2 "server/module/export/model"
	eUtils "server/module/export/utils"
	model3 "server/module/load/model"
	"server/module/pro_manager/model"
)

// RelationConf 关系对应
var relationConf = model2.RelationConf{
	InvoiceFieldCode: model2.TypeCode{FieldCode: []string{"fc154"}, BlockCode: []string{"bc027"}},
	// QingDanFieldCode:          model2.TypeCode{FieldCode: []string{"fc158"}, BlockCode: []string{"bc019"}},
	BaoXiaoDanFieldCode: model2.TypeCode{FieldCode: []string{"fc156"}, BlockCode: []string{"bc028"}},
	// ThirdBaoXiaoDan1FieldCode: model2.TypeCode{FieldCode: []string{"fc160"}, BlockCode: []string{"bc021"}},
	OtherTempType:    map[string]string{"10": "10", "11": "11", "12": "12", "13": "13", "1": "1", "3": "3", "16": "16", "4": "4", "6": "6", "19": "19", "91": "91", "92": "92", "93": "93", "94": "94", "95": "95", "96": "96"},
	TempTypeField:    "fc003",
	InvoiceNumField:  []string{"fc042", ""},
	MoneyField:       []string{"fc046", ""},
	InvoiceTypeField: "fc326",
	//InvoiceTypeField: "fc003",
}

// ResultData B0108
func ResultData(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField) (err error, obj model2.ResultDataBill) {
	obj = eUtils.RelationDealManyInvoiceType(bill, blocks, fieldMap, relationConf)
	//常量
	// constMap := constDeal(bill.ProCode)

	fieldLocationMap := make(map[string][][]int)
	n := 0
	for i := 0; i < len(obj.Invoice); i++ {
		//同一发票
		invoiceMap := obj.Invoice[i]
		eleLen := reflect.ValueOf(invoiceMap).NumField()
		for j := 0; j < eleLen; j++ {
			if reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.String && reflect.ValueOf(invoiceMap).Field(j).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(invoiceMap).Field(j).Interface().([][]model3.ProjectField)
				for k := 0; k < len(fieldsArr); k++ {
					fields := fieldsArr[k]
					for l := 0; l < len(fields); l++ {
						//fields[l].Issues = nil
						//i:发票index j:发票结构体字段index k:字段二维数组x的index l:字段二维数组y的index
						fieldLocationMap[fields[l].Code] = append(fieldLocationMap[fields[l].Code], []int{i, j, k, l, fields[l].BlockIndex})
						n++
					}
				}
			}
		}
	}
	fmt.Println(n)

	return nil, obj
}

func GetFieldValue(fields []model3.ProjectField, code string, isFinal bool, bidx int) string {
	for _, field := range fields {
		if field.Code == code && (field.BlockIndex == bidx || bidx == -1) {
			if isFinal {
				return field.FinalValue
			} else {
				return field.ResultValue
			}
		}
	}
	return ""
}
