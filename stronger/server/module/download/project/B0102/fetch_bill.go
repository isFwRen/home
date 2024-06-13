package B0102

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
	projectBill.DownloadPath = fmt.Sprintf("%v%v/download/%v/%v/%v/%v/",
		global.GConfig.LocalUpload.FilePath,
		projectBill.ProCode,
		createdAt.Format("2006"),
		createdAt.Format("01"),
		createdAt.Format("02"),
		projectBill.BillName)
	global.GLog.Info("DownloadPath", zap.Any("", projectBill.DownloadPath))
	err = os.MkdirAll(projectBill.DownloadPath, os.ModePerm)
	if err != nil {
		return err, projectBill
	}

	//curl -o '%v' 'ftp://192.168.202.3/HXLP_gexian/BPO_SEND/%v'  -o '%v' 'ftp://192.168.202.3/HXLP_gexian/BPO_SEND/%v'  -u myftp:myftp  -f -s -k
	fileName := projectBill.BillName + ".zip"
	localFile := projectBill.DownloadPath + fileName
	cmd := fmt.Sprintf(global.GProConf[global.GConfig.System.ProCode].DownloadPaths.FetchBill, fileName, localFile)
	//cmd := fmt.Sprintf("curl 'ftp://58.252.228.250:32221/SF/B0113/%v' -f -s -k -o '%v' -u wuserftp:123456", projectBill.Files[0], localFile)
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
