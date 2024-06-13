package project

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
	/* 这是我的第一个简单的程序 */
	fmt.Println("B0118 scan!")
	items := []model.ProjectBill{}
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
	//err, stdout, _ := project.ShellOut("curl -u myftp:myftp -f -s -k ftp://192.168.202.3/ZYLP/")
	//err, stdout, _ := project.ShellOut("curl -u ICSBPOzhuhaihuiliu:ICSBPOzhuhaihuiliu#123 -f -s -k ftps://202.108.196.190:991/input/")
	//err, stdout, _ := project.ShellOut("curl -u wuserftp:123456 -f -s -k ftp://192.168.3.45/SF/")
	//cmd := "docker exec -it download_system /usr/bin/curl -u ICSBPOzhuhaihuiliu:test -f -s -k ftps://202.108.196.163:990/input/"
	//cmd := "docker exec -it download_system /usr/bin/curl -u ICSBPOzhuhaihuiliu:ICSBPOzhuhaihuiliu#123 -f -s -k ftps://202.108.196.190:991/input/"
	err, stdout, stderr := project.SpecialShell(cmd)
	//err, stdout, _ := project.ShellOut(cmd)
	fmt.Println("stderr:", stderr)
	fmt.Println("errerrerr:", err)
	fmt.Println("stdout:", stdout)
	if err != nil {
		fmt.Printf("cmd.Run failed with %s\n", err)
		return items
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
			var projectBill model.ProjectBill
			projectBill.BillName = name
			// ProjectBill.BillNum = name
			parse, err := time.Parse("2006 Jan 02 15:04", strconv.Itoa(time.Now().Year())+" "+arrs[2])
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
			// ProjectBill.DownloadAt = time.Now()
			items = append(items, projectBill)

		}
	}
	// fmt.Println("items:", items)
	return items
	// fmt.Println(lines[0])
	// cmd := ""
}
