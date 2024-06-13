package B0110

import (
	"server/global"
	load "server/module/load/model"
	"server/module/load/service"
	"server/module/pro_manager/model"
	"strings"

	uuid "github.com/satori/go.uuid"
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
	var proFields []load.ProjectField

	//未定义的案件默认为优先单
	projectBill.StickLevel = 2

	for ii, block := range blocks {
		// GetSysFields
		global.GLog.Info("block", zap.Any(block.Code, ii))
		projectBlock.ID = strings.Replace(uuid.NewV4().String(), "-", "", -1)
		projectBlock.Code = block.Code
		projectBlock.CreatedAt = projectBill.ScanAt
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
		//if projectBlock.Level == 0 {
		//	projectBlock.Level = 99
		//}
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
		for ff, field := range fields {
			global.GLog.Info("field", zap.Any(field.FCode, ff))
			var proField load.ProjectField
			proField.ID = strings.Replace(uuid.NewV4().String(), "-", "", -1)
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
	}

	projectBill.Stage = 2

	err = service.DelBlocksFieldsAndUpdateBillStageByBillID(projectBill.ProCode, projectBill)
	if err != nil {
		return err, projectBill
	}
	err = service.DelBlocksFieldsAndInsertBillDetail(projectBill.ProCode+"_task", proFields, projectBlock, projectBill)
	return err, projectBill
}
