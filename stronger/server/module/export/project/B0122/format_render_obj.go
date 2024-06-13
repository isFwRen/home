package B0122

import (
	"fmt"
	"reflect"
	"regexp"
	_func "server/global/func"
	model2 "server/module/export/model"
	"server/module/load/model"
	model3 "server/module/pro_manager/model"
	"server/utils"
	"strings"
	"time"
)

type FormatObj struct {
	Bill    model3.ProjectBill   `json:"bill"`   //单的信息
	Fields  []model.ProjectField `json:"fields"` //所有字段
	Medical []YiYuanObj          `json:"medical"`
}

type YiYuanObj struct {
	Name        string               `json:"name"`        //医院名称
	Fields      []model.ProjectField `json:"fields"`      //所有字段
	Hospital    []FieldsMap          `json:"hospital"`    //住院发票
	Clinc       []FieldsMap          `json:"clinc"`       //门诊发票
	Pharmacy    []FieldsMap          `json:"pharmacy"`    //药房发票
	Nonmedical  []FieldsMap          `json:"nonmedical"`  //非医疗发票
	Injury      []FieldsMap          `json:"injury"`      //诊断
	Injurybj    []FieldsMap          `json:"injurybj"`    //诊断北京
	Operation   []FieldsMap          `json:"operation"`   //手术
	Operationbj []FieldsMap          `json:"operationbj"` //手术北京
}

type FieldsMap struct {
	Id         string               `json:"id"`         //发票属性
	Code       string               `json:"code"`       //账单号
	Money      float64              `json:"money"`      //账单金额
	Fields     []model.ProjectField `json:"fields"`     //发票
	QingDan    []FieldsMap          `json:"qingDan"`    //清单
	BaoXiaoDan []model.ProjectField `json:"BaoXiaoDan"` //报销单

}

