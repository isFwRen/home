package B0122

import (
	"fmt"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"

	"go.uber.org/zap"
)

func Clean(projectBill model.ProjectBill) (error, model.ProjectBill) {
	//return nil, projectBill
	fileName := projectBill.BillName + ".zip"
	rmCMD := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.DownloadClean, fileName)
	global.GLog.Info("cmd", zap.Any("rmCMD", rmCMD))
	err, stdout, stderr := project.ShellOut(rmCMD)
	if err != nil {
		return err, projectBill
	}
	global.GLog.Info("stderr", zap.Any("", stderr))
	global.GLog.Info("stdout", zap.Any("", stdout))
	return err, projectBill
}
