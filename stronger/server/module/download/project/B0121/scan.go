/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023年03月10日17:07:15
 */

package B0121

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
	//err, stdout, _ := project.ShellOut("curl -u wuserftp:123456 -f -s -k ftp://58.252.228.250:32221/SF/B0113/")
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
			var projectBill model.ProjectBill
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

			//202303081051508306_5_8060000004641368_1_8641_0.zip
			//批次号_该批次下的总案件数_赔案号_该案件在该批次下的顺序号_机构号_直付标记.zip
			projectBill.BillName = strings.Replace(name, ".zip", "", -1)
			zipNameArr := strings.Split(projectBill.BillName, "_")
			if len(zipNameArr) != 6 {
				global.GLog.Error("案件号压缩包名字有误", zap.Any("", name))
			}
			projectBill.BatchNum = zipNameArr[0]
			projectBill.BillNum = zipNameArr[2]

			re = regexp.MustCompile("^8065")
			if !re.MatchString(projectBill.BillNum) {
				global.GLog.Error("该单据不是团险：" + projectBill.BillName)
				continue
			}
			projectBill.Agency = zipNameArr[4]
			projectBill.Files = []string{name}
			projectBill.CreatedAt = parse
			items = append(items, projectBill)

		}
	}
	return items
}
