package utils

import (
	"fmt"
	"strings"
)

func ArrayToString(array []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
}

func HasItem(arr []string, item string) bool {
	for _, i2 := range arr {
		if i2 == item {
			return true
		}
	}
	return false
}

func ItemHasNotInArr(arr []string, item string) bool {
	for _, i2 := range arr {
		if i2 != item {
			return true
		}
	}
	return false
}
