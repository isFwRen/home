package utils

/**
 * @Author: 6727
 * @Description: 结构体赋值结构体, 如果有其他需要请在相应结构体打算标签,在下面代码继续新增相应的判断和赋值
 */

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"time"
)

// StructAssign
// binding type interface 要修改的结构体
// value type interace 有数据的结构体
func StructAssign(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem() //获取reflect.Type类型
	vVal := reflect.ValueOf(value).Elem()   //获取reflect.Type类型
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		// 在要修改的结构体中查询有数据结构体中相同属性的字段，有则修改其值
		name := vTypeOfT.Field(i).Name
		if ok := bVal.FieldByName(name).IsValid(); ok {
			if vTypeOfT.Field(i).Tag.Get("flag") == "time" {
				//时间类型格式化
				bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface().(time.Time).Format("2006-01-02")))
			} else if vTypeOfT.Field(i).Tag.Get("flag") == "rate" {
				if vVal.Field(i).Interface().(float64) > 0.001 {
					fmt.Println("sa")
					r, _ := decimal.NewFromFloat(vVal.Field(i).Interface().(float64)).Mul(decimal.NewFromFloat(100)).Float64()
					bVal.FieldByName(name).SetString(fmt.Sprintf("%.2f", r) + "%")
				} else {
					bVal.FieldByName(name).SetString(decimal.NewFromFloat(vVal.Field(i).Interface().(float64)).Mul(decimal.NewFromFloat(100)).RoundCeil(2).String() + "%")
				}
			} else {
				bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
			}
		}
	}
}
