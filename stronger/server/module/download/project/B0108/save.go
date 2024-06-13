package project

import (
	"server/global"
	"server/module/download/service"
	"server/module/pro_manager/model"

	"go.uber.org/zap"
)

func Save(projectBill model.ProjectBill) (error, model.ProjectBill) {
	global.GLog.Info("Save", zap.Any("projectBill", projectBill))
	if projectBill.SaleChannel == "理赔" {
		projectBill.DownloadPath = projectBill.DownloadPath + projectBill.BatchNum + "/" + projectBill.BillNum + "/"
	}
	if projectBill.SaleChannel == "秒赔" {
		projectBill.DownloadPath = projectBill.DownloadPath + projectBill.BatchNum + "/"
	}
	//CSB0108RC0004000
	//CSB0108RC0005000
	//秒赔案件：
	//影像数量≤20页的，案件默认为紧急单；
	//影像数 量＞20页的，案件默认为优先单；
	//2023年10月11日11:56:29 取消需求
	//projectBill.StickLevel = 1
	//global.GLog.Info("projectBill.SaleChannel", zap.Any(projectBill.SaleChannel, len(projectBill.Images)))
	//if projectBill.SaleChannel == "秒赔" {
	//	if len(projectBill.Images) > 20 {
	//		projectBill.StickLevel = 2
	//		//global.GLog.Info("s", zap.Any("s", "i111111111111"))
	//	} else {
	//		projectBill.StickLevel = 1
	//		//global.GLog.Info("s", zap.Any("s", "i222222222222"))
	//	}
	//}
	projectBill.Agency = projectBill.BillName[0:8]
	projectBill.Stage = 1
	err, c := service.InsertBill(projectBill.ProCode, projectBill)
	if err != nil {
		global.GLog.Error("InsertBill-err", zap.Error(err))
		return err, projectBill
	}
	global.GLog.Info("InsertBill", zap.Any("count", c))
	return err, projectBill
}
