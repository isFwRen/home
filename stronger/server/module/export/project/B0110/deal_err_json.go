/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/19 10:53
 */

package B0110

import (
	"bytes"
	"github.com/bitly/go-simplejson"
)

func DealErrJson(obj interface{}, jsonValue []byte) (newJsonValue []byte, err error) {
	//mv, err := mxj.NewMapJson(jsonValue)
	//err = mv.SetValueForPath(nil, "rcptInfoList.id")
	//if err != nil {
	//	return nil, err
	//}
	//err = mv.SetValueForPath(nil, "rcptInfoList.rcptType")
	//if err != nil {
	//	return nil, err
	//}
	//err = mv.SetValueForPath(nil, "rcptInfoList.socialPayAmnt")
	//if err != nil {
	//	return nil, err
	//}
	//return mv.Json()

	jsonData, err := simplejson.NewFromReader(bytes.NewReader(jsonValue))
	if err != nil {
		return nil, err
	}
	jsonData.Set("rcptInfoList", []interface{}{jsonData.Get("rcptInfoList").GetIndex(0)})
	return jsonData.MarshalJSON()
}
