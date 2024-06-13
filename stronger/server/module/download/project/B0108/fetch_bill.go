package project

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
	"server/global"
	"server/module/download/project"
	"server/module/download/service"
	"server/module/pro_manager/model"
	"time"
)

func FetchBill(projectBill model.ProjectBill) (error, model.ProjectBill) {
	err, total := service.CountBillByBatchNum(projectBill.ProCode, projectBill.BillName)
	if err != nil {
		global.GLog.Info("FetchBill", zap.Any("CountBill", err))
		return err, projectBill
	}
	if total > 0 {
		global.GLog.Error("FetchBill", zap.Error(errors.New("已下载过")))
		_, projectBill := Clean(projectBill)
		return errors.New("已下载过"), projectBill
	}
	global.GLog.Info("FetchBill", zap.Any("CreatedAt0", projectBill.CreatedAt))
	createdAt := time.Unix(projectBill.CreatedAt.Unix(), 0)
	global.GLog.Info("FetchBill", zap.Any("CreatedAt1", projectBill.CreatedAt))
	projectBill.DownloadPath = projectBill.ProCode + "/download/" + createdAt.Format("2006") + "/" + createdAt.Format("01-02") + "/"
	global.GLog.Info("FetchBill", zap.Any("DownloadPath", projectBill.DownloadPath))
	err = os.MkdirAll(global.GConfig.LocalUpload.FilePath+projectBill.DownloadPath, os.ModePerm)
	if err != nil {
		global.GLog.Error("FetchBill-MkdirAll", zap.Error(err))
		return err, model.ProjectBill{}
	}
	downFile := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + projectBill.BillName + ".zip"
	fileName := projectBill.BillName + ".zip"
	fetchBillCmd := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.FetchBill, global.GConfig.System.DownloadPath, fileName, downFile)
	//fetchBillCmd := fmt.Sprintf("curl 'ftp://192.168.3.45/%v/%v' -f -s -k -o '%v' -u wuserftp:123456", global.GConfig.System.DownloadPath, fileName, downFile)
	global.GLog.Info("FetchBill", zap.Any("fetchBillCmd", fetchBillCmd))
	err, stdout, _ := project.ShellOut(fetchBillCmd)
	global.GLog.Info("ShellOut", zap.Any("stdout", stdout))
	if err != nil {
		global.GLog.Error("ShellOut", zap.Error(err))
		return err, projectBill
	}

	return err, projectBill
}