func FormatRenderObj(obj model2.ResultDataBill) (error, interface{}) {
	var formatObj FormatObj
	formatObj.Bill = obj.Bill

	isNotFaPiao := true
	for _, item := range obj.Invoice {
		for _, fields := range item.Invoice {
			for _, field := range fields {
				if field.Code == "fc153" && utils.RegIsMatch(`^(5|6)$`, field.ResultValue) {
					isNotFaPiao = false
				}
			}
		}
	}
	fc040s := []string{}
	for _, item := range obj.Invoice {
		//将所有字段丢到Fields
		var yiYuanFields []model.ProjectField
		eleLen := reflect.ValueOf(item).NumField()
		for j := 0; j < eleLen; j++ {
			if reflect.ValueOf(item).Field(j).Kind() != reflect.String && reflect.ValueOf(item).Field(j).Kind() != reflect.Float64 {
				//每张发票每种类型的字段
				fieldsArr := reflect.ValueOf(item).Field(j).Interface().([][]model.ProjectField)
				for k := 0; k < len(fieldsArr); k++ {
					fields := fieldsArr[k]
					formatObj.Fields = append(formatObj.Fields, fields...)
					yiYuanFields = append(yiYuanFields, fields...)
				}
			}
		}
		fc037 := utils.GetValByFields(item.Invoice[0], "fc037", true)
		var infoMap = FieldsMap{
			Id:    item.Id,
			Code:  item.Code,
			Money: item.Money,
		}
		if item.Id != "other" {
			infoMap.Fields = item.Invoice[0]
			for _, fields := range item.BaoXiaoDan {
				infoMap.Fields = append(infoMap.Fields, fields...)
			}
			for _, fields := range item.ThirdBaoXiaoDan1 {
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
				// qinDanFieldsMap = append(qinDanFieldsMap, FieldsMap{
				// 	Fields: qinDan,
				// })
			}
			// infoMap.QingDan = append(infoMap.QingDan, qinDanFieldsMap...)

			if fc037 != "" {
				fc279 := utils.GetValByFields(item.Invoice[0], "fc279", true)
				var yiYuanMap = YiYuanObj{
					Name:   fc037,
					Fields: append(yiYuanFields),
				}
				for mm, medical := range formatObj.Medical {
					if medical.Name == fc037 {
						yiYuanMap.Name = ""
						if fc279 == "1" {
							formatObj.Medical[mm].Hospital = append(formatObj.Medical[mm].Hospital, infoMap)
						} else if fc279 == "2" {
							formatObj.Medical[mm].Clinc = append(formatObj.Medical[mm].Clinc, infoMap)
						} else if fc279 == "3" {
							formatObj.Medical[mm].Pharmacy = append(formatObj.Medical[mm].Pharmacy, infoMap)
						} else {
							formatObj.Medical[mm].Nonmedical = append(formatObj.Medical[mm].Nonmedical, infoMap)
						}
						formatObj.Medical[mm].Fields = append(formatObj.Medical[mm].Fields, yiYuanFields...)
					}
				}
				if yiYuanMap.Name != "" {
					if fc279 == "1" {
						yiYuanMap.Hospital = append(yiYuanMap.Hospital, infoMap)
					} else if fc279 == "2" {
						yiYuanMap.Clinc = append(yiYuanMap.Clinc, infoMap)
					} else if fc279 == "3" {
						yiYuanMap.Pharmacy = append(yiYuanMap.Pharmacy, infoMap)
					} else {
						yiYuanMap.Nonmedical = append(yiYuanMap.Nonmedical, infoMap)
					}
					formatObj.Medical = append(formatObj.Medical, yiYuanMap)
				}
			}

		} else {
			for _, invoice := range item.Invoice {
				fc478 := utils.GetValByFields(invoice, "fc478", true)
				fc040 := utils.GetValByFields(invoice, "fc040", true)
				fc041 := utils.GetValByFields(invoice, "fc041", true)
				if isNotFaPiao && fc040 != "" {
					if !utils.HasItem(fc040s, fc040) {
						var yiYuanMap = YiYuanObj{
							Name:   fc040,
							Fields: []model.ProjectField{},
						}
						formatObj.Medical = append(formatObj.Medical, yiYuanMap)
					}
					fc040s = append(fc040s, fc040)
				}
				fmt.Println("-------------fc040---fc041------------", fc040, fc041)
				if fc040 == "" && fc041 == "" && fc478 == "" {
					continue
				}
				for mm, medical := range formatObj.Medical {
					fmt.Println("-------------Name-----------", medical.Name)
					infoMap.Fields = append(invoice)
					if medical.Name == fc040 {
						formatObj.Medical[mm].Fields = append(formatObj.Medical[mm].Fields, invoice...)
						formatObj.Medical[mm].Injury = append(formatObj.Medical[mm].Injury, infoMap)
						formatObj.Medical[mm].Injurybj = append(formatObj.Medical[mm].Injurybj, infoMap)
					} else if strings.Index(medical.Name, fc041) != -1 && fc041 != "" {
						if medical.Name == fc041 {
							formatObj.Medical[mm].Fields = append(formatObj.Medical[mm].Fields, invoice...)
							formatObj.Medical[mm].Operation = append(formatObj.Medical[mm].Operation, infoMap)
							formatObj.Medical[mm].Operationbj = append(formatObj.Medical[mm].Operationbj, infoMap)
						} else {
							fc189 := utils.GetValByFields(invoice, "fc189", true)
							re := regexp.MustCompile("_")
							lines := re.Split(medical.Name, -1)
							fmt.Println("------------lines--------------------", medical.Name, lines, fc189)
							a, _ := time.Parse("20060102", fc189)
							b, _ := time.Parse("20060102", lines[1])
							c, _ := time.Parse("20060102", lines[2])
							if (a.After(b) && a.Before(c)) || a.Equal(b) || a.Equal(c) {
								formatObj.Medical[mm].Fields = append(formatObj.Medical[mm].Fields, invoice...)
								formatObj.Medical[mm].Operation = append(formatObj.Medical[mm].Operation, infoMap)
								formatObj.Medical[mm].Operationbj = append(formatObj.Medical[mm].Operationbj, infoMap)

							}
						}
					} else if strings.Index(medical.Name, fc478) != -1 && fc478 != "" {
						if medical.Name == fc478 {
							formatObj.Medical[mm].Fields = append(formatObj.Medical[mm].Fields, invoice...)
						} else {
							fc473 := utils.GetValByFields(invoice, "fc473", true)
							re := regexp.MustCompile("_")
							lines := re.Split(medical.Name, -1)
							a, _ := time.Parse("20060102", fc473)
							b, _ := time.Parse("20060102", lines[1])
							c, _ := time.Parse("20060102", lines[2])
							if (a.After(b) && a.Before(c)) || a.Equal(b) || a.Equal(c) {
								formatObj.Medical[mm].Fields = append(formatObj.Medical[mm].Fields, invoice...)
							}
						}
					}
				}
			}

		}
	}

	fmt.Println("-------------Medical---------------", len(formatObj.Medical))

	// if len(formatObj.ClinicInfo) == 0 {
	// 	formatObj.ClinicInfo = append(formatObj.ClinicInfo, FieldsMap{})
	// }
	// for qq, ClinicInfo := range formatObj.ClinicInfo {
	// 	if len(ClinicInfo.QingDan) == 0 {
	// 		formatObj.ClinicInfo[qq].QingDan = append(formatObj.ClinicInfo[qq].QingDan, FieldsMap{})
	// 	}
	// }
	// if len(formatObj.InpatientInfo) == 0 {
	// 	formatObj.InpatientInfo = append(formatObj.InpatientInfo, FieldsMap{})
	// }
	// for qq, InpatientInfo := range formatObj.InpatientInfo {
	// 	if len(InpatientInfo.QingDan) == 0 {
	// 		formatObj.InpatientInfo[qq].QingDan = append(formatObj.InpatientInfo[qq].QingDan, FieldsMap{})
	// 	}
	// }
	//global.GLog.Info("单号:::" + obj.Bill.BillNum + "	id:::" + obj.Bill.ID)

	return nil, formatObj
}
