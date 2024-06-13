package B0103

import (
	"errors"
	"fmt"
	"server/global"
	_func "server/global/func"
	"server/module/load/model"
	"server/module/load/service"
	pModel "server/module/pro_conf/model"
	unitFunc "server/module/unit"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"gorm.io/gorm"
)

func TypeCrop(proCode string, block model.ProjectBlock, tempBlocks []pModel.SysProTempB, CacheFieldConf map[string]pModel.SysProField) error {

	proID := global.ProCodeId[proCode]
	global.GLog.Info("TypeCrop", zap.Any("proID", proID))
	err, fields := service.SelectFieldsByBlockID(proCode+"_task", block.ID)
	if err != nil {
		global.GLog.Error("TypeCrop-SelectFieldsByBlockID", zap.Error(err))
		return err
	}
	//项目自定义校验需求
	err = proCustomizeValidation(fields, proCode)
	if err != nil {
		global.GLog.Error("TypeCrop-proCustomizeValidation", zap.Error(err))
		return err
	}
	err, bill := service.SelectBillByID(proCode+"_task", block.BillID)
	if err != nil {
		global.GLog.Error("TypeCrop-SelectBillByID", zap.Error(err))
		return err
	}
	opFields := _func.FieldsFormat(fields)

	fc091 := ""
	for _, fFields := range opFields {
		for _, field := range fFields {
			if field.Code == "fc091" {
				fc091 = field.ResultValue
			}
		}
	}
	fmt.Println("----------fc091-----------", fc091)

	return global.ProDbMap[proCode+"_task"].Transaction(func(tx *gorm.DB) error {
		for ii, fFields := range opFields {
			value := getValue(fFields, "模板类型字段", CacheFieldConf)
			page := getValue(fFields, "图片页码", CacheFieldConf)
			w := getValue(fFields, "显示范围", CacheFieldConf)
			global.GLog.Info("", zap.Any("模板类型字段", value))
			global.GLog.Info("", zap.Any("图片页码", page))
			global.GLog.Info("", zap.Any("显示范围", w))
			blocks := getBlocks(tempBlocks, value)
			for gg, gBlock := range blocks {
				err, bFields := service.GetTempBFRelationByBId(gBlock.ID)
				global.GLog.Info("GetTempBFRelationByBId", zap.Any("len(bFields)", len(bFields)))
				if err != nil {
					return err
				}
				iBlock := initBlock(gBlock, block.BillID)
				iBlock.CreatedAt = bill.ScanAt
				iBlock.Zero = ii
				if w != "" {
					iBlock.Picture = w
				} else if page != "" {
					idx, _ := strconv.Atoi(page)
					pic := bill.Pictures[idx]
					fmt.Println("------------------------------------------------- ", iBlock.Code, iBlock.WCoordinate)
					if iBlock.WCoordinate == nil || len(iBlock.WCoordinate) == 0 {
						iBlock.Picture = pic
					} else if len(iBlock.WCoordinate) == 4 && iBlock.WCoordinate[0] == "" && iBlock.WCoordinate[1] == "" && iBlock.WCoordinate[2] == "" && iBlock.WCoordinate[3] == "" {
						iBlock.Picture = pic
					} else if len(iBlock.WCoordinate) == 4 && iBlock.WCoordinate[0] != "" && iBlock.WCoordinate[1] != "" && iBlock.WCoordinate[2] != "" && iBlock.WCoordinate[3] != "" {
						iBlock.Picture = bill.BillName + ".crop." + strconv.Itoa(ii) + "." + strconv.Itoa(gg) + ".png"
						// cmd := "gmic -v - src -crop[-1] x0%,y0%,x1%,y1% -o[-1] dst"
						src := bill.DownloadPath + pic
						//解密切图前的原图
						//decrypt(bill,iBlock)

						dst := bill.DownloadPath + iBlock.Picture
						cmd := fmt.Sprintf(`gmic -v - %s -crop[-1] %s%%,%s%%,%s%%,%s%% -o[-1] %s`, src, iBlock.WCoordinate[0], iBlock.WCoordinate[1], iBlock.WCoordinate[2], iBlock.WCoordinate[3], dst)
						fmt.Println("-----cmd--------", cmd)
						err, stdout, _ := utils.ShellOut(cmd)
						if err != nil {
							global.GLog.Error("gmic", zap.Error(err))
							global.GLog.Error(stdout)
						}

						//加密分块图片
						//crypto := utils.CreateRandomString(16)
						//iBlock.Crypto = crypto
						//err = encImage(bill, crypto, iBlock.Picture)
						//if err != nil {
						//	fmt.Println("encrypt block image error", err, stdout)
						//}
					}
					// iBlock.PicPage = bill.Pictures[ page ]
				} else {
					continue
				}
				if len(bFields) == 0 {
					continue
				}
				iBlock.Level = 99
				//二轮加载时，类型为“有报销”或“空”的案件，自动将该单号放入优先单池
				if bill.ClaimType == 5 || bill.ClaimType == 0 {
					iBlock.Level = 2
				}
				err = tx.Model(&model.ProjectBlock{}).Create(&iBlock).Error
				if err != nil {
					return err
				}
				sFields, isOp2 := initFields(bFields, iBlock, 0, CacheFieldConf)

				if iBlock.Code == "bc001" && fc091 == "1" {
					// fc067 == 1
					dst := bill.DownloadPath + iBlock.Picture
					fmt.Println("-------dst----------", dst)
					err, qrcode := unitFunc.OcrZxing(dst)
					if err == nil {
						bodyData := make(map[string]interface{})
						bodyData["code"] = qrcode
						err, respData := unitFunc.Invoice("/v2/invoice/qrcode", bodyData)
						fc097 := ""
						fmt.Println("--------respData--------------", err, respData)
						if err != nil {
							fmt.Println("--------respData--------------", err, respData)
							// response.FailWithMessage(fmt.Sprintf("%v", err), c)
							// return
						} else {
							del := respData.Data.Del
							if del == "0" {
								fc097 = "01"
							}
							if del == "2" {
								fc097 = "99"
							}
							if del == "3" || del == "7" || del == "8" {
								fc097 = "02"
							}
						}
						sFields = set_fields(sFields, "fc097", fc097, 0)
					}
					// return errors.New("test")
				}

				if !isOp2 && iBlock.Op2Stage == "op2" && !iBlock.IsCompetitive {
					// iBlock.Op2Stage = ""
					iBlock.Op2Stage = "done"
					iBlock.Op2Code = "0"
					iBlock.Op2ApplyAt = time.Now()
					iBlock.Op2SubmitAt = time.Now()
					err = tx.Model(&model.ProjectBlock{}).Where("id = ?", iBlock.ID).Updates(iBlock).Error
					if err != nil {
						return err
					}
				}
				err = tx.Create(&sFields).Error
				if err != nil {
					return err
				}
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
	global.GLog.Info("getValue", zap.Any("name", name))
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

func getBlocks(tempBlocks []pModel.SysProTempB, value string) []pModel.SysProTempB {
	var blocks []pModel.SysProTempB
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
		projectBlock.Op1Stage = "opCache"
		projectBlock.Op2Stage = "opCache"
	} else if block.IsLoop {
		if len(projectBlock.PreBCode) == 0 && len(projectBlock.Op1PreBCode) == 0 {
			projectBlock.Op1Stage = "opCache"
		}
		projectBlock.Op2Stage = "done"
		projectBlock.Op2Code = "0"
		projectBlock.Op2ApplyAt = time.Now()
		projectBlock.Op2SubmitAt = time.Now()
	} else {
		if len(projectBlock.PreBCode) == 0 {
			projectBlock.Op1Stage = "opCache"
			if len(projectBlock.Op1PreBCode) == 0 {
				projectBlock.Op2Stage = "opCache"
			}
		}
	}

	return projectBlock
}

func initFields(bFields []pModel.TempBFRelation, block model.ProjectBlock, blockIndex int, CacheFieldConf map[string]pModel.SysProField) ([]model.ProjectField, bool) {
	var fields []model.ProjectField
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
	//CSB0106RC0082000
	//同一个案件中，业务清单的“理赔类型”这一列的值根据以下条件赋值:
	//a.fc003(循环)录入值不包含3时，“理赔类型”的值为“无报销”；
	//b.fc003(循环)录入值不包含2时，“理赔类型”的值为“有报销”；
	//c.fc003(循环)录入值同时包含2和3时，“理赔类型”的值为“混合型”；
	//d.不符合以上3点条件的，“理赔类型”的值为空；
	//修改2023年09月12日11:12:31黄文雅
	//d.fc003(循环)录入值不包含1时，“理赔类型”的值为空；
	if len(fields) < 1 {
		return errors.New("没有初审字段")
	}
	billId := fields[0].BillID
	val := getFieldsValueByCode(fields, "fc003")
	var claimType = -1
	if !utils.IsContain(val, "3") {
		claimType = 4
	}
	if !utils.IsContain(val, "2") {
		claimType = 5
	}
	if utils.IsContain(val, "3") && utils.IsContain(val, "2") {
		claimType = 6
	}
	if !utils.IsContain(val, "1") {
		claimType = -1
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

func set_fields(fFields []model.ProjectField, code string, value string, bb int) []model.ProjectField {
	for ff, field := range fFields {
		if field.Code == code {
			fFields[ff].Op1Value = value
			fFields[ff].Op1Input = "no_if"
			fFields[ff].Op2Value = value
			fFields[ff].Op2Input = "no_if"
			fFields[ff].ResultValue = value
			fFields[ff].ResultInput = "no_if"
		}
		fFields[ff].BlockIndex = bb
	}
	return fFields
}
