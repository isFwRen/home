/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/31 10:00
 */

package B0114

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"server/global"
	"server/module/download/project"
	utils2 "server/module/monitor/utils"
	"server/module/pro_conf/model"
	"strconv"
	"strings"
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
			parse, err := time.Parse("2006 Jan  2 15:04", strconv.Itoa(time.Now().Year())+" "+arr[2])
			// parse = parse.Add(-8 * time.Hour)
			if err != nil {
				global.GLog.Error("单据解析时间错误:::"+name, zap.Error(err))
				continue
			}

			//扫描文件夹里面是否有xml和zip
			scanXmlZipCmd := scanCmd + "/" + name + "/"
			global.GLog.Info("scanXmlZipCmd", zap.Any("", scanXmlZipCmd))
			err, stdoutXmlZip, _ := project.ShellOut(scanXmlZipCmd)
			if err != nil {
				global.GLog.Error("scan run cmd failed", zap.Error(err))
				continue
			}
			reXmlZip := regexp.MustCompile("[\r\n]+")
			linesXmlZip := reXmlZip.Split(stdoutXmlZip, -1)
			global.GLog.Info("", zap.Any("lines ", linesXmlZip))
			lineRegXmlZip := regexp.MustCompile(`(\S+)\s+(\S+\s+\S+\s+\S+)\s+(\S+)$`)
			flag := 0
			for _, lineXmlZip := range linesXmlZip {
				global.GLog.Info("", zap.Any(lineXmlZip, ""))
				if lineXmlZip == "" {
					continue
				}
				arrXmlZip := lineRegXmlZip.FindStringSubmatch(lineXmlZip)
				global.GLog.Info("", zap.Any("arrXmlZip", arrXmlZip))
				matched, _ := regexp.MatchString(`(.zip|.xml)$`, lineXmlZip)
				if matched && len(arrXmlZip) >= 4 {
					global.GLog.Info("", zap.Any("arrXmlZip[3]", arrXmlZip[3]))
					if strings.Index(arrXmlZip[3], ".xml") != -1 {
						flag++
					}
					if strings.Index(arrXmlZip[3], ".zip") != -1 {
						flag++
					}
				}
			}

			if flag != 2 {
				global.GLog.Error(name, zap.Error(errors.New("没有压缩包或xml")))
				continue
			}

			endTime := time.Now().Add(-time.Duration(conf.Frequency) * time.Minute)
			if flag == 2 && parse.Before(endTime) {
				billNameArr = append(billNameArr, name)
			}
		}
	}

	err = utils2.SendMsg(billNameArr, conf)

	return err
}
