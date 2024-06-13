/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/31 14:16
 */

package project

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"server/global"
	"server/module/download/project"
	"server/module/download/service"
	model2 "server/module/pro_manager/model"
	"server/utils"

	"go.uber.org/zap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//NoticeFetchBill 通知函下载
func NoticeFetchBill(proCode string, item model2.ProjectBill) (err error) {
	//获取这张单是否已接收
	err, bill := service.FindBillByStage(proCode, item.BillNum)
	if err != nil {
		global.GLog.Info("FetchBill", zap.Any("NoticeFetchBill", err))
		// err = NoticeClean(item)
		// if err != nil {
		// 	return errors.New("删除失败")
		// }
		return nil
	}
	if bill.ID == "" {
		global.GLog.Error("FetchBill", zap.Error(errors.New("不存在这张单")))
		//err = NoticeClean(item)
		return errors.New("没有找到该单据")
	}

	//下载xml
	downFile := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + item.BillNum + ".xml"
	fileName := bill.BillNum + ".xml"
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
	global.GLog.Info("FetchBill", zap.Any("下载成功", ""))

	//解析xml
	data, err := ioutil.ReadFile(downFile)
	if err != nil {
		global.GLog.Error(downFile+" File reading error", zap.Error(err))
		return err
	}
	readByte, err := ioutil.ReadAll(transform.NewReader(bytes.NewBuffer(data), simplifiedchinese.GBK.NewDecoder()))
	if err != nil {
		global.GLog.Error(downFile+" File read Byte  error", zap.Error(err))
		return err
	}
	packCode := utils.GetNodeValue(string(readByte), "pack_code")
	errorValue := utils.GetNodeValue(string(readByte), "error_value")
	global.GLog.Info("FetchBill", zap.Any("packCode", packCode))
	global.GLog.Info("FetchBill", zap.Any("errorValue", errorValue))
	if packCode == "" {
		packCode = "true"
	}
	bill.PackCode = packCode
	bill.WrongNote = errorValue
	global.GLog.Info("FetchBill", zap.Any("解析成功", ""))

	//更新数据库
	err = service.UpdateBillPackCode(bill)
	if err != nil {
		global.GLog.Error("更新数据库", zap.Error(err))
		return err
	}
	global.GLog.Info("FetchBill", zap.Any("更新成功", ""))

	//删除ftp文件
	err = NoticeClean(item)

	return err
}
