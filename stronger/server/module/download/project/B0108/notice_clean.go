/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/31 14:47
 */

package project

import (
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"server/module/download/project"
	model2 "server/module/pro_manager/model"
)

//NoticeClean 通知函删除
func NoticeClean(bill model2.ProjectBill) (err error) {
	fileName := bill.BillNum + ".xml"
	fetchBillCmd := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.DownloadClean, global.GConfig.System.DownloadPath, fileName)
	global.GLog.Info("NoticeClean", zap.Any("NoticeClean cmd", fetchBillCmd))
	//下载文件
	err, stdout, _ := project.ShellOut(fetchBillCmd)
	global.GLog.Info("ShellOut", zap.Any("stdout", stdout))
	if err != nil {
		global.GLog.Error("ShellOut", zap.Error(err))
		return err
	}
	global.GLog.Info("NoticeClean", zap.Any("删除成功", ""))
	return err
}
