/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/3/7 10:45 上午
 */

package utils

import (
	"errors"
	"server/module/load/model"
)

func GetFieldValue(items []model.ProjectField, code string) (string, error) {
	for _, item := range items {
		if item.Code == code {
			return item.FinalValue, nil
		}
	}
	return "", errors.New("not exit " + code)
}
