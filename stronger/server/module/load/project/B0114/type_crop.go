package B0114

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"server/global"
	_func "server/global/func"
	"server/module/load/model"
	"server/module/load/service"
	pModel "server/module/pro_conf/model"
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
	constMap := constDeal(proCode)

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
						src := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + pic
						//解密切图前的原图
						//decrypt(bill,iBlock)

						dst := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + iBlock.Picture
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
				fmt.Println("------------------------bc004-------------------------------", iBlock.Code, value, getValue(fFields, "是否OCR", CacheFieldConf))
				if iBlock.Code == "bc010" && value == "bc004" && getValue(fFields, "是否OCR", CacheFieldConf) == "1" {
					_, fc008 := service.SelectBillFields(proCode+"_task", iBlock.BillID, -1, "fc008")
					_, fc003 := service.SelectBillFields(proCode+"_task", iBlock.BillID, -1, "fc003")
					_, fc004 := service.SelectBillFields(proCode+"_task", iBlock.BillID, -1, "fc004")
					fmt.Println("------------------------fc008-------------------------------", fc008.ResultValue, fc003.ResultValue, fc004.ResultValue)
					// if fc084.ResultValue == "7" && fc218.ResultValue == "1" {
					policyBranchCode := utils.GetNodeValue(bill.OtherInfo, "policyBranchCode")
					yiLiaoMuLu, yiLiaoMuLuInPolicyBranchCode := GetYiLiaoMuLu(proCode, policyBranchCode, fc008.ResultValue, fc003.ResultValue, fc004.ResultValue, constMap)
					fmt.Println("------------------------yiLiaoMuLu-------------------------------", policyBranchCode, yiLiaoMuLu, yiLiaoMuLuInPolicyBranchCode)
					dst := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + iBlock.Picture
					sFields = ocr_fields(block, sFields, dst, yiLiaoMuLu, yiLiaoMuLuInPolicyBranchCode, constMap)
					// }
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

				fmt.Println("------------------------------LLLLLLLLLLLLL=")
				fmt.Println("sFields=", sFields)
				fmt.Println("block=", block)
				fmt.Println("------------------------------LLLLLLLLLLLLL=")

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
		// return errors.New("测试")
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

	//fc179（初审循环）录入内容只包含4，不包含5时（可以有1、2等内容），业务清单中的类型为“无报销”
	//fc179（初审循环）录入内容只包含5，不包含4时（可以有1、2等内容），业务清单中的类型为“有报销”
	//fc179（初审循环）录入内容即包含4又包含5时（可以有1、2等内容），业务清单中的类型为“混合型”
	//4.不符合以上3点条件的，“理赔类型”的值为空；
	if len(fields) < 1 {
		return errors.New("没有初审字段")
	}
	billId := fields[0].BillID
	val := getFieldsValueByCode(fields, "fc179")
	var claimType = -1
	if utils.IsContain(val, "4") && !utils.IsContain(val, "5") {
		claimType = 4
	}
	if utils.IsContain(val, "5") && !utils.IsContain(val, "4") {
		claimType = 5
	}
	if utils.IsContain(val, "5") && utils.IsContain(val, "4") {
		claimType = 6
	}
	stickLevel := 99
	//二轮加载时，类型为“有报销”或“空”的案件，自动将该单号放入优先单池
	if claimType == 5 || claimType == 0 {
		stickLevel = 2
	}
	err := service.UpdateClaimTypeStickLevel(proCode, billId, claimType, stickLevel)
	if err != nil {
		return err
	}
	return service.UpdateClaimTypeStickLevel(proCode+"_task", billId, claimType, stickLevel)
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

func ocr_fields(block model.ProjectBlock, bFields []model.ProjectField, pic string, yiLiaoMuLu, yiLiaoMuLuInPolicyBranchCode string, constMap map[string]map[string]string) []model.ProjectField {
	sFields := []string{"fc138", "fc142", "fc146", "fc150", "fc154", "fc158", "fc162", "fc166"}
	cmd := fmt.Sprintf(`curl -X POST -F image=@%s http://192.168.202.17:18003/ocr`, pic)
	err, stdout, _ := utils.ShellOut(cmd)
	fmt.Println("-------ocr_fields------", err, cmd)
	if err != nil {
		global.GLog.Error("ocr", zap.Error(err))
		global.GLog.Error(stdout)
		return bFields
	}
	var ocrData OcrData
	datas := []string{}
	json.Unmarshal([]byte(stdout), &ocrData)
	if len(ocrData.Images) == 0 {
		return bFields
	}
	if len(ocrData.Images[0].Pages) == 0 {
		return bFields
	}
	block_datas := ocrData.Images[0].Pages[0].Blocks
	// fmt.Println("-------block_datas------", block_datas)
	for bb, block_data := range block_datas {
		if block_data.N != 0 {
			continue
		}
		x := (block_data.Area[0] - block_data.Area[2]/2) - (block_datas[0].Area[0] - block_datas[0].Area[2]/2)
		if x > 0.1 || x < -0.1 {
			continue
		}
		text := check_ocr_value(block_data.Text)
		if bb+1 < len(block_datas) && RegIsMatch(check_ocr_value(block_datas[bb+1].Text), `^_`) {
			text = text + check_ocr_value(block_datas[bb+1].Text)
		}
		datas = append(datas, text)
	}
	fmt.Println("-----------------------LLLLLLLLLLLLLLLLLLLLLLLLLL---------------------------------------------")
	fmt.Println("---------------------ocrData=", ocrData)
	fmt.Println("*******************************************************************")
	fmt.Println("---------------------datas=", datas)
	fmt.Println("*******************************************************************")
	fmt.Println("---------------------block_datas=", block_datas)
	fmt.Println("-----------------------LLLLLLLLLLLLLLLLLLLLLLLLLL---------------------------------------------")
	fmt.Println("-------datas------", datas)
	cloneFields := make([]model.ProjectField, len(bFields))
	cacheFields := make([]model.ProjectField, len(bFields))
	copy(cloneFields, bFields)
	copy(cacheFields, bFields)
	bFields = []model.ProjectField{}
	fIndex := 0
	bIndex := 0
	for dd, data := range datas {
		fmt.Println("-------------------da-ta-------------------", data, dd, fIndex, bIndex, len(bFields), len(cloneFields))
		_, is1 := constMap[yiLiaoMuLuInPolicyBranchCode][data]
		_, is2 := constMap[yiLiaoMuLu][data]
		if fIndex == 7 {
			if is1 || is2 {
				cloneFields = set_ocr_fields(cloneFields, sFields[fIndex], data, bIndex, true)
			} else {
				cloneFields = set_ocr_fields(cloneFields, sFields[fIndex], data, bIndex, false)
			}
			bFields = append(bFields, cloneFields...)
			// cloneFields = cacheFields
			bIndex++
			copy(cloneFields, cacheFields)
			for ff, _ := range cloneFields {
				cloneFields[ff].BlockIndex = bIndex
			}
			fIndex = 0
		} else {
			if is1 || is2 {
				cloneFields = set_ocr_fields(cloneFields, sFields[fIndex], data, bIndex, true)
			} else {
				cloneFields = set_ocr_fields(cloneFields, sFields[fIndex], data, bIndex, false)
			}
			if dd == len(datas)-1 {
				bFields = append(bFields, cloneFields...)
			}
			fIndex++
		}
	}
	if len(bFields) == 0 {
		bFields = cacheFields
	}
	// stdout, _ = strconv.Unquote(stdout)
	// json.Unmarshal([]byte(stdout), &ocrData)

	return bFields
}

func set_ocr_fields(fFields []model.ProjectField, code string, value string, bb int, isOCR bool) []model.ProjectField {
	fmt.Println("-------------------set_ocr_fields-------------------", code, value)
	for ff, field := range fFields {
		if field.Code == code {
			fFields[ff].Op0Value = value
			fFields[ff].Op0Input = "ocr"
			fFields[ff].Op1Value = value
			fFields[ff].Op2Value = value
			fFields[ff].ResultValue = value
			if isOCR {
				fFields[ff].Op0Input = "ocr_no_if"
				fFields[ff].Op1Input = "no_if"
				fFields[ff].Op2Input = "no_if"
				fFields[ff].ResultInput = "no_if"
			}

		}
		fFields[ff].BlockIndex = bb
	}
	return fFields
}

func check_ocr_value(value string) string {
	if RegIsMatch(value, `(合计|小计|以上|汇总|护理)`) {
		return ""
	}
	if RegIsMatch(value, `(\*)`) {
		return RegReplace(value, `\*`, "")
	}
	return value
}

func RegReplace(data string, query string, value string) string {
	reg := regexp.MustCompile(query)
	return reg.ReplaceAllString(data, value)
}

func HasKey(data map[string]string, key string) bool {
	_, isOK := data[key]
	return isOK
}

func RegIsMatch(value string, query string) bool {
	// reg := regexp.MustCompile(query)
	matched, _ := regexp.MatchString(query, value)
	return matched
}

func GetYiLiaoMuLu(ProCode, policyBranchCode, fc008, fc003, fc004 string, constMap map[string]map[string]string) (string, string) {
	// constMap := constDeal(ProCode)
	specialConstMap := specialConst(ProCode)
	//获取医疗目录
	// policyBranchCode := utils.GetNodeValue(obj.Bill.OtherInfo, "policyBranchCode")
	yiLiaoMuLu := ""
	yiLiaoMuLuInPolicyBranchCode := ""
	valArr := specialConstMap["diZhiMap"][fc008]
	if len(valArr) == 1 {
		yiLiaoMuLu = valArr[0]
	} else if len(valArr) > 1 {
		str := fc008
		if fc003 != "" {
			str = fc003 + "_" + str
		}
		if fc004 != "" {
			str = fc004 + "_" + str
		}
		val, ok := constMap["diZhiShengShiQuMap"][str]
		if ok {
			yiLiaoMuLu = val
		}
	} else {
		valArr, ok := specialConstMap["diZhiShiMap"][fc003]
		if ok {
			yiLiaoMuLu = valArr[0]
		} else {
			valArr, ok = specialConstMap["diZhiShengMap"][fc003]
			if ok {
				yiLiaoMuLu = valArr[0]
			}
		}
	}
	val, ok := constMap["jiGouBianMaMap"][policyBranchCode]
	if ok {
		yiLiaoMuLuInPolicyBranchCode = val
	}
	return yiLiaoMuLu, yiLiaoMuLuInPolicyBranchCode
}

// 初始化常量
func constDeal(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"diZhiShengShiQuMap", "B0114_华夏理赔_华夏理赔地址库", "4", "1"},
		{"jiGouBianMaMap", "B0114_华夏理赔_机构编码对应表", "1", "0"},
	}
	for i, item := range nameMap {
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		tempMap := make(map[string]string, 0)
		tempNumMap := make(map[string]string, 0)
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi(item[2])
				v, _ := strconv.Atoi(item[3])

				//省市区匹配唯一CSB0114RC0076000
				if i == 1 {
					tempMap[arr[2]+"_"+arr[3]+"_"+arr[4]] = arr[v]
					tempMap[arr[3]+"_"+arr[4]] = arr[v]
					tempMap[arr[2]+"_"+arr[4]] = arr[v]

					//同一个key有多少个值
					num := 1
					numStr, has := tempNumMap[strings.TrimSpace(arr[k])]
					if has {
						num, _ = strconv.Atoi(numStr)
						num++
					}
					tempNumMap[strings.TrimSpace(arr[k])] = fmt.Sprintf("%d", num)
				}
				//else {
				tempMap[strings.TrimSpace(arr[k])] = arr[v]
				//}
			}
		}
		constObj[item[0]] = tempMap
		constObj[item[0]+"Num"] = tempNumMap
	}

	//地址
	for k, v := range global.GProConf[proCode].ConstTable {
		if strings.Index(k, "省份") != -1 {
			tempMap := make(map[string]string, 0)
			for _, arr := range v {
				if len(arr) < 3 {
					continue
				}
				//获取最小的一项
				if val, ok := tempMap[arr[2]]; ok {
					a, err := strconv.ParseFloat(val, 64)
					if err != nil {
						continue
					}
					b, err := strconv.ParseFloat(arr[0], 64)
					if err != nil {
						continue
					}
					if a < b {
						arr[0] = val
					}
				}
				tempMap[arr[2]] = arr[0]
			}
			constObj[strings.Replace(k, "B0114_华夏理赔_省份-", "", -1)] = tempMap
		}
	}

	return constObj
}

// specialConst 特别的常量 需要判断是否唯一的
func specialConst(proCode string) map[string]map[string][]string {
	constObj := make(map[string]map[string][]string, 0)
	nameMap := [][]string{
		{"diZhiMap", "B0114_华夏理赔_华夏理赔地址库", "4", "0"},
		{"diZhiShiMap", "B0114_华夏理赔_华夏理赔地址库", "3", "0"},
		{"diZhiShengMap", "B0114_华夏理赔_华夏理赔地址库", "2", "0"},
	}
	for _, item := range nameMap {
		tempMap := make(map[string][]string, 0)
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi(item[2])
				v, _ := strconv.Atoi(item[3])
				tempMap[strings.TrimSpace(arr[k])] = append(tempMap[strings.TrimSpace(arr[k])], arr[v])
			}
		} else {
			global.GLog.Error(item[1], zap.Error(errors.New("没有该常量")))
			fmt.Println("------------------------proCode-------------------------------", proCode, global.GProConf[proCode].ConstTable)
		}
		constObj[item[0]] = tempMap
	}
	return constObj
}
