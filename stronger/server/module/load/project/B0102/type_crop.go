package B0102

import (
	"errors"
	"fmt"
	"server/global"
	_func "server/global/func"
	"server/module/load/model"
	"server/module/load/service"
	pModel "server/module/pro_conf/model"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func TypeCrop(proCode string, block model.ProjectBlock, tempBlocks []pModel.SysProTempB, CacheFieldConf map[string]pModel.SysProField) error {
	/* 这是我的第一个简单的程序 */
	proID := global.ProCodeId[proCode]
	fmt.Println("proIDproID", proID)
	// err, tFields := service.GetTempBFRelationByBId(block.ID)
	// fmt.Println("GetTempBFRelationByBId", err, len(tFields), block.ID)
	// err, fields := service.SelectCropFields(proCode, block.ID, "模板类型字段")
	err, bill := service.SelectBillByID(proCode+"_task", block.BillID)
	fmt.Println("111111111", err)
	err, fields := service.SelectFieldsByBlockID(proCode+"_task", block.ID)
	fmt.Println("2222222222", err)
	if err != nil {
		return err
	}
	//项目自定义校验需求
	err = proCustomizeValidation(fields, proCode)
	if err != nil {
		return err
	}
	fc040s := []string{}
	opFields := _func.FieldsFormat(fields)
	return global.ProDbMap[proCode+"_task"].Transaction(func(tx *gorm.DB) error {
		for ii, fFields := range opFields {
			value := getValue(fFields, "模板类型字段", CacheFieldConf)
			page := getValue(fFields, "图片页码", CacheFieldConf)
			w := getValue(fFields, "显示范围", CacheFieldConf)
			fmt.Println("------valuevalue---------:", value, page, w)
			blocks := getBlocks(tempBlocks, value)
			fc040 := getCodeValue(fFields, "fc040")
			for gg, gBlock := range blocks {
				fmt.Println("gBlockgBlock:", gBlock)
				err, bFields := service.GetTempBFRelationByBId(gBlock.ID)
				fmt.Println("------bFields---------:", len(bFields))
				if err != nil {
					return err
				}
				iblock := initBlock(gBlock, block.BillID)
				iblock.Zero = ii
				if w != "" {
					iblock.Picture = w
				} else if page != "" {
					idx, _ := strconv.Atoi(page)
					pic := bill.Pictures[idx]
					// fmt.Println("-----------iblock.WCoordinate----------- ", iblock.WCoordinate, iblock.WCoordinate[0] == "")
					if iblock.WCoordinate == nil || len(iblock.WCoordinate) == 0 {
						iblock.Picture = pic
					} else if len(iblock.WCoordinate) == 4 && iblock.WCoordinate[0] == "" && iblock.WCoordinate[1] == "" && iblock.WCoordinate[2] == "" && iblock.WCoordinate[3] == "" {
						iblock.Picture = pic
					} else if len(iblock.WCoordinate) == 4 && iblock.WCoordinate[0] != "" && iblock.WCoordinate[1] != "" && iblock.WCoordinate[2] != "" && iblock.WCoordinate[3] != "" {
						iblock.Picture = bill.BillName + ".crop." + strconv.Itoa(ii) + "." + strconv.Itoa(gg) + ".png"
						// cmd := "gmic -v - src -crop[-1] x0%,y0%,x1%,y1% -o[-1] dst"
						src := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + pic
						//解密切图前的原图
						//decrypt(bill,iblock)

						dst := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + iblock.Picture
						cmd := fmt.Sprintf(`gmic -v - %s -crop[-1] %s%,%s%,%s%,%s% -o[-1] %s`, src, iblock.WCoordinate[0], iblock.WCoordinate[1], iblock.WCoordinate[2], iblock.WCoordinate[3], dst)
						fmt.Println("cmd ", cmd)
						err, stdout, _ := utils.ShellOut(cmd)
						if err != nil {
							fmt.Println("gmic  error", err, stdout)
						}

						//加密分块图片
						//crypto := utils.CreateRandomString(16)
						//iblock.Crypto = crypto
						//err = encImage(bill, crypto, iblock.Picture)
						//if err != nil {
						//	fmt.Println("encrypt block image error", err, stdout)
						//}
					}
					// iblock.PicPage = bill.Pictures[ page ]
				} else {
					continue
				}
				if len(bFields) == 0 {
					continue
				}
				iblock.Level = block.Level
				if iblock.Code == "bc022" || iblock.Code == "bc023" || iblock.Code == "bc024" || iblock.Code == "bc025" || iblock.Code == "bc026" {
					if utils.HasItem(fc040s, fc040) {
						iblock.Op1Stage = "done"
						iblock.Op1Code = "0"
						iblock.Op1ApplyAt = time.Now()
						iblock.Op1SubmitAt = time.Now()
						iblock.Op2Stage = "done"
						iblock.Op2Code = "0"
						iblock.Op2ApplyAt = time.Now()
						iblock.Op2SubmitAt = time.Now()
						iblock.OpqStage = "done"
						iblock.OpqCode = "0"
						iblock.OpqApplyAt = time.Now()
						iblock.OpqSubmitAt = time.Now()
					}

				}

				err = tx.Model(&model.ProjectBlock{}).Create(&iblock).Error
				if err != nil {
					return err
				}
				sFields, isOp2 := initFields(bFields, iblock, ii, CacheFieldConf)
				if !isOp2 && iblock.Op2Stage == "op2" && !iblock.IsCompetitive {
					// iblock.Op2Stage = ""
					iblock.Op2Stage = "done"
					iblock.Op2Code = "0"
					iblock.Op2ApplyAt = time.Now()
					iblock.Op2SubmitAt = time.Now()
					err = tx.Model(&model.ProjectBlock{}).Where("id = ?", iblock.ID).Updates(iblock).Error
					if err != nil {
						return err
					}
				}
				err = tx.Create(&sFields).Error
				if err != nil {
					return err
				}
			}
			if fc040 != "" {
				fc040s = append(fc040s, fc040)
			}

		}
		block.Op0Stage = "done"
		err = tx.Model(&model.ProjectBlock{}).Where("id = ?", block.ID).Updates(block).Error
		if err != nil {
			return err
		}

		//加密单据原图和缩列图
		//err, bill = Encrypt(bill)
		//err = tx.Model(&model2.ProjectBill{}).Where("id = ?", bill.ID).Updates(bill).Error
		return nil
	})
}

func getValue(fFields []model.ProjectField, name string, CacheFieldConf map[string]pModel.SysProField) string {
	fmt.Println("getValuegetValue ", name)
	for _, field := range fFields {
		if field.Name == name {
			if name == "模板类型字段" {
				cacheValue, ok := CacheFieldConf[field.Code]
				if ok {
					return valChange(field.Op0Value, cacheValue.ValChange)
				} else {
					return field.Op0Value
				}
			} else {
				return field.Op0Value
			}

			// if field.ResultValue != "" {
			// 	fmt.Println("ResultValue ", field.ResultValue)
			// 	return field.ResultValue
			// } else {
			// 	fmt.Println("Op0Value ", field.Op0Value)
			// 	return field.Op0Value
			// }

		}
	}
	return ""
}

func getCodeValue(fFields []model.ProjectField, code string) string {
	for _, field := range fFields {
		if field.Code == code {
			return field.Op0Value
		}
	}
	return ""
}

func getBlocks(tempBlocks []pModel.SysProTempB, value string) []pModel.SysProTempB {
	blocks := []pModel.SysProTempB{}
	for _, block := range tempBlocks {
		if block.Relation == "MB001-"+value {
			blocks = append(blocks, block)
		}
	}
	return blocks
}

func initBlock(block pModel.SysProTempB, billID string) model.ProjectBlock {
	projectBlock := model.ProjectBlock{}
	projectBlock.Temp = "MB002"
	projectBlock.Code = block.Code
	projectBlock.Name = block.Name
	projectBlock.BillID = billID
	projectBlock.FEight = block.FEight
	projectBlock.Ocr = block.Ocr
	projectBlock.FreeTime = block.FreeTime
	projectBlock.IsLoop = block.IsLoop
	projectBlock.PicPage = block.PicPage
	projectBlock.WCoordinate = block.WCoordinate
	projectBlock.MCoordinate = block.MCoordinate
	// projectBlock.Stage = "opCache"
	projectBlock.Status = 2
	projectBlock.Op1Stage = "opCache"
	projectBlock.Op2Stage = ""
	projectBlock.OpqStage = ""
	projectBlock.IsCompetitive = block.IsCompetitive
	// projectBlock.Mark = ""
	projectBlock.PreBCode = []string{}
	projectBlock.LinkBCode = []string{}
	projectBlock.Op1PreBCode = []string{}
	_, tempBlockRelationList := service.GetBlockRelationsByBId(block.ID)
	for _, tempBlock := range tempBlockRelationList {
		if tempBlock.MyType == 0 {
			projectBlock.PreBCode = append(projectBlock.PreBCode, tempBlock.PreBCode)
		} else if tempBlock.MyType == 1 {
			projectBlock.LinkBCode = append(projectBlock.LinkBCode, tempBlock.PreBCode)
		} else if tempBlock.MyType == 2 {
			projectBlock.Op1PreBCode = append(projectBlock.Op1PreBCode, tempBlock.PreBCode)
		}
	}
	// if len(projectBlock.PreBCode) == 0 && len(projectBlock.Op1PreBCode) == 0 {
	if block.IsCompetitive {
		// projectBlock.Mark = "1"
		projectBlock.Op1Stage = "op1"
		projectBlock.Op2Stage = "op2"
	} else if block.IsLoop {
		if len(projectBlock.PreBCode) == 0 && len(projectBlock.Op1PreBCode) == 0 {
			projectBlock.Op1Stage = "op1"
		}
		projectBlock.Op2Stage = "done"
		projectBlock.Op2Code = "0"
		projectBlock.Op2ApplyAt = time.Now()
		projectBlock.Op2SubmitAt = time.Now()
	} else {
		if len(projectBlock.PreBCode) == 0 {
			projectBlock.Op1Stage = "op1"
			if len(projectBlock.Op1PreBCode) == 0 {
				projectBlock.Op2Stage = "op2"
			}
		}
		// projectBlock.Mark = "2"
	}
	// else if
	// }

	return projectBlock
}

func initFields(bFields []pModel.TempBFRelation, block model.ProjectBlock, blockIndex int, CacheFieldConf map[string]pModel.SysProField) ([]model.ProjectField, bool) {
	fields := []model.ProjectField{}
	isOp2 := false
	for ff, bField := range bFields {
		proField := model.ProjectField{}
		proField.BillID = block.BillID
		proField.BlockID = block.ID
		proField.Code = bField.FCode
		proField.Name = bField.FName
		proField.BlockIndex = blockIndex
		proField.FieldIndex = ff
		value, ok := CacheFieldConf[proField.Code]
		if ok {
			if value.InputProcess == 2 {
				proField.Op1Input = "yes"
				proField.Op2Input = "no"
			} else if value.InputProcess == 3 {
				proField.Op1Input = "yes"
				proField.Op2Input = "yes"
				isOp2 = true
			} else {
				proField.Op1Input = "no"
				proField.Op2Input = "no"
			}
		}
		proField.OpqInput = "no"
		fields = append(fields, proField)
	}
	return fields, isOp2
}

func valChange(value string, valChange string) string {
	if valChange == "" {
		return value
	}
	arr := strings.Split(valChange, ";")
	for _, str := range arr {
		if str == "" {
			continue
		}
		end := strings.Index(str, "=")
		if end == -1 {
			continue
		}
		if str[:end] != value {
			continue
		}
		value = str[end+1:]
		break
	}
	return value
}

// proCustomizeValidation 项目自定义校验需求
func proCustomizeValidation(fields []model.ProjectField, proCode string) error {
	//（fc153为模板类型字段）
	//当fc153录入值同时不存在5、6、7、8时，“案件列表-理赔类型”一栏处显示“无发票”；
	//当fc153录入值存在7，不存在8时，“案件列表-理赔类型”一栏处显示“无报销”；
	//当fc153录入值存在8，不存在7时，“案件列表-理赔类型”一栏处显示“有报销”；
	//当fc153录入值同时存在7和8时，“案件列表-理赔类型”一栏处显示“混合型”；
	//不符合以上4种情况，“理赔类型”的值为空；
	if len(fields) < 1 {
		return errors.New("没有初审字段")
	}
	billId := fields[0].BillID
	val := getFieldsValueByCode(fields, "fc056")
	// 1: "医疗",
	// 2: "非医疗",
	// 3: "无发票",
	// 4: "无报销",
	// 5: "有报销",
	// 6: "混合型",
	// 7: "简易",
	var claimType = -1
	if !utils.IsContain(val, "5") && !utils.IsContain(val, "6") && !utils.IsContain(val, "7") {
		claimType = 3
	}
	if utils.IsContain(val, "7") && !utils.IsContain(val, "8") {
		claimType = 4
	}
	if utils.IsContain(val, "8") && !utils.IsContain(val, "7") {
		claimType = 5
	}
	if utils.IsContain(val, "7") && utils.IsContain(val, "8") {
		claimType = 6
	}
	return service.UpdateClaimType(proCode, billId, claimType)
}

// getFieldsValueByCode 在字段数组中根据code获取的值
func getFieldsValueByCode(fields []model.ProjectField, code string) (val []string) {
	val = make([]string, 0)
	for _, field := range fields {
		if field.Code == code {
			val = append(val, field.ResultValue)
		}
	}
	return val
}
