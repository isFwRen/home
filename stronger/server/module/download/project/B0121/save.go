package B0121

import (
	"go.uber.org/zap"
	"server/global"
	"server/module/download/service"
	"server/module/pro_manager/model"
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
