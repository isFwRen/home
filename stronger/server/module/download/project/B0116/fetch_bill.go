package B0116

import (
	"errors"
	"fmt"
	"os"
	"server/global"
	"server/module/download/project"
	"server/module/download/service"
	"server/module/pro_manager/model"
	"time"

	"go.uber.org/zap"
)

func FetchBill(projectBill model.ProjectBill) (error, model.ProjectBill) {
	err, total := service.CountBill(projectBill.ProCode, projectBill.BillName)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		return err, projectBill
	}
	if total > 0 {
		err, projectBill = Clean(projectBill)
		if err != nil {
			return err, projectBill
		}
		return errors.New("已下载过"), projectBill
	}
	createdAt := time.Unix(projectBill.CreatedAt.Unix(), 0)
	global.GLog.Info("createdAt", zap.Any("", createdAt))
	projectBill.DownloadPath = fmt.Sprintf("%v/download/%v/%v/%v/",
		projectBill.ProCode,
		createdAt.Format("2006"),
		createdAt.Format("01-02"),
		projectBill.BillName)
	global.GLog.Info("DownloadPath", zap.Any("", projectBill.DownloadPath))
	err = os.MkdirAll(global.GConfig.LocalUpload.FilePath+projectBill.DownloadPath, os.ModePerm)
	if err != nil {
		return err, projectBill
	}

	remoteFile0 := projectBill.BillName + ".zip"
	//curl -o '%v' 'ftp://192.168.202.3/HXLP_gexian/BPO_SEND/%v'  -o '%v' 'ftp://192.168.202.3/HXLP_gexian/BPO_SEND/%v'  -u myftp:myftp  -f -s -k
	localFile0 := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + "/" + remoteFile0
	// localFile1 := global.GConfig.LocalUpload.FilePath + projectBill.DownloadPath + "/" + projectBill.Files[1]

	// remoteFile1 := projectBill.BillName + "/" + projectBill.Files[1]
	cmd := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.FetchBill, remoteFile0, localFile0)
	global.GLog.Info("fetch cmd", zap.Any("", cmd))
	err, stdout, stderr := project.ShellOut(cmd)
	global.GLog.Info("stdout", zap.Any("", stdout))
	global.GLog.Error("stderr", zap.Any("", stderr))
	if err != nil {
		global.GLog.Info("", zap.Any("", err))
		return err, projectBill
	}

	return err, projectBill
}
