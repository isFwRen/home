package B0108

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
	model2 "server/module/pro_manager/model"
	"server/utils"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/wxnacy/wgo/arrays"
	"go.uber.org/zap"

	"gorm.io/gorm"
)

func TypeCrop(proCode string, block model.ProjectBlock, tempBlocks []pModel.SysProTempB, CacheFieldConf map[string]pModel.SysProField) error {
	proID := global.ProCodeId[proCode]
	global.GLog.Info("TypeCrop", zap.Any("proID", proID))
	err, bill := service.SelectBillByID(proCode+"_task", block.BillID)
	if err != nil {
		global.GLog.Error("TypeCrop-SelectBillByID", zap.Error(err))
		return err
	}
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
	c, ok := global.GProConf[proCode].ConstTable["B0108_太平理赔_全国"]
	tempMap := []string{}
	if ok {
		for _, arr := range c {
			tempMap = append(tempMap, strings.TrimSpace(arr[0]))
			// tempMap[strings.TrimSpace(arr[0])] = arr[2]
		}
	}
	opFields := _func.FieldsFormat(fields)
	bcDis := true
	fc084s := []string{}
	for _, fFields := range opFields {
		for _, field := range fFields {
			if field.Code == "fc084" {
				fc084s = append(fc084s, field.ResultValue)
				matched, _ := regexp.MatchString(`(5|6)`, field.ResultValue)
				if matched {
					bcDis = false
				}
			}
		}
	}
	// bc007 := false
	// bc008 := false
	bcArrs := []string{}
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
				iBlock := initBlock(gBlock, block.BillID, bill)
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
				iBlock.Level = block.Level
				if ((iBlock.Code == "bc005" || iBlock.Code == "bc007" || iBlock.Code == "bc008") && bcDis) || ((iBlock.Code == "bc007" || iBlock.Code == "bc008") && !utils.HasItem(fc084s, "6")) {
					// iBlock.Op1Stage = "done"
					// iBlock.Op1Code = "0"
					// iBlock.Op1ApplyAt = time.Now()
					// iBlock.Op1SubmitAt = time.Now()
					// iBlock.Op2Stage = "done"
					// iBlock.Op2Code = "0"
					// iBlock.Op2ApplyAt = time.Now()
					// iBlock.Op2SubmitAt = time.Now()
					// iBlock.OpqStage = "done"
					// iBlock.OpqCode = "0"
					// iBlock.OpqApplyAt = time.Now()
					// iBlock.OpqSubmitAt = time.Now()
					iBlock = disableBlock(iBlock)
				}
				if iBlock.Code == "bc007" || iBlock.Code == "bc008" {
					if arrays.Contains(bcArrs, iBlock.Code) != -1 {
						iBlock = disableBlock(iBlock)
					} else {
						bcArrs = append(bcArrs, iBlock.Code)
					}
				}

				err = tx.Model(&model.ProjectBlock{}).Create(&iBlock).Error
				if err != nil {
					return err
				}
				sFields, isOp2 := initFields(bFields, iBlock, 0, CacheFieldConf)
				if iBlock.Code == "bc004" {
					_, fc084 := service.SelectBillFields(proCode+"_task", iBlock.BillID, iBlock.Zero, "fc084")
					_, fc218 := service.SelectBillFields(proCode+"_task", iBlock.BillID, iBlock.Zero, "fc218")
					fmt.Println("------------------------bc004-------------------------------", fc084, fc218)
					if fc084.ResultValue == "7" && fc218.ResultValue == "1" {
						dst := global.GConfig.LocalUpload.FilePath + bill.DownloadPath + iBlock.Picture
						sFields = ocr_fields(block, sFields, dst, tempMap)
					}
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

func disableBlock(iBlock model.ProjectBlock) model.ProjectBlock {
	iBlock.Op1Stage = "done"
	iBlock.Op1Code = "0"
	iBlock.Op1ApplyAt = time.Now()
	iBlock.Op1SubmitAt = time.Now()
	iBlock.Op2Stage = "done"
	iBlock.Op2Code = "0"
	iBlock.Op2ApplyAt = time.Now()
	iBlock.Op2SubmitAt = time.Now()
	iBlock.OpqStage = "done"
	iBlock.OpqCode = "0"
	iBlock.OpqApplyAt = time.Now()
	iBlock.OpqSubmitAt = time.Now()
	return iBlock
}

func ocr_fields(block model.ProjectBlock, bFields []model.ProjectField, pic string, tempMap []string) []model.ProjectField {
	sFields := []string{"fc080", "fc143", "fc145", "fc147", "fc149", "fc151", "fc153", "fc155"}
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
		if bb+1 < len(block_datas) && RegIsMatch(check_ocr_value(block_datas[bb+1].Text), `^[_\(\（\）\)]`) {
			text = text + check_ocr_value(block_datas[bb+1].Text)
		}
		datas = append(datas, text)
	}
	fmt.Println("-------datas------", datas)
	cloneFields := make([]model.ProjectField, len(bFields))
	cacheFields := make([]model.ProjectField, len(bFields))
	copy(cloneFields, bFields)
	copy(cacheFields, bFields)
	bFields = []model.ProjectField{}
	fIndex := 0
	bIndex := 0
	for dd, data := range datas {
		data = utils.RegReplace(data, `[\*<>（）\(\)\-—★。●•▪※【】\[\]《》_]`, "")
		fmt.Println("-------------------da-ta-------------------", data, dd, fIndex, bIndex, len(bFields), len(cloneFields))
		if fIndex == 7 {
			if arrays.ContainsString(tempMap, data) != -1 {
				cloneFields = set_ocr_fields(cloneFields, sFields[fIndex], data, bIndex, true)
			} else {
				value := searchValue(data, tempMap)
				cloneFields = set_ocr_fields(cloneFields, sFields[fIndex], value, bIndex, false)
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
			if arrays.ContainsString(tempMap, data) != -1 {
				cloneFields = set_ocr_fields(cloneFields, sFields[fIndex], data, bIndex, true)
			} else {
				value := searchValue(data, tempMap)
				cloneFields = set_ocr_fields(cloneFields, sFields[fIndex], value, bIndex, false)
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

func initBlock(block pModel.SysProTempB, billID string, bill model2.ProjectBill) model.ProjectBlock {
	projectBlock := model.ProjectBlock{}
	projectBlock.CreatedAt = bill.ScanAt
	if bill.SaleChannel == "理赔" {
		projectBlock.CreatedAt = bill.ScanAt.Add(20 * time.Minute)
	}
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
	//CSB0108RC0009000
	//当fc084录入值同时不存在5、6、7、8时（fc084为循环字段），“案件列表-理赔类型”一栏处显示“无发票”；
	//当fc084录入值存在7，不存在8时（fc084为循环字段），“案件列表-理赔类型”一栏处显示“无报销”；
	//当fc084录入值存在8，不存在7时（fc084为循环字段），“案件列表-理赔类型”一栏处显示“有报销”；
	//当fc084录入值同时存在7和8时（fc084为循环字段），“案件列表-理赔类型”一栏处显示“混合型”；
	//不符合以上4种情况，“理赔类型”的值为空；
	if len(fields) < 1 {
		return errors.New("没有初审字段")
	}
	billId := fields[0].BillID
	val := getFieldsValueByCode(fields, "fc084")
	global.GLog.Info("项目自定义校验需求", zap.Any("fc084", val))
	var claimType = -1
	if !utils.IsContain(val, "5") && !utils.IsContain(val, "6") &&
		!utils.IsContain(val, "7") && !utils.IsContain(val, "8") {
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

func searchValue(baseStr string, compareStrs []string) (res string) {
	ii := 100
	value := ""
	for _, str := range compareStrs {
		aa := fuzzy.LevenshteinDistance(baseStr, str)
		// fmt.Println("-------------", str, aa)
		if aa < ii {
			ii = aa
			value = str
		}
	}
	if utf8.RuneCountInString(baseStr) >= utf8.RuneCountInString(value) {
		ii = 0
		len := utf8.RuneCountInString(value)
		for i, str := range []rune(baseStr) {

			if i >= len {
				res += "?"
			} else {
				if string(str) == substr(value, i, i+1) {
					ii++
					res += string(str)
				} else {
					res += "?"
				}
			}
		}
		if ii < 4 {
			res = ""
		}
	} else {
		ii = 0
		len := utf8.RuneCountInString(baseStr)
		for i, str := range []rune(value) {
			if i >= len {
				res += "?"
			} else {
				if string(str) == substr(baseStr, i, i+1) {
					ii++
					res += string(str)
				} else {
					res += "?"
				}
			}
		}
		if ii < 4 {
			res = ""
		}
	}

	return res
}

func substr(str string, start, end int) string {
	data := []rune(str)
	if start == -1 && end != -1 {
		return string(data[:end])
	} else if start != -1 && end == -1 {
		return string(data[start:])
	} else if start != -1 && end != -1 {
		return string(data[start:end])
	}
	return str

}
