/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/26 15:08
 */

package project

import (
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/msg_manager/model"
	"strconv"
	"strings"
	"time"
)

func CustomerNoticeScan() []model.CustomerNotice {
	global.GLog.Warn("B0108-download-scan",
		zap.Any(global.GConfig.System.DownloadBillEnd,
			global.GConfig.System.DownloadPath),
	)
	if global.GConfig.System.DownloadPath == "-1" {
		global.GConfig.System.DownloadPath = ""
	}
	var items []model.CustomerNotice
	cache, ok := global.GProConf[global.GConfig.System.ProCode]
	if !ok {
		global.GLog.Error("没有该项目[" + global.GConfig.System.ProCode + "]配置")
		return nil
	}
	if cache.DownloadPaths.Scan == "" {
		global.GLog.Error("该项目[" + global.GConfig.System.ProCode + "]扫描配置为空配置")
		return nil
	}
	scanCmd := fmt.Sprintf(cache.DownloadPaths.Scan, global.GConfig.System.DownloadPath)
	global.GLog.Info("scanCmd", zap.Any("", scanCmd))
	err, stdout, _ := project.ShellOut(scanCmd)
	if err != nil {
		global.GLog.Error("scan run cmd failed", zap.Error(err))
		return items
	}
	global.GLog.Info("B0108-download CustomerNotice", zap.Any("scan ", stdout))
	re := regexp.MustCompile("[\r\n]+")
	lines := re.Split(stdout, -1)
	global.GLog.Info("B0108-download CustomerNotice", zap.Any("lines ", lines))
	lineReg := regexp.MustCompile(`(\S+)\s+(\S+\s+\S+\s+\S+)\s+(\S+)$`)
	for ii, line := range lines {
		global.GLog.Info("B0108-download CustomerNotice", zap.Any(line, ii))
		if line == "" {
			continue
		}
		matched, _ := regexp.MatchString(`.xml$`, line)
		arr := lineReg.FindStringSubmatch(line)
		global.GLog.Info("B0108-download CustomerNotice", zap.Any("arr", arr))
		if matched && len(arr) >= 4 {
			global.GLog.Info("B0108-download CustomerNotice", zap.Any("arr[3]", arr[3]))
			name := strings.Replace(arr[3], ".xml", "", -1)
			var customerNotice model.CustomerNotice
			if strings.Index(global.GConfig.System.DownloadPath, "tpbb") != -1 {
				customerNotice.MsgType = 2
			}
			if strings.Index(global.GConfig.System.DownloadPath, "Claim") != -1 {
				customerNotice.MsgType = 1
			}
			customerNotice.FileName = name
			parse, err := time.Parse("2006 Jan 02 15:04", strconv.Itoa(time.Now().Year())+" "+arr[2])
			if err != nil {
				fmt.Println(strconv.Itoa(time.Now().Year()))
				fmt.Println(arr[2])
				fmt.Println(strconv.Itoa(time.Now().Year()) + " " + arr[2])
				global.GLog.Error("单据解析时间错误:::"+name, zap.Error(err))
				continue
			}
			if parse.Before(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-2, 23, 55, 01, 0, time.Now().Location())) {
				global.GLog.Error("两天前的单不用下载了:::"+name, zap.Any("时间", parse))
				continue
			}
			customerNotice.SendTime = parse
			customerNotice.ProCode = global.GConfig.System.ProCode
			customerNotice.Status = 1
			customerNotice.Content = "当前能否多处理一些业务"
			items = append(items, customerNotice)

		}
	}
	return items
}
