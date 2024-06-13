package B0102

import (
	"server/global"
	"server/module/download/service"
	"server/module/pro_manager/model"

	"go.uber.org/zap"
)

func Save(projectBill model.ProjectBill) (error, model.ProjectBill) {
	global.GLog.Info("", zap.Any("projectBill", projectBill))
	//projectBill.DownloadPath = projectBill.DownloadPath + projectBill.BillName + "/"
	projectBill.Stage = 1
	err, insertBill := service.InsertBill(projectBill.ProCode, projectBill)
	if err != nil {
		return err, insertBill
	}

	return err, insertBill
}
