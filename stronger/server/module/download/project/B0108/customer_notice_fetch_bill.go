/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/26 15:23
 */

package project

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"server/global"
	"server/module/download/project"
	"server/module/download/service"
	"server/module/msg_manager/model"
	"server/utils"
	"strconv"
	"time"
)

func CustomerNoticeFetchBill(customerNotice model.CustomerNotice) error {
	err, total := service.CountBillByName(customerNotice.ProCode, customerNotice.FileName)
	if err != nil {
		global.GLog.Info("FetchBill", zap.Any("CountBill", err))
		return err
	}
	if total > 0 {
		global.GLog.Error("FetchBill", zap.Error(errors.New("已下载过")))
		err = CustomerNoticeClean(customerNotice)
		return errors.New("已下载过")
	}
	global.GLog.Info("FetchBill", zap.Any("SendTime", customerNotice.SendTime))
	sendTime := time.Unix(customerNotice.SendTime.Unix(), 0)
	global.GLog.Info("FetchBill", zap.Any("sendTime", customerNotice.SendTime))
	customerNotice.DownloadPath = customerNotice.ProCode + "/download/" + sendTime.Format("2006") + "/" + sendTime.Format("01-02") + "/" + strconv.Itoa(customerNotice.MsgType) + "/"
	global.GLog.Info("FetchBill", zap.Any("DownloadPath", customerNotice.DownloadPath))
	err = os.MkdirAll(global.GConfig.LocalUpload.FilePath+customerNotice.DownloadPath, os.ModePerm)
	if err != nil {
		global.GLog.Error("FetchBill-MkdirAll", zap.Error(err))
		return err
	}
	downFile := global.GConfig.LocalUpload.FilePath + customerNotice.DownloadPath + customerNotice.FileName + ".xml"
	fileName := customerNotice.FileName + ".xml"
	fetchBillCmd := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.FetchBill, global.GConfig.System.DownloadPath, fileName, downFile)
	//fetchBillCmd := fmt.Sprintf("curl 'ftp://192.168.3.45/%v/%v' -f -s -k -o '%v' -u wuserftp:123456", global.GConfig.System.DownloadPath, fileName, downFile)
	global.GLog.Info("FetchBill", zap.Any("fetchBillCmd", fetchBillCmd))
	//下载文件
	err, stdout, _ := project.ShellOut(fetchBillCmd)
	global.GLog.Info("ShellOut", zap.Any("stdout", stdout))
	if err != nil {
		global.GLog.Error("ShellOut", zap.Error(err))
		return err
	}

	//保存到数据库
	err = service.InsertCustomerNotice(global.GConfig.System.ProCode, customerNotice)
	if err != nil {
		global.GLog.Error("InsertCustomerNotice", zap.Error(err))
		return err
	}
	global.GLog.Info("下载成功")
	fmt.Println(global.GConfig.System.CommonPort)

	err, count := service.GetCustomerNoticeByStatus(global.GConfig.System.ProCode)
	if err != nil {
		return err
	}
	//广播通知
	err, res := utils.HttpRequest("http://localhost:"+strconv.Itoa(global.GConfig.System.CommonPort)+"/sys-socket-io-notice/customer-notice", map[string]interface{}{
		"fileName": customerNotice.FileName,
		"sendTime": customerNotice.SendTime,
		"proCode":  customerNotice.ProCode,
		"content":  customerNotice.Content,
		"count":    count,
	})
	if err != nil {
		global.GLog.Error("通知失败::" + ":::" + err.Error())
	}
	global.GLog.Error("通知", zap.Any("res", res))
	global.GLog.Info("通知到/sys-socket-io-notice/customer-notice成功")

	//删除客户路径文件
	err = CustomerNoticeClean(customerNotice)
	return err
}
