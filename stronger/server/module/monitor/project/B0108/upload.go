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
	"server/module/monitor/service"
	"server/utils"
	"strconv"
	"strings"
	"time"
)

func ScanUpload() error {
	//CSB0108RC0284000

	//获取时间范围 2天前 ~ 3分钟前
	startTime := time.Now().Add(3 * -24 * 60 * time.Minute)
	endTime := time.Now().Add(-3 * time.Minute)

	//配置获取
	global.GLog.Warn("B0108-upload-Monitor",
		zap.Any(global.GConfig.System.DownloadBillEnd,
			global.GConfig.System.DownloadPath),
	)
	if global.GConfig.System.DownloadPath == "-1" {
		global.GConfig.System.DownloadPath = ""
	}
	if strings.Index(global.GConfig.System.DownloadPath, "claim_error") == -1 {
		global.GLog.Warn("s", zap.Any("s", "回传监控在通知函下载脚本"))
		return nil
	}
	cache, ok := global.GProConf[global.GConfig.System.ProCode]
	if !ok {
		global.GLog.Error("没有该项目[" + global.GConfig.System.ProCode + "]配置")
		return nil
	}
	if cache.DownloadPaths.Backup == "" {
		global.GLog.Error("该项目[" + global.GConfig.System.ProCode + "]backup配置为空配置")
		return nil
	}

	//获取ftp路径超过三分钟的zip文件
	//cache.DownloadPaths.Backup = "curl -u wuserftp:123456 -f -s -k ftp://192.168.3.45/TPLP/%v"
	scanCmd := fmt.Sprintf(cache.DownloadPaths.Backup, "")
	global.GLog.Info("Backup CMD", zap.Any("", scanCmd))
	err, stdout, _ := project.ShellOut(scanCmd)
	if err != nil {
		global.GLog.Error("Backup run cmd failed", zap.Error(err))
		return err
	}
	global.GLog.Info("B0108-upload notice", zap.Any("scan ", stdout))
	billBatchNums := make([]string, 0)
	re := regexp.MustCompile("[\r\n]+")
	lines := re.Split(stdout, -1)
	global.GLog.Info("B0108-upload notice", zap.Any("lines ", lines))
	lineReg := regexp.MustCompile(`(\S+)\s+(\S+\s+\S+\s+\S+)\s+(\S+)$`)
	for ii, line := range lines {
		global.GLog.Info("B0108-upload notice", zap.Any(line, ii))
		if line == "" {
			continue
		}
		matched, _ := regexp.MatchString(`.zip$`, line)
		arr := lineReg.FindStringSubmatch(line)
		global.GLog.Info("B0108-upload notice", zap.Any("arr", arr))
		if matched && len(arr) >= 4 {
			global.GLog.Info("B0108-upload notice", zap.Any("arr[3]", arr[3]))
			global.GLog.Info("B0108-upload notice", zap.Any("arr[2]", arr[2]))
			//global.GLog.Info("B0108-upload notice", zap.Any("arr[1]", arr[1]))
			arrTime, err := time.ParseInLocation("2006 Jan 2 15:04", strconv.Itoa(time.Now().Year())+" "+arr[2], time.Local)
			if err != nil {
				global.GLog.Error("", zap.Error(err))
				continue
			}
			if arrTime.Before(endTime) {
				billBatchNums = append(billBatchNums, strings.ReplaceAll(arr[3], ".zip", ""))
			}
		}
	}

	//查询该单是否在系统是已回传状态
	err, bills := service.GetIsUploadUnNoticeBills(global.GConfig.System.ProCode, startTime, endTime, billBatchNums)
	if err != nil {
		global.GLog.Error("get bills", zap.Error(err))
		return err
	}

	billBatchNums = []string{}
	for _, bill := range bills {
		billBatchNums = append(billBatchNums, bill.BatchNum)
	}
	global.GLog.Info("B0108-upload notice", zap.Any("billBatchNums", billBatchNums))
	if len(billBatchNums) > 0 {
		msg := "单号:" + strings.Join(billBatchNums, ",") + "，超过3分钟未被客户取走，请核查!"
		obj := map[string]interface{}{
			"title":   "回传通知",                        // 标题
			"msg":     msg,                           // 内容
			"type":    2,                             // 消息类型1:下载2:上传
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
