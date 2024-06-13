/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/31 10:00
 */

package B0113

import (
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"server/global"
	"server/module/download/project"
	utils2 "server/module/monitor/utils"
	"server/module/pro_conf/model"
	"strconv"
	"time"
)

func ScanDownload(conf model.SysFtpMonitor) error {
	billNameArr := make([]string, 0)

	global.GLog.Warn("-scan",
		zap.Any(global.GConfig.System.DownloadBillEnd,
			global.GConfig.System.DownloadPath),
	)
	if global.GConfig.System.DownloadPath == "-1" {
		global.GConfig.System.DownloadPath = ""
	}
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
	//err, stdout, _ := project.ShellOut("curl -u wuserftp:123456 -f -s -k ftp://58.252.228.250:32221/SF/B0113/")
	if err != nil {
		global.GLog.Error("scan run cmd failed", zap.Error(err))
		return err
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
			if global.GConfig.System.DownloadBillEnd != "-1" {
				reg := regexp.MustCompile("[" + global.GConfig.System.DownloadBillEnd + "]$")
				if !reg.MatchString(name) {
					global.GLog.Warn("只下载这些结尾的单", zap.Any("", global.GConfig.System.DownloadBillEnd))
					continue
				}
			}
			parse, err := time.Parse("2006 Jan  2 15:04", strconv.Itoa(time.Now().Year())+" "+arr[2])
			// parse = parse.Add(-8 * time.Hour)
			if err != nil {
				global.GLog.Error("单据解析时间错误:::"+name, zap.Error(err))
				continue
			}

			endTime := time.Now().Add(-time.Duration(conf.Frequency) * time.Minute)
			if parse.Before(endTime) {
				billNameArr = append(billNameArr, name)
			}

		}
	}

	err = utils2.SendMsg(billNameArr, conf)

	return err
}
