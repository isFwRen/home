package B0116

import (
	"fmt"
	"regexp"
	"server/global"
	"server/module/download/project"
	"server/module/download/service"
	"server/module/pro_manager/model"

	"go.uber.org/zap"
)

func Save(projectBill model.ProjectBill) (error, model.ProjectBill) {
	global.GLog.Info("", zap.Any("projectBill", projectBill))
	//projectBill.DownloadPath = projectBill.DownloadPath + projectBill.BillName + "/"
	projectBills := []model.ProjectBill{}
	projectBill.Stage = 1
	for _, bill_num := range projectBill.Files {
		projectBillNew := projectBill
		projectBillNew.BillNum = bill_num
		projectBillNew.DownloadPath = projectBillNew.DownloadPath + projectBill.BatchNum + "/" + bill_num + "/"
		cmd := fmt.Sprintf(`ls %s`, global.GConfig.LocalUpload.FilePath+projectBillNew.DownloadPath)
		global.GLog.Info("", zap.Any("cmd", cmd))
		_, stdout, _ := project.ShellOut(cmd)
		reg := regexp.MustCompile(`[\r\n]+`)
		lines := reg.Split(stdout, -1)
		global.GLog.Info("lines", zap.Any("", lines))
		for _, line := range lines {
			projectBillNew.Images = append(projectBillNew.Images, line)
		}
		projectBillNew.IsAutoUpload = true
		projectBills = append(projectBills, projectBillNew)
	}

	err, _ := service.InsertBills(projectBill.ProCode, projectBills)
	if err != nil {
		return err, projectBill
	}

	return nil, projectBill
}
