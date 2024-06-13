/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/2/14 9:43 上午
 */

package utils

import "server/module/load/model"

func CopyFieldsSlice(src []model.ProjectField) []model.ProjectField {
	if src == nil {
		return nil
	}
	dup := make([]model.ProjectField, len(src))
	copy(dup, src)
	return dup
}

func CopyBlocksSlice(src []model.ProjectBlock) []model.ProjectBlock {
	if src == nil {
		return nil
	}
	dup := make([]model.ProjectBlock, len(src))
	copy(dup, src)
	return dup
}

// IsContain 判断string切片包不包含某个元素
func IsContain(arr []string, item string) bool {
	for _, s := range arr {
		if s == item {
			return true
		}
	}
	return false
}

// 不添加重复数据
func IsRepeateds(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return false
		}
	}
	return true
}

// RemoveElem 将string切片中的item元素去除 返回一个新的切片
func RemoveElem(arr []string, item string) (bool, []string) {
	for i, s := range arr {
		if s == item {
			arr = append(arr[:i], arr[i+1:]...)
			return true, arr
		}
	}
	return false, arr
}
