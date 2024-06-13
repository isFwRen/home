/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/8/31 10:00
 */

package B0108

import (
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/utils"
	"strconv"
	"strings"
	"time"
)

func ScanDownload() error {
	endTime := time.Now().Add(-5 * time.Minute)
	address := ""
	global.GLog.Warn("B0108-download-Monitor",
		zap.Any(global.GConfig.System.DownloadBillEnd,
			global.GConfig.System.DownloadPath),
	)
	if global.GConfig.System.DownloadPath == "-1" {
		global.GConfig.System.DownloadPath = ""
	}
	if strings.Index(global.GConfig.System.DownloadPath, "tpbb") != -1 {
		address = "秒赔"
	}
	if strings.Index(global.GConfig.System.DownloadPath, "Claim") != -1 {
		address = "理赔"
	}
	if strings.Index(global.GConfig.System.DownloadPath, "claim_error") != -1 {
		address = "通知函"
	}
	//msg := "下载监控 服务器:" + global.GProConf["B0118"].InnerIp + ",项目:太平理赔/" + address + ";"
	msg := ""
	cache, ok := global.GProConf[global.GConfig.System.ProCode]
	if !ok {
		global.GLog.Error("没有该项目[" + global.GConfig.System.ProCode + "]配置")
		return nil
	}
	if cache.DownloadPaths.Scan == "" {
		global.GLog.Error("该项目[" + global.GConfig.System.ProCode + "]扫描配置为空配置")
		return nil
	}
	//cache.DownloadPaths.Scan = "curl -u wuserftp:123456 -f -s -k ftp://192.168.3.45/TPLP/%v"
	scanCmd := fmt.Sprintf(cache.DownloadPaths.Scan, global.GConfig.System.DownloadPath)
	global.GLog.Info("scanCmd", zap.Any("", scanCmd))
	err, stdout, _ := project.ShellOut(scanCmd)
	if err != nil {
		global.GLog.Error("scan run cmd failed", zap.Error(err))
		return err
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
		matched, _ := regexp.MatchString(`.zip$`, line)
		if address == "通知函" {
			matched, _ = regexp.MatchString(`.xml$`, line)
		}
		arr := lineReg.FindStringSubmatch(line)
		global.GLog.Info("B0108-download", zap.Any("arr", arr))
		if matched && len(arr) >= 4 {
			global.GLog.Info("B0108-download", zap.Any("arr[3]", arr[3]))
			global.GLog.Info("B0108-download", zap.Any("arr[1]", arr[1]))
			f, _ := strconv.ParseFloat(arr[1], 64)
			sizeMB := f / 1048576
			arrTime, err := time.ParseInLocation("2006 Jan 2 15:04", strconv.Itoa(time.Now().Year())+" "+arr[2], time.Local)
			if err != nil {
				global.GLog.Error("", zap.Error(err))
				continue
			}
			if arrTime.Before(endTime) {
				if address == "通知函" {
					msg += "\n单号:" + arr[3] + "，超过" + strconv.FormatFloat(time.Now().Sub(arrTime).Minutes(), 'G', 5, 64) + "分钟未下载!"
				} else {
					msg += "\n批次号:" + arr[3] + "，大小:" + strconv.FormatFloat(sizeMB, 'G', 5, 64) + "MB；"
				}
			}
		}
	}

	global.GLog.Info(address)

	if m, _ := regexp.MatchString(`(批次号|单号)`, msg); m {
		if strings.Index(msg, "批次号") != -1 {
			msg += "\n描述:超过5分钟未下载!"
		}
		//发消息
		//robot := utils.NewRobot("b72bc04ad782bdeea40828d064df3869fd49c4121ed0bbcbb1ec8295508c1b01", "SEC467698ef65ac33aa5e2f04e2c12eb9b5a55325446df8be33594927b91412ab30")
		//err = robot.SendTextMessage(msg, []string{}, true)
		////更新消息
		//if err != nil {
		//	global.GLog.Error("B0118 download Monitor的消息发送失败", zap.Error(err))
		//}
		//global.GLog.Info("B0118 download Monitor的消息发送成功", zap.Any("", msg))

		obj := map[string]interface{}{
			"title":   "下载通知",                        // 标题
			"msg":     msg,                           // 内容
			"type":    1,                             // 消息类型1:下载2:上传
			"stage":   0,                             // '状态（0，正常，1已删除）',
			"proCode": global.GConfig.System.ProCode, // 项目编码
		}
		global.GLog.Info("s", zap.Any("s", obj))
		//广播通知
		err, res := utils.HttpRequest("http://localhost:"+strconv.Itoa(global.GConfig.System.CommonPort)+"/sys-socket-io-notice/business-push", obj)
		if err != nil {
			global.GLog.Error("通知失败::" + ":::" + err.Error())
		}
		global.GLog.Error("通知", zap.Any("res", res))
		global.GLog.Info("通知到/sys-socket-io-notice/business-push成功")
	}

	return err
}
