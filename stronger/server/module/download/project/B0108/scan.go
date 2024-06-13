package project

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func Scan() []model.ProjectBill {
	global.GLog.Warn("B0108-download-scan",
		zap.Any(global.GConfig.System.DownloadBillEnd,
			global.GConfig.System.DownloadPath),
	)
	if global.GConfig.System.DownloadPath == "-1" {
		global.GConfig.System.DownloadPath = ""
	}
	var items []model.ProjectBill
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
	global.GLog.Info("B0108-download", zap.Any("scan ", stdout))
	re := regexp.MustCompile("[\r\n]+")
	lines := re.Split(stdout, -1)
	global.GLog.Info("B0108-download", zap.Any("lines ", lines))
	lineReg := regexp.MustCompile(`(\S+)\s+(\S+\s+\S+\s+\S+)\s+(\S+)$`)
	for ii, line := range lines {
		global.GLog.Info("B0108-download", zap.Any(line, ii))
		if line == "" {
			continue
		}
		matched, _ := regexp.MatchString(`.log$`, line)
		arr := lineReg.FindStringSubmatch(line)
		global.GLog.Info("B0108-download", zap.Any("arr", arr))
		if matched && len(arr) >= 4 {
			global.GLog.Info("B0108-download", zap.Any("arr[3]", arr[3]))
			name := strings.Replace(arr[3], ".log", "", -1)
			var projectBill model.ProjectBill
			if strings.Index(global.GConfig.System.DownloadPath, "tpbb") != -1 {
				projectBill.SaleChannel = "秒赔"
			}
			if strings.Index(global.GConfig.System.DownloadPath, "Claim") != -1 {
				projectBill.SaleChannel = "理赔"
			}
			projectBill.BillName = name
			projectBill.BatchNum = name
			if global.GConfig.System.DownloadBillEnd != "-1" {
				reg := regexp.MustCompile("[" + global.GConfig.System.DownloadBillEnd + "]$")
				if !reg.MatchString(name) {
					global.GLog.Warn("只下载这些结尾的单", zap.Any("", global.GConfig.System.DownloadBillEnd))
					continue
				}
			}
			parse, err := time.Parse("2006 Jan  2 15:04", strconv.Itoa(time.Now().Year())+" "+arr[2])
			parse = parse.Add(-8 * time.Hour)
			if err != nil {
				global.GLog.Error("单据解析时间错误:::"+name, zap.Error(err))
				continue
			}
			if parse.Before(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day()-2, 23, 55, 01, 0, time.Now().Location())) {
				global.GLog.Error("两天前的单不用下载了:::"+name, zap.Any("时间", parse))
				continue
			}
			projectBill.CreatedAt = parse
			items = append(items, projectBill)

		}
	}
	sort.SliceStable(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	return items
}
