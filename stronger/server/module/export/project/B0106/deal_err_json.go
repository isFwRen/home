/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/19 10:53
 */

package B0106

import (
	"bytes"
	"github.com/bitly/go-simplejson"
)

func DealErrJson(obj interface{}, jsonValue []byte) (newJsonValue []byte, err error) {
	data := string(jsonValue)
	// fmt.Println("------------data-------------------", data)
	rcptOtherPayDetailList := `"rcptOtherPayDetailList": ""`
	data = RegReplace(data, rcptOtherPayDetailList, `"rcptOtherPayDetailList": []`)

	costCategoryList := `"costCategoryList": ""`
	data = RegReplace(data, costCategoryList, `"costCategoryList": []`)

	jsonData, err := simplejson.NewFromReader(bytes.NewReader([]byte(data)))
	if err != nil {
		return nil, err
	}
	jsonData.Set("bankAccInfoList", []interface{}{jsonData.Get("bankAccInfoList").GetIndex(0)})
	jsonData.Set("rcptInfoList", []interface{}{jsonData.Get("rcptInfoList").GetIndex(0)})
	return jsonData.MarshalJSON()
}
