/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/11/25 14:22
 */

package B0108

import (
	"regexp"
	"server/module/pro_manager/model"
	"strings"
	"time"
)

// CalculateBackTimeAndTimeRemaining 计算最晚回传时间和剩余时间
func CalculateBackTimeAndTimeRemaining(bill model.ProjectBill) (backAtTheLatestStr, timeRemaining string, second float64) {
	//CSB0108RC0003000
	//回传时效：调整时效表的最晚回传时间和剩余时间
	//下载路径为\Claim的案件（理赔），时效为1小时；
	//下载路径为\tpbb时的案件（秒赔），时效为30分钟。
	deadLine := time.Duration(60)
	if bill.SaleChannel == "秒赔" {
		deadLine = time.Duration(30)
	}
	backAtTheLatest := bill.ScanAt.Add(deadLine * time.Minute)

	backAtTheLatestStr = backAtTheLatest.Format("2006-01-02 15:04:05")
	timeRemaining = backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
	re, _ := regexp.Compile("[h|m]")
	timeRemaining = strings.Replace(re.ReplaceAllString(timeRemaining, ":"), "s", "", -1)
	second = backAtTheLatest.Sub(time.Now()).Seconds()
	return backAtTheLatestStr, timeRemaining, second
}

// CalculateBackTimeAndTimeRemaining 计算最晚回传时间和剩余时间(剩余时间使用)
func CalculateBackTimeAndRemainder(bill model.ProjectBill) (backAtTheLatestStr, timeRemaining string, second float64) {
	//CSB0108RC0003000
	//回传时效：调整时效表的最晚回传时间和剩余时间
	//下载路径为\Claim的案件（理赔），时效为1小时；
	//下载路径为\tpbb时的案件（秒赔），时效为30分钟。
	deadLine := time.Duration(60)
	if bill.SaleChannel == "秒赔" {
		deadLine = time.Duration(30)
	}
	backAtTheLatest := bill.ScanAt.Add(deadLine * time.Minute)

	backAtTheLatestStr = backAtTheLatest.Format("2006-01-02 15:04:05")
	timeRemaining = backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
	// 时效格式 00:00:00
	re, _ := regexp.Compile("[h|m]")
	timeRemaining = strings.Replace(re.ReplaceAllString(timeRemaining, ":"), "s", "", -1)
	second = backAtTheLatest.Sub(time.Now()).Seconds()
	timeRemaining = GetRemainder(timeRemaining)
	return backAtTheLatestStr, timeRemaining, second
}

// CalculateRequirementOfAging 获取回传时效要求
func CalculateRequirementOfAging(bill model.ProjectBill) (requirementOfAging string) {
	requirementOfAging = "60"
	if bill.SaleChannel == "秒赔" {
		requirementOfAging = "30"
	}
	return requirementOfAging
}

// 获取剩余时间
func GetRemainder(sf string) string {
	split := strings.Split(sf, ":")
	var reStr string
	switch len(split) {
	case 3:
		reStr = ProcessingFormat(split, len(split))
		return reStr
	case 2:
		reStr = ProcessingFormat(split, len(split))
		return reStr
	case 1:
		reStr = ProcessingFormat(split, len(split))
		return reStr
	}
	return ""
}
func ProcessingFormat(arr []string, length int) string {
	var reStr string
	if length == 3 {
		if strings.HasPrefix(arr[0], "-") && len(arr) != 0 && len(arr[0]) == 2 {
			arr[0] = arr[0][1:]
			arr[0] = "-0" + arr[0]
		}
		for _, item := range arr {
			if len(item) < 2 {
				item = "0" + item
			}
			reStr = reStr + item + ":"
		}
		if strings.HasSuffix(reStr, ":") {
			reStr = reStr[:len(reStr)-1]
		}
		return reStr
	} else if length == 2 {
		if strings.HasPrefix(arr[0], "-") && len(arr) != 0 && len(arr[0]) > 1 {
			reStr = "-00:"
			arr[0] = arr[0][1:]
		} else {
			reStr = "00:"
		}
		for _, item := range arr {
			if len(item) < 2 {
				item = "0" + item
			}
			reStr = reStr + item + ":"
		}
		if strings.HasSuffix(reStr, ":") {
			reStr = reStr[:len(reStr)-1]
		}
		return reStr
	} else if length == 1 {
		if strings.HasPrefix(arr[0], "-") && len(arr) != 0 && len(arr[0]) > 1 {
			reStr = "-00:00:"
			arr[0] = arr[0][1:]
		} else {
			reStr = "00:00:"
		}
		for _, item := range arr {
			if len(item) < 2 {
				item = "0" + item
			}
			reStr = reStr + item + ":"
		}
		if strings.HasSuffix(reStr, ":") {
			reStr = reStr[:len(reStr)-1]
		}
		return reStr
	}
	return ""
}
