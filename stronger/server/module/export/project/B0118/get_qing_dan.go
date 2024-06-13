/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/21 9:25 上午
 */

package B0118

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"server/global"
	model2 "server/module/export/model"
	"server/module/load/model"
	model3 "server/module/pro_manager/model"
	"strconv"
	"strings"
	"time"
)

func GetQingDan(form model3.QingDanForm, obj model2.ResultDataBill) (err error, list []model2.QingDan) {
	//bc002
	//医疗项目分类：fc150清单大类属性
	//项目名称：fc084清单项目名称1、fc085清单项目名称2、fc086清单项目名称3、fc087清单项目名称4、fc088清单项目名称5、fc089清单项目名称6、fc090清单项目名称7、fc091清单项目名称8
	//项目类型：fc154清单项目类型1、fc155清单项目类型2、fc156清单项目类型3、fc157清单项目类型4、fc158清单项目类型5、fc159清单项目类型6、fc160清单项目类型7、fc161清单项目类型8
	//项目金额：fc092清单项目名称金额1、fc093清单项目名称金额2、fc094清单项目名称金额3、fc095清单项目名称金额4、fc096清单项目名称金额5、fc097清单项目名称金额6、fc098清单项目名称金额7、fc099清单项目名称金额8
	//数量：默认1
	//项目比例：fc162清单项目自费比例1、fc163清单项目自费比例2、fc164清单项目自费比例3、fc165清单项目自费比例4、fc166清单项目自费比例5fc167清单项目自费比例6、fc168清单项目自费比例7、fc169清单项目自费比例8
	//项目自付：fc172清单项目自付（先付、自费）金额1、fc173清单项目自付（先付、自费）金额2、fc174清单项目自付（先付、自费）金额3、fc175清单项目自付（先付、自费）金额4、fc176清单项目自付（先付、自费）金额5、fc177清单项目自付（先付、自费）金额6、fc178清单项目自付（先付、自费）金额7、fc179清单项目自付（先付、自费）金额8
	timeStart := time.Now()
	arr1 := [][]string{
		{"fc084", "fc154", "fc092", "", "fc162", "fc172"},
		{"fc085", "fc155", "fc093", "", "fc163", "fc173"},
		{"fc086", "fc156", "fc094", "", "fc164", "fc174"},
		{"fc087", "fc157", "fc095", "", "fc165", "fc175"},
		{"fc088", "fc158", "fc096", "", "fc166", "fc176"},
		{"fc089", "fc159", "fc097", "", "fc167", "fc177"},
		{"fc090", "fc160", "fc098", "", "fc168", "fc178"},
		{"fc091", "fc161", "fc099", "", "fc169", "fc179"},
	}
	list = make([]model2.QingDan, 0)
	for _, invoiceMap := range obj.Invoice {
		fieldsArr := invoiceMap.QingDan
		var fields []model.ProjectField
		for _, field := range fieldsArr {
			fields = append(fields, field...)
		}
		code := invoiceMap.Code
		var temp = make([]model2.QingDan, 0)
		temp = getQingDanList(arr1, fields, code, form.ItemName)
		if form.InvoiceNum != "" {
			if form.InvoiceNum == code {
				return nil, temp
			}
		} else {
			list = append(list, temp...)
		}
	}

	global.GLog.Warn("GetQingDan费时:::", zap.Any("毫秒", time.Since(timeStart).Milliseconds()))
	return nil, list
}

//getQingDanList 获取同一发票的清单list
func getQingDanList(arr1 [][]string, fields []model.ProjectField, code string, itemName string) []model2.QingDan {
	var list = make([]model2.QingDan, 0)
	fieldMap := make(map[string]model.ProjectField, len(fields))
	fc152Val := ""
	for _, field := range fields {
		fieldMap[field.BlockID+"_"+field.Code+"_"+strconv.Itoa(field.BlockIndex)] = field
		if field.Code == "fc152" {
			fc152Val = field.ResultValue
		}
	}

	for _, f1 := range fields {
		for _, fieldCode := range arr1 {
			if f1.Code == fieldCode[0] {
				if itemName != "" && strings.Index(f1.FinalValue, itemName) == -1 {
					continue
				}
				f2, _ := decimal.NewFromString(getFinalValue(fieldCode[2], fieldMap, f1))
				f4, _ := decimal.NewFromString(getFinalValue(fieldCode[4], fieldMap, f1))
				qingDan := model2.QingDan{
					InvoiceNum:  code,
					MedicalType: fc152Val,
					Name:        f1.FinalValue,
					Type:        getFinalValue(fieldCode[1], fieldMap, f1),
					Price:       getFinalValue(fieldCode[2], fieldMap, f1),
					Count:       "1",
					Percent:     getFinalValue(fieldCode[4], fieldMap, f1),
					Pay:         f2.Mul(f4).Round(2).String(),
				}
				if f1.FinalValue == "" {
					list = append(list, qingDan)
				}
			}
		}
	}
	return list
}

func getFinalValue(fieldCode string, fieldMap map[string]model.ProjectField, f model.ProjectField) string {
	f2, ok := fieldMap[f.BlockID+"_"+fieldCode+"_"+strconv.Itoa(f.BlockIndex)]
	if !ok {
		return ""
	}
	return f2.FinalValue
}

func getResultValue(fieldCode string, fieldMap map[string]model.ProjectField, f model.ProjectField) string {
	f2, ok := fieldMap[f.BlockID+"_"+fieldCode+"_"+strconv.Itoa(f.BlockIndex)]
	if !ok {
		return ""
	}
	return f2.ResultValue
}
