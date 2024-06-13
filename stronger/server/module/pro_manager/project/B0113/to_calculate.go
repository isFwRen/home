/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/11/25 14:22
 */

package B0113

import (
	"regexp"
	"server/module/pro_manager/model"
	pro108 "server/module/pro_manager/project/B0108"
	"strings"
	"time"
)

// CalculateBackTimeAndTimeRemaining 计算最晚回传时间和剩余时间
func CalculateBackTimeAndTimeRemaining(bill model.ProjectBill) (backAtTheLatestStr, timeRemaining string, second float64) {
	//CSB0113RC0081000
	//根据下载文件压缩包命名规则，直付标记为1的案件，时效默认更改为20min
	//下载文件命名格式：批次号_该批次下的总案件数_赔案号_该案件在该批次下的顺序号_机构号_直付标记.zip
	deadLine := time.Duration(20)
	backAtTheLatest := bill.ScanAt.Add(deadLine * time.Minute)
	backAtTheLatestStr = backAtTheLatest.Format("2006-01-02 15:04:05")
	timeRemaining = backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
	re, _ := regexp.Compile("[h|m]")
	timeRemaining = strings.Replace(re.ReplaceAllString(timeRemaining, ":"), "s", "", -1)
	second = backAtTheLatest.Sub(time.Now()).Seconds()
	return backAtTheLatestStr, timeRemaining, second
}

// CalculateBackTimeAndRemainder 计算最晚回传时间和剩余时间 (剩余时间使用)
func CalculateBackTimeAndRemainder(bill model.ProjectBill) (backAtTheLatestStr, timeRemaining string, second float64) {
	//CSB0113RC0081000
	//根据下载文件压缩包命名规则，直付标记为1的案件，时效默认更改为20min
	//下载文件命名格式：批次号_该批次下的总案件数_赔案号_该案件在该批次下的顺序号_机构号_直付标记.zip
	deadLine := time.Duration(20)
	backAtTheLatest := bill.ScanAt.Add(deadLine * time.Minute)
	backAtTheLatestStr = backAtTheLatest.Format("2006-01-02 15:04:05")
	timeRemaining = backAtTheLatest.Sub(time.Now()).Round(time.Second).String()
	re, _ := regexp.Compile("[h|m]")
	timeRemaining = strings.Replace(re.ReplaceAllString(timeRemaining, ":"), "s", "", -1)
	timeRemaining = pro108.GetRemainder(timeRemaining)
	second = backAtTheLatest.Sub(time.Now()).Seconds()
	return backAtTheLatestStr, timeRemaining, second
}
