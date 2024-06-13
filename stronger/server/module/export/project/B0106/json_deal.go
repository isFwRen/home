/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/12 14:05
 */

package B0106

import (
	"fmt"
	"server/global"
	"server/utils"
)

// JsonDeal 处理json
func JsonDeal(o interface{}, jsonValue []byte) (err error, newJsonValue []byte) {
	global.GLog.Info("------------------B0106:::JsonDeal-----------------------")
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

	query := `(,|)\s*{\s*"array_list"\s*:\s*""\s*}(,|)`
	data = utils.RegReplace(data, query, "")

	return err, []byte(data)
}
