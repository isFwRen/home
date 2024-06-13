/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/29 4:01 下午
 */

package utils

import (
	"regexp"
	"server/global"
	"time"

	"github.com/shopspring/decimal"
)

func Validator(validationMap map[int]string, text string, xmlNode string) (note string) {
	reg := ""
	if _, ok := validationMap[1]; ok {
		reg = "^(\\d+)$"
	}
	if _, ok := validationMap[2]; ok {
		reg = "^([\u4e00-\u9fa5]+)$"
	}
	if _, ok := validationMap[3]; ok {
		reg = "^([A-Za-z]+)$"
	}
	//if _, ok := validationMap[4]; ok {
	//	reg = "^-?\\d+\\.?\\d*$"
	//}
	if reg != "" {
		global.GLog.Info("reg::" + reg)
		r := regexp.MustCompile(reg)
		if !r.MatchString(text) {
			note += xmlNode + "数据格式错误;"
		}

	}

	r := regexp.MustCompile("^-?\\d+\\.?\\d*$")
	//r := regexp.MustCompile("^-?\\d+\\.\\d\\d$")
	if _, ok := validationMap[4]; !r.MatchString(text) && ok {
		note += xmlNode + "数据不为金额;"
	}
	if _, ok := validationMap[5]; ok {
		_, err := time.Parse("2006-01-02", text)
		if err != nil {
			// fmt.Println("-----------err---------------", err, text)
			note += xmlNode + "日期格式错误;"
		}
	}
	r = regexp.MustCompile("^1[3|4|5|8][0-9]\\d{4,8}$")
	if _, ok := validationMap[9]; !r.MatchString(text) && ok {
		note += xmlNode + "手机号码错误;"
	}
	r = regexp.MustCompile("\\w{1,}[@][\\w\\-]{1,}([.]([\\w\\-]{1,})){1,3}$")
	if _, ok := validationMap[10]; !r.MatchString(text) && ok {
		note += xmlNode + "邮件格式错误;"
	}

	//r = regexp.MustCompile("^-([1-9]\\d*\\.\\d*|0\\.\\d*[1-9]\\d*)$")
	//if _, ok := validationMap[8]; !r.MatchString(text) && ok {
	//	note += "不为负数;"
	//}

	if _, ok := validationMap[8]; ok {
		fromString, err := decimal.NewFromString(text)
		if err != nil {
			note += err.Error() + ";"
			return
		}
		if fromString.Cmp(decimal.Zero) >= 0 {
			note += xmlNode + "不为负数;"
		}
	}
	return note
}
