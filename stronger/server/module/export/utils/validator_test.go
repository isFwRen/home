/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/29 4:05 下午
 */

package utils

import "testing"

func TestValidator(t *testing.T) {
	var validationMap = map[string]string{
		"2": "2",
	}
	Validator(validationMap, "1111")
}
