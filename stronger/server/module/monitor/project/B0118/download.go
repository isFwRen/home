/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/31 10:00
 */

package B0118

import (
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

	cache, ok := global.GProConf[global.GConfig.System.ProCode]
	if !ok {
		global.GLog.Error("没有该项目[" + global.GConfig.System.ProCode + "]配置")
		return nil
	}
	if cache.DownloadPaths.Scan == "" {
		global.GLog.Error("该项目[" + global.GConfig.System.ProCode + "]扫描配置为空配置")
		return nil
	}
	cmd := cache.DownloadPaths.Scan
	global.GLog.Info("cmd", zap.Any("", cmd))
	err, stdout, stderr := project.SpecialShell(cmd)
	//err, stdout, _ := project.ShellOut(cmd)
	fmt.Println("stderr:", stderr)
	fmt.Println("errerrerr:", err)
	fmt.Println("stdout:", stdout)
	if err != nil {
		fmt.Printf("cmd.Run failed with %s\n", err)
		return err
	}
	re := regexp.MustCompile("[\r\n]+")
	lines := re.Split(stdout, -1)
	fmt.Println("lines:", lines)
	lineReg := regexp.MustCompile(`(\S+)\s+(\S+\s+\S+\s+\S+)\s+(\S+)$`)
	for ii, line := range lines {
		// line := lines[i]
		fmt.Println("iiiiii:", ii, line)
		if line == "" {
			continue
		}
		// fmt.Println("line:", line)
		matched, _ := regexp.MatchString(`.zip$`, line)
		arrs := lineReg.FindStringSubmatch(line)
		fmt.Println("arrs:", arrs, len(arrs))
		if matched && len(arrs) >= 4 {
			fmt.Println("arrs:", arrs[3])
			// items[0] = arrs[3]
			name := strings.Replace(arrs[3], ".zip", "", -1)
			// ProjectBill.BillNum = name
			parse, err := time.Parse("2006 Jan 02 15:04", strconv.Itoa(time.Now().Year())+" "+arrs[2])
			parse = parse.Add(-8 * time.Hour)
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
