/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/12 14:05
 */

package B0102

import (
	"server/global"
)

// JsonDeal 处理json
func JsonDeal(o interface{}, jsonValue []byte) (err error, newJsonValue []byte) {
	global.GLog.Info("------------------B0102:::JsonDeal-----------------------")
	data := string(jsonValue)

	// itemNodes := []string{"beneficiary_records", "medical_records", "bill_records", "fee_records", "injury_records", "operation_records", "questions"}
	// for _, itemNode := range itemNodes {
	// query := `"` + itemNode + `": {\s*""}(,|)`
	// query := `(,|)\s*{\s*"array_list"\s*:\s*""\s*}(,|)`
	// data = utils.RegReplace(data, query, "")
	// }
	// costCategoryList := `"costCategoryList": {\s*"medicalCategoryName": "",\s*"classifiedAmount": ""\s*}`
	// data = utils.RegReplace(data, costCategoryList, `"costCategoryList": []`)

	return err, []byte(data)
}
