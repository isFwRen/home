package B0102

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func Scan() []model.ProjectBill {
	global.GLog.Warn("B0102-download-scan",
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
	global.GLog.Info("B0102-download", zap.Any("scan ", stdout))
	re := regexp.MustCompile("[\r\n]+")
	lines := re.Split(stdout, -1)
	global.GLog.Info("B0102-download", zap.Any("lines ", lines))
	lineReg := regexp.MustCompile(`(\S+)\s+(\S+\s+\S+\s+\S+)\s+(\S+)$`)
	for ii, line := range lines {
		global.GLog.Info("B0102-download", zap.Any(line, ii))
		if line == "" {
			continue
		}
		matched, _ := regexp.MatchString(`.zip$`, line)
		arr := lineReg.FindStringSubmatch(line)
		global.GLog.Info("B0102-download", zap.Any("arr", arr))
		if matched && len(arr) >= 4 {
			global.GLog.Info("B0102-download", zap.Any("arr[3]", arr[3]))
			name := strings.Replace(arr[3], ".zip", "", -1)
			var projectBill model.ProjectBill
			projectBill.BillName = name
			projectBill.BatchNum = name
			parse, err := time.Parse("2006 Jan  2 15:04", strconv.Itoa(time.Now().Year())+" "+arr[2])
			if err != nil {
				global.GLog.Error("单据解析时间错误:::"+name, zap.Error(err))
				continue
			}

			projectBill.CreatedAt = parse
			items = append(items, projectBill)

		}
	}
	return items
}
