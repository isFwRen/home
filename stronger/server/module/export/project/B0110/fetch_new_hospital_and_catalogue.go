/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/15 10:49
 */

package B0110

import (
	"fmt"
	"reflect"
	"server/global"
	"server/module/export/model"
	"server/module/export/service"
	model3 "server/module/load/model"
	"server/utils"
)

func FetchNewHospitalAndCatalogue(obj model.ResultDataBill) error {
	global.GLog.Info("B0110:::FetchNewHospitalAndCatalogue")

	////特殊的常量
	//constSpecialMap := constSpecialDeal(obj.Bill.ProCode)

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

	myCodes := [][]string{
		{"fc012", "fc013", "fc014"},
		{"fc025", "fc032", "fc039"},
		{"fc026", "fc033", "fc040"},
		{"fc027", "fc034", "fc041"},
		{"fc028", "fc035", "fc042"},
		{"fc029", "fc036", "fc043"},
		{"fc030", "fc037", "fc044"},
		{"fc031", "fc038", "fc045"},
	}
	hospitalCatalogues := make([]model.HospitalCatalogue, 0)
	for _, codes := range myCodes {
		for _, oneLoc := range fieldLocationMap[codes[0]] {
			err, one := utils.GetFieldByLoc(obj, oneLoc)
			if err != nil || one.FinalInput != "no_match_const" || one.ResultValue == "" {
				continue
			}
			twoVal := utils.FetchFieldValBySaveBlockIndex(obj, fieldLocationMap[codes[1]], oneLoc[0], oneLoc[1], oneLoc[2], oneLoc[4], false)
			threeVal := utils.FetchFieldValBySaveBlockIndex(obj, fieldLocationMap[codes[2]], oneLoc[0], oneLoc[1], oneLoc[2], oneLoc[4], false)
			catalogue := model.HospitalCatalogue{
				BillId:  obj.Bill.ID,
				BillNum: obj.Bill.BillNum,
				Agency:  obj.Bill.Agency,
				Type:    2,
				Name:    one.ResultValue,
				Num:     twoVal,
				Money:   threeVal,
				Date:    obj.Bill.CreatedAt,
			}
			hospitalCatalogues = append(hospitalCatalogues, catalogue)
		}
	}

	for _, loc := range fieldLocationMap["fc006"] {
		err, one := utils.GetFieldByLoc(obj, loc)
		if err != nil || one.FinalInput != "no_match_const" || one.ResultValue == "" {
			continue
		}
		hospital := model.HospitalCatalogue{
			BillId:  obj.Bill.ID,
			BillNum: obj.Bill.BillNum,
			Agency:  obj.Bill.Agency,
			Type:    1,
			Name:    one.ResultValue,
			Date:    obj.Bill.CreatedAt,
		}
		hospitalCatalogues = append(hospitalCatalogues, hospital)
	}

	return service.UpdateHospitalCatalogue(obj.Bill.ProCode, obj.Bill.ID, hospitalCatalogues)
}
