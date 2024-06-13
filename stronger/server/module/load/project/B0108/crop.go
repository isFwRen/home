package B0108

import (
	"server/global"
	load "server/module/load/model"
	"server/module/load/service"
	"server/module/pro_manager/model"
	"time"

	"go.uber.org/zap"
)

func Crop(proCode string, projectBill model.ProjectBill) (error, model.ProjectBill) {
	proID := global.ProCodeId[proCode]
	global.GLog.Info("crop", zap.Any("proId", proID))
	err, temp := service.GetSysProTempByProIdAndName(proID, "未定义")
	global.GLog.Info("crop", zap.Any("temp", temp))
	if err != nil {
		global.GLog.Error("crop-GetSysProTempByProIdAndName", zap.Error(err))
	}

	err, blocks := service.GetSysProTempBlockByTempId(temp.ID)
	global.GLog.Info("crop", zap.Any("len(blocks)", len(blocks)))
	if err != nil {
		global.GLog.Error("crop-GetSysProTempBlockByTempId", zap.Error(err))
	}

	var projectBlock load.ProjectBlock
	err = service.DelBlocksByBillID(projectBill.ProCode, projectBill.ID)
	err = service.DelFieldsByBillID(projectBill.ProCode, projectBill.ID)
	err = service.DelBlocksByBillID(projectBill.ProCode+"_task", projectBill.ID)
	err = service.DelFieldsByBillID(projectBill.ProCode+"_task", projectBill.ID)

	for ii, block := range blocks {
		// GetSysFields
		global.GLog.Info("block", zap.Any(block.Code, ii))
		var proFields []load.ProjectField
		projectBlock.Code = block.Code
		projectBlock.CreatedAt = projectBill.ScanAt
		if projectBill.SaleChannel == "理赔" {
			projectBlock.CreatedAt = projectBill.ScanAt.Add(20 * time.Minute)
		}
		projectBlock.Temp = temp.Name
		projectBlock.Name = block.Name
		projectBlock.BillID = projectBill.ID
		projectBlock.FEight = block.FEight
		projectBlock.Ocr = block.Ocr
		projectBlock.FreeTime = block.FreeTime
		projectBlock.IsLoop = block.IsLoop
		projectBlock.PicPage = block.PicPage
		projectBlock.WCoordinate = block.WCoordinate
		projectBlock.MCoordinate = block.MCoordinate
		projectBlock.Op0Stage = "op0"
		projectBlock.Level = projectBill.StickLevel
		if projectBlock.Level == 0 {
			projectBlock.Level = 99
		}
		projectBlock.Status = 1
		err, fields := service.GetTempBFRelationByBId(block.ID)
		global.GLog.Info("crop-GetTempBFRelationByBId", zap.Any("len(fields)", len(fields)))
		global.GLog.Info("crop-GetTempBFRelationByBId", zap.Any("block.ID", block.ID))
		if err != nil {
			global.GLog.Error("crop-GetTempBFRelationByBId", zap.Error(err))
		}
		if len(fields) == 0 {
			continue
		}
		err, projectBlock = service.InsertBlock(projectBill.ProCode+"_task", projectBlock)
		global.GLog.Info("crop", zap.Any("InsertBlock", projectBlock.ID))
		if err != nil {
			global.GLog.Error("crop-InsertBlock", zap.Error(err))
		}
		for ff, field := range fields {
			global.GLog.Info("field", zap.Any(field.FCode, ff))
			var proField load.ProjectField
			proField.BillID = projectBill.ID
			proField.BlockID = projectBlock.ID
			proField.Code = field.FCode
			proField.Name = field.FName
			proField.BlockIndex = 0
			proField.FieldIndex = ff
			proField.Op0Input = "no"
			if field.FName == "模板类型字段" {
				proField.Op0Input = "yes"
			}
			proFields = append(proFields, proField)
		}
		err = service.InsertFields(projectBill.ProCode+"_task", proFields)
		if err != nil {
			global.GLog.Error("crop-InsertFields", zap.Error(err))
		}
	}

	projectBill.Stage = 2
	err = service.UpdateBill(projectBill.ProCode, projectBill, projectBill.ID)
	err, _ = service.InsertBill(projectBill.ProCode+"_task", projectBill)

	return err, projectBill
}
