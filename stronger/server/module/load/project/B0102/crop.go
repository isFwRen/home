package B0102

import (
	"fmt"
	"server/global"
	load "server/module/load/model"
	"server/module/load/service"
	"server/module/pro_manager/model"
)

func Crop(proCode string, projectBill model.ProjectBill) (error, model.ProjectBill) {
	/* 这是我的第一个简单的程序 */
	proID := global.ProCodeId[proCode]
	fmt.Println("proIDproID:", proID)
	err, temp := service.GetSysProTempByProIdAndName(proID, "未定义")
	fmt.Println("temptemp", err, temp)

	err, blocks := service.GetSysProTempBlockByTempId(temp.ID)

	fmt.Println("err", err, len(blocks))

	var projectBlock load.ProjectBlock

	err = service.DelBlocksByBillID(projectBill.ProCode, projectBill.ID)
	err = service.DelFieldsByBillID(projectBill.ProCode, projectBill.ID)

	// make()
	// blocks = blocks.([]sysConfModel.SysProTempB)
	// return global.ProDbMap[proCode+"_task"].Transaction(func(tx *gorm.DB) error {
	for ii, block := range blocks {
		fmt.Println("block", ii, block.Code)
		// GetSysFields
		var proFields []load.ProjectField
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
		if projectBlock.Level == 0 {
			projectBlock.Level = 99
		}
		projectBlock.Status = 1
		err, fields := service.GetTempBFRelationByBId(block.ID)
		fmt.Println("GetTempBFRelationByBId", err, len(fields), block.ID)
		if len(fields) == 0 {
			continue
		}
		// err, projectBlock = service.InsertBlock(projectBill.ProCode, projectBlock)
		err, projectBlock = service.InsertBlock(projectBill.ProCode+"_task", projectBlock)
		fmt.Println("InsertBlock", err, projectBlock.ID)

		for ff, field := range fields {
			fmt.Println("ffffff", ff, field.FCode)
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
		// service.InsertFields(projectBill.ProCode, proFields)
		service.InsertFields(projectBill.ProCode+"_task", proFields)
		// fmt.Println("InsertBlock", err, projectBlock.ID)
	}

	// err = service.DelBlock("")
	// err = service.DelField("")

	// pictures := projectBill.Picture
	// for ii, block := range blocks {
	// 	fmt.Println("iiii", ii, block)
	// 	// coordinates :=  [] block.Coordinates
	// 	// picPage := block.PicPage
	// 	// if coordinates == nil || picPage == "" {
	// 	// 	continue
	// 	// }
	// 	// src := ""
	// 	// dst := ""
	// 	cmd := "gmic -c - src -crop[-1] x0,y0,x1,y1 -o[-1] dst"
	// 	err, stdout, _ := utils.ShellOut(cmd)
	// 	if err != nil {
	// 		fmt.Println("gmicgmicgmic", err, stdout)
	// 		return err
	// 	}
	// 	// err, _ = service.InsertBlock(block)
	// 	// err, _ = service.InsertField(field)
	// }

	projectBill.Stage = 2
	//1是紧急 2是优先
	if projectBill.Agency == "1" {
		projectBill.StickLevel = 1
	}
	if projectBill.Agency == "2" {
		projectBill.StickLevel = 2
	}
	err = service.UpdateBill(projectBill.ProCode, projectBill, projectBill.ID)
	err, _ = service.InsertBill(projectBill.ProCode+"_task", projectBill)
	fmt.Println("---------------errerrerrerrerr----------------------", err)

	return err, projectBill
}
