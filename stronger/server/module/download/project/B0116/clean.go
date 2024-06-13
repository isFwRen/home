package B0116

import (
	"fmt"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
	"time"

	"go.uber.org/zap"
)

func Clean(projectBill model.ProjectBill) (error, model.ProjectBill) {

	dayDir := time.Now().Format("20060102")
	mkdirDay := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.Backup, dayDir)
	global.GLog.Info("cmd", zap.Any("mkdirDay", mkdirDay))
	err, stdout, stderr := project.ShellOut(mkdirDay)
	global.GLog.Info("stdout", zap.Any("mkdirDay", stdout))
	if err != nil || stderr != "" {
		global.GLog.Error("err", zap.Error(err))
		global.GLog.Error("stderr", zap.Any("", stderr))
	}
	reName := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.DownloadClean,
		projectBill.BillName, dayDir, projectBill.BillName)
	global.GLog.Info("cmd", zap.Any("reName", reName))
	err, stdout, _ = project.ShellOut(reName)
	if err != nil {
		return err, projectBill
	}
	global.GLog.Info("stdout", zap.Any("reName", stdout))
	return err, projectBill
}
