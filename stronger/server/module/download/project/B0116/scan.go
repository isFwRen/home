package B0116

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

func Scan() []model.ProjectBill {
	global.GLog.Warn("-scan",
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
	global.GLog.Info("", zap.Any("scan ", stdout))
	re := regexp.MustCompile("[\r\n]+")
	lines := re.Split(stdout, -1)
	global.GLog.Info("", zap.Any("lines ", lines))
	lineReg := regexp.MustCompile(`(\S+)\s+(\S+\s+\S+\s+\S+)\s+(\S+)$`)
	for ii, line := range lines {
		global.GLog.Info("", zap.Any(line, ii))
		if line == "" {
			continue
		}
		arr := lineReg.FindStringSubmatch(line)
		global.GLog.Info("", zap.Any("arr", arr))
		if len(arr) >= 4 {
			global.GLog.Info("", zap.Any("arr[3]", arr[3]))
			name := arr[3]
			if !utils.RegIsMatch(`\.zip$`, name) {
				continue
			}
			name = strings.Replace(name, ".zip", "", -1)
			var projectBill model.ProjectBill
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

			//扫描文件夹里面是否有xml和zip
			// scanXmlZipCmd := scanCmd + "/" + name + "/"
			// global.GLog.Info("scanXmlZipCmd", zap.Any("", scanXmlZipCmd))
			// err, stdoutXmlZip, _ := project.ShellOut(scanXmlZipCmd)
			// if err != nil {
			// 	global.GLog.Error("scan run cmd failed", zap.Error(err))
			// 	continue
			// }
			// reXmlZip := regexp.MustCompile("[\r\n]+")
			// linesXmlZip := reXmlZip.Split(stdoutXmlZip, -1)
			// global.GLog.Info("", zap.Any("lines ", linesXmlZip))
			// lineRegXmlZip := regexp.MustCompile(`(\S+)\s+(\S+\s+\S+\s+\S+)\s+(\S+)$`)
			// projectBill.Files = make([]string, 2)
			// for _, lineXmlZip := range linesXmlZip {
			// 	global.GLog.Info("", zap.Any(lineXmlZip, ""))
			// 	if lineXmlZip == "" {
			// 		continue
			// 	}
			// 	arrXmlZip := lineRegXmlZip.FindStringSubmatch(lineXmlZip)
			// 	global.GLog.Info("", zap.Any("arrXmlZip", arrXmlZip))
			// 	matched, _ := regexp.MatchString(`(.zip|.xml)$`, lineXmlZip)
			// 	if matched && len(arrXmlZip) >= 4 {
			// 		global.GLog.Info("", zap.Any("arrXmlZip[3]", arrXmlZip[3]))
			// 		if strings.Index(arrXmlZip[3], ".xml") != -1 {
			// 			projectBill.Files[0] = arrXmlZip[3]
			// 		}
			// 		if strings.Index(arrXmlZip[3], ".zip") != -1 {
			// 			projectBill.Files[1] = arrXmlZip[3]
			// 		}
			// 	}
			// }
			// global.GLog.Info("", zap.Any("", projectBill))
			// if projectBill.Files[0] == "" || projectBill.Files[1] == "" {
			// 	global.GLog.Error(name, zap.Error(errors.New("没有压缩包或xml")))
			// 	continue
			// }
			projectBill.BillName = name
			projectBill.CreatedAt = parse
			projectBill.BatchNum = name
			items = append(items, projectBill)

		}
	}
	return items
}
