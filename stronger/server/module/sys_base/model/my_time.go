package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

//MyTime 自定义时间 yyyy-MM-dd HH:mm:ss
type MyTime time.Time

//MyTimeHHMMSS 自定义时间 HH:mm:ss
type MyTimeHHMMSS time.Time

var yyyyMMddHHmmss = "2006-01-02 15:04:05"
var HHmmss = "15:04:05"

func (t *MyTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.Parse(yyyyMMddHHmmss, timeStr)
	*t = MyTime(t1)
	return err
}

func (t MyTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format(yyyyMMddHHmmss))
	return []byte(formatted), nil
}

func (t MyTime) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(yyyyMMddHHmmss), nil
}

func (t *MyTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = MyTime(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *MyTime) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}

func (t *MyTimeHHMMSS) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	reflect.TypeOf(t)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	timeStr = strings.TrimPrefix(timeStr, "\"")
	t1, err := time.Parse(yyyyMMddHHmmss, "2021-01-01 "+timeStr)
	*t = MyTimeHHMMSS(t1)
	return err
}

func (t MyTimeHHMMSS) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%v\"", time.Time(t).Format(HHmmss))
	return []byte(formatted), nil
}

func (t MyTimeHHMMSS) Value() (driver.Value, error) {
	// MyTime 转换成 time.Time 类型
	tTime := time.Time(t)
	return tTime.Format(yyyyMMddHHmmss), nil
}

func (t *MyTimeHHMMSS) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = MyTimeHHMMSS(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

func (t *MyTimeHHMMSS) String() string {
	return fmt.Sprintf("hhh:%s", time.Time(*t).String())
}
