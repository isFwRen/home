package B0121

import (
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"server/module/download/project"
	"server/module/pro_manager/model"
)

func Clean(projectBill model.ProjectBill) (error, model.ProjectBill) {
	//return nil, projectBill
	rmCMD := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.DownloadClean, projectBill.Files[0], projectBill.Files[0])
	global.GLog.Info("cmd", zap.Any("rmCMD", rmCMD))
	err, stdout, stderr := project.ShellOut(rmCMD)
	if err != nil {
		return err, projectBill
	}
	global.GLog.Info("stderr", zap.Any("", stderr))
	global.GLog.Info("stdout", zap.Any("", stdout))
	return err, projectBill
}