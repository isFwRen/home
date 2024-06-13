/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/12 14:05
 */

package B0110

import (
	"fmt"
	"server/global"
	"server/utils"
)

// JsonDeal 处理json
func JsonDeal(o interface{}, jsonValue []byte) (err error, newJsonValue []byte) {
	global.GLog.Info("------------------B0110:::JsonDeal-----------------------")
	// obj := o.(FormatObj)
	// data := make(map[string]interface{})
	// err = json.Unmarshal(jsonValue, &data)
	// if err != nil {
	// 	fmt.Println("------------Unmarshal-------------------", err, obj.Bill.BillNum)
	// }

	data := string(jsonValue)
	fmt.Println("------------data-------------------", data)
	rcptOtherPayDetailList := `"rcptOtherPayDetailList": {\s*"otherSpiltAmnt": "",\s*"uedUnitNo": ""\s*}`
	data = RegReplace(data, rcptOtherPayDetailList, `"rcptOtherPayDetailList": []`)
	rcptOtherPayDetailList = `"rcptOtherPayDetailList": {\s*"uedUnitNo": "",\s*"otherSpiltAmnt": ""\s*}`
	data = RegReplace(data, rcptOtherPayDetailList, `"rcptOtherPayDetailList": []`)

	costCategoryList := `"costCategoryList": {\s*"medicalCategoryName": "",\s*"classifiedAmount": ""\s*}`
	data = RegReplace(data, costCategoryList, `"costCategoryList": []`)
	costCategoryList = `"costCategoryList": {\s*"classifiedAmount": "",\s*"medicalCategoryName": ""\s*}`
	data = RegReplace(data, costCategoryList, `"costCategoryList": []`)

	// fmt.Println("-------------rcptInfoList---------------", jsonValue["rcptInfoList"])

	// fmt.Println("-------------rcptInfoList---------------", data["rcptInfoList"])
	// data.(map[string]interface{})
	// fmt.Println("-------------rcptInfoLists---------------", len(rcptInfoLists))
	// for _, rcptInfoList := range rcptInfoLists {
	// 	costCategoryList := rcptInfoList["costCategoryList"].(map[string]string)
	// 	fmt.Println("-------------classifiedAmount---------------", rcptInfoList)
	// 	// if costCategoryList["classifiedAmount"] == "" && costCategoryList["medicalCategoryName"] == "" {
	// 	// 	rcptInfoLists[ii]["costCategoryList"] = []string{}
	// 	// }
	// 	// rcptOtherPayDetailList := rcptInfoList["rcptOtherPayDetailList"].(map[string]string)
	// 	// if rcptOtherPayDetailList["otherSpiltAmnt"] == "" && rcptOtherPayDetailList["uedUnitNo"] == "" {
	// 	// 	rcptInfoLists[ii]["costCategoryList"] = []string{}
	// 	// }
	// }
	// data["rcptInfoList"] = rcptInfoLists

	query := `(,|)\s*{\s*"array_list"\s*:\s*""\s*}(,|)`
	data = utils.RegReplace(data, query, "")

	// jsonValue, err = json.Marshal(data)
	return err, []byte(data)
}
