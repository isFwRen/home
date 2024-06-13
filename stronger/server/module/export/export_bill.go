/**
 * @Author: xingqiyi
 * @Description:导出
 * @Date: 2021/11/23 4:03 下午
 */

package export

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"runtime/debug"
	"server/global"
	model4 "server/module/export/model"
	"server/module/export/project"
	"server/module/export/service"
	utils2 "server/module/export/utils"
	model3 "server/module/load/model"
	"server/module/pro_conf/const_data"
	model2 "server/module/pro_conf/model"
	"server/module/pro_manager/model"
	managerService "server/module/pro_manager/service"
	"server/utils"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/flosch/pongo2/v4"
	"go.uber.org/zap"
)

// Process 导出流程
func Process(proCode string) error {
	var wg sync.WaitGroup
	var billList []model.ProjectBill
	var fieldConfList []model2.SysProField

	//1.获取单据 GetTaskBillList
	wg.Add(1)
	go service.GetTaskBillList(proCode, &wg, &billList)

	//2.获取字段配置 getFields
	wg.Add(1)
	go service.GetFieldConf(proCode, &wg, &fieldConfList)

	wg.Wait()
	flag := "2"
	//3.循环bills 每一张都做导出一张单的流程 BillExport
	for _, bill := range billList {
		fmt.Println(bill)
		//3.1获取分块
		err, blocks := service.GetBlocks(bill)
		if err != nil {
			global.GLog.Error("单据：：："+bill.ID+"获取分块错误", zap.Error(err))
			continue
		}
		copyBlocks := utils.CopyBlocksSlice(blocks)

		//获取字段
		err, fields := service.GetTaskBillFields(bill)
		if err != nil {
			global.GLog.Error("单据：：："+bill.ID+"获取字段错误", zap.Error(err))
			continue
		}
		copyFields := utils.CopyFieldsSlice(fields)

		//3.2导出处理流程
		obj, err := BillExport(bill, blocks, fieldConfList, fields, flag)
		if err != nil {
			global.GLog.Error("单据：：："+bill.ID+"导出错误", zap.Error(err))
			continue
		}

		//将分块和字段搬到历史库
		err = service.SaveBlocksAndFields(proCode, copyBlocks, copyFields)
		if err != nil {
			global.GLog.Error("单据：：："+bill.ID+"保存分块和字段到历史库错误", zap.Error(err))
			continue
		}

		//3.3删除任务库的该单据
		err = service.DelBill(bill.ID, bill.ProCode+"_task")
		if err != nil {
			return err
		}

		err = service.AutoExportFieldUpdate(obj)
		if err != nil {
			return err
		}

	}

	//退出
	return nil
}

// BillExport 处理导出保存一张单
func BillExport(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldConfList []model2.SysProField, fields []model3.ProjectField, flag string) (obj model4.ResultDataBill, err error) {
	//生成结果数据
	obj, err = ChangeData(bill, blocks, fieldConfList, fields)
	if err != nil {
		return obj, err
	}

	//3.根据结果值和导出配置生成xml BillToXml
	err, xmlValue, formatObj := BillToXml(bill.ProCode, &obj)
	if err != nil {
		return obj, err
	}
	//global.GLog.Info(xmlValue)

	//CheckXml 系统配置审核配置
	err, wrongNoteInSys := CheckXml(bill, xmlValue)
	if err != nil {
		return obj, err
	}
	obj.Bill.WrongNote += wrongNoteInSys

	//4.项目自定义审核配置 导出校验 CheckXml
	err, wrongNoteInSelf := project.BillExportCheckXmlAdapter(bill.ProCode, formatObj, xmlValue)
	if err != nil {
		return obj, err
	}

	//wrongNote = saveXml的导出校验()+系统配置导出校验+自定义导出校验
	obj.Bill.WrongNote += wrongNoteInSelf

	//5.提取目录外数据 FetchNewHospitalAndCatalogue
	err = project.BillExportFetchNewHospitalAndCatalogueAdapter(bill.ProCode, obj)
	if err != nil {
		return obj, err
	}

	//6.机构抽取 FetchAgency
	err = project.BillExportFetchAgencyAdapter(bill.ProCode, obj, xmlValue)
	if err != nil {
		return obj, err
	}

	//7.错误统计
	global.GLog.Info(fmt.Sprintf("单号：%s, 发票数量：%v", bill.BillNum, len(obj.Invoice)-1))
	invoiceLen := 0
	if flag == "2" {
		invoiceLen = len(obj.Invoice) - 1
	}
	err, row, _ := WrongSum(model.ProCodeAndId{
		ProCode: bill.ProCode,
		ID:      bill.ID,
	}, invoiceLen, bill.ID, flag, fields)

	global.GLog.Info(fmt.Sprintf("一共错误统计数据有%v条", row))
	if err != nil {
		fmt.Println("------------------err-------------------------------: ", err)
		debug.PrintStack()
		// tx.Rollback()
		global.GLog.Error(fmt.Sprintf("错误统计id:::%v,错误:::%v", bill.ID, err))
		return obj, err
	}

	err, row = OcrSum(blocks, obj)
	if err != nil {
		return obj, err
	}

	//8.保存单据信息 SaveBill
	err = service.SaveBill(obj)
	if err != nil {
		return obj, err
	}

	// 扣费明细统计报表
	err = project.BillExportDeductionDetails(bill.ProCode, obj, xmlValue)
	if err != nil {
		fmt.Println("--------------扣费明细统计报表-----------------", err)
	}

	// 销毁报告
	managerService.SaveDelBills(bill, "下载文件")
	managerService.SaveDelBills(bill, "结果数据")

	// if err = tx.Commit().Error; err != nil {
	// 	tx.Rollback()
	// 	return obj, err
	// }

	return obj, err
}

// CheckXml 审核配置系统功能
func CheckXml(bill model.ProjectBill, xmlValue string) (err error, wrongNote string) {
	err, inspects := service.GetInspectConf(bill.ProCode)
	if err != nil {
		return err, wrongNote
	}
	for _, item := range inspects {
		if item.Msg != "" {
			wrongNote += item.Msg
		} else if item.XmlNodeCode != "" {
			//解释正则表达式
			reg := regexp.MustCompile(fmt.Sprintf("(<%v>.*</%v>|<%v/>)", item.XmlNodeCode, item.XmlNodeCode, item.XmlNodeCode))
			if reg == nil {
				fmt.Println("MustCompile err")
				return errors.New("编译xml错误"), wrongNote
			}
			//提取关键信息
			result := reg.FindAllString(xmlValue, -1)
			var validationMap = make(map[int]string, 0)
			for _, s := range item.Validation {
				validationMap[int(s)] = strconv.Itoa(int(s))
			}

			//不可缺节点
			if _, ok := validationMap[7]; len(result) == 0 && ok {
				global.GLog.Error("没有找到该节点信息:::" + item.XmlNodeCode)
				wrongNote += item.XmlNodeCode + "节点不存在;"
				continue
			}
			for i, text := range result {
				//fmt.Println("text = ", text)
				var val = ""
				if text == "<"+item.XmlNodeCode+"/>" {
					val = ""
				} else {
					val = strings.Replace(text, "<"+item.XmlNodeCode+">", "", -1)
					val = strings.Replace(val, "</"+item.XmlNodeCode+">", "", -1)
				}

				//if val == "" {
				//	continue
				//}
				//空值判断
				if _, ok := validationMap[6]; val == "" && ok {
					//global.GLog.Error("节点不能为空:::" + item.XmlNodeCode)
					wrongNote += item.XmlNodeCode + "【" + strconv.Itoa(i) + "】不能为空;"
					continue
				}
				if val == "" {
					continue
				}

				//只能录入
				if item.OnlyInput != "" {
					input := strings.Split(item.OnlyInput, ";")
					inMap := make(map[string]string, 0)
					for _, v := range input {
						inMap[v] = v
					}
					if _, ok := inMap[val]; !ok {
						wrongNote += item.XmlNodeCode + "录入内容错误;"
						continue
					}
				}

				//不能录入
				if item.NotInput != "" {
					input := strings.Split(item.NotInput, ";")
					inMap := make(map[string]string, 0)
					for _, v := range input {
						inMap[v] = v
					}
					if _, ok := inMap[val]; !ok {
						wrongNote += item.XmlNodeCode + "录入内容错误;"
						continue
					}
				}
				maxLen, _ := strconv.Atoi(item.MaxLen)
				minLen, _ := strconv.Atoi(item.MinLen)
				maxVal, _ := strconv.Atoi(item.MaxVal)
				minVal, _ := strconv.Atoi(item.MinVal)
				if maxLen > 0 && len(val) > maxLen {
					wrongNote += "长度超过限制;"
					continue
				}
				if minLen > 0 && len(val) < minLen {
					wrongNote += "长度小于限制;"
					continue
				}
				if maxVal > 0 {
					f, err := strconv.ParseFloat(val, 64)
					if err == nil && f > float64(maxVal) {
						wrongNote += "超过最大值;"
						continue
					}
				}
				if minVal > 0 {
					f, err := strconv.ParseFloat(val, 64)
					if err == nil && f < float64(minVal) {
						wrongNote += "小于最小值;"
						continue
					}
				}
				wrongNote += utils2.Validator(validationMap, val, item.XmlNodeName)
			}
		}

	}
	return err, wrongNote
}

// ChangeData 数据转换处理，生成结果数据
func ChangeData(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldConfList []model2.SysProField, fields []model3.ProjectField) (model4.ResultDataBill, error) {
	bill.WrongNote = ""
	timeStart := time.Now()
	//1.字段根据配置处理转换，替换，插入等 ProcessFields
	err, fieldMap := ProcessFields(fields, fieldConfList)
	if err != nil {
		return model4.ResultDataBill{}, err
	}

	//2.根据需求字段的转换 ResultData
	err, obj := project.BillExportResultDataAdapter(bill, blocks, fieldMap)
	if err != nil {
		return model4.ResultDataBill{}, err
	}
	global.GLog.Warn("ChangeData费时:::", zap.Any("毫秒", time.Since(timeStart).Milliseconds()))
	return obj, nil
}

// BillToXml 根据结果数据和导出配置生成xml
func BillToXml(proCode string, obj *model4.ResultDataBill) (err error, xmlValue string, formatObj interface{}) {
	timeStart := time.Now()
	err, export := service.GetExportConf(proCode)
	if err != nil {
		return err, xmlValue, formatObj
	}

	//格式化渲染xml的数据
	err, formatObj = project.FormatRenderObjTempFilterAdapter(proCode, *obj)
	if err != nil {
		return err, xmlValue, formatObj
	}

	err, tpl := project.BillExportTempFilterAdapter(proCode, export)
	if err != nil {
		return err, xmlValue, formatObj
	}
	xmlValue, err = tpl.Execute(pongo2.Context{"items0": formatObj})
	if err != nil {
		return err, xmlValue, formatObj
	}
	//global.GLog.Info(out)

	err, xmlValue = project.BillExportXmlDealAdapter(obj.Bill.ProCode, formatObj, xmlValue)
	if err != nil {
		return err, "", formatObj
	}

	err = project.SaveXml(&obj.Bill, formatObj, xmlValue, export.XmlType, false)
	if err != nil {
		return err, "", formatObj
	}

	global.GLog.Warn("BillToXml费时:::", zap.Any("毫秒", time.Since(timeStart).Milliseconds()))
	return nil, xmlValue, formatObj
}

// ProcessFields 字段根据配置处理转换，替换，插入等
func ProcessFields(fields []model3.ProjectField, fieldConfList []model2.SysProField) (err error, fieldMap map[string][]model3.ProjectField) {
	fieldMap = make(map[string][]model3.ProjectField, 0)
	var fieldsConf = make(map[string]model2.SysProField)
	for _, field := range fieldConfList {
		fieldsConf[field.Code] = field
	}

	for _, field := range fields {
		value, ok := fieldsConf[field.Code]
		if !ok {
			continue
		}
		field.FinalValue = field.ResultValue
		field.Issues = nil
		valQuestion(&field, value.QuestionChange)
		valInsert(&field, value.ValInsert)
		valChange(&field, value.ValChange)
		valIssue(&field, value.SysIssues)
		fieldMap[field.BlockID] = append(fieldMap[field.BlockID], field)
	}
	return nil, fieldMap

}

// valQuestion 问题件转换  A=值=编码;
func valQuestion(field *model3.ProjectField, questionChange string) {
	field.Issues = nil
	if questionChange == "" {
		return
	}
	arr := strings.Split(questionChange, ";")
	for _, str := range arr {
		if str == "" {
			continue
		}
		iArr := strings.Split(str, "=")
		if len(iArr) != 3 {
			continue
		}
		if iArr[0] == field.ResultValue {
			issue := model3.Issue{
				Type:    iArr[2],
				Code:    iArr[2],
				Message: const_data.IssueMap[iArr[2]],
			}
			field.Issues = append(field.Issues, issue)
			field.FinalValue = iArr[1]
		}
	}
}

// valChange 数值转换  A=50;1=4
func valChange(field *model3.ProjectField, valChange string) {
	if valChange == "" {
		return
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
		if str[:end] != field.FinalValue {
			continue
		}
		field.FinalValue = str[end+1:]
		break
	}
}

// valInsert 数值插入 4(插入位置)=5;1=A
func valInsert(field *model3.ProjectField, valInsert string) {
	if valInsert == "" || field.ResultValue == "" {
		return
	}
	arr := strings.Split(valInsert, ";")

	offset := 0
	for _, str := range arr {
		if str == "" {
			continue
		}
		end := strings.Index(str, "=")
		if end == -1 {
			continue
		}
		loc, err := strconv.Atoi(str[:end])
		loc += offset
		if err != nil || loc > len(field.FinalValue) {
			continue
		}
		field.FinalValue = field.FinalValue[:loc] + str[end+1:] + field.FinalValue[loc:]
		offset += len(str[end+1:])
	}
}

// valIssue 问题件配置
func valIssue(field *model3.ProjectField, sysIssues []model2.SysIssue) {
	//field.Issues = nil
	for _, sysIssue := range sysIssues {
		if sysIssue.InputVal == field.ResultValue {
			field.FinalValue = sysIssue.ChangeVal
			issue := model3.Issue{
				Type:    sysIssue.IssueType,
				Code:    sysIssue.Code,
				Message: sysIssue.Desc,
			}
			field.Issues = append(field.Issues, issue)
		}
	}
}

func CalculateFieldNum(fields []model3.ProjectField, op string) int {
	fieldCharacter := 0
	for _, field := range fields {
		mfield := reflect.ValueOf(&field).Elem()
		input := mfield.FieldByName(op + "Input").Interface().(string)
		value := mfield.FieldByName(op + "Value").Interface().(string)
		fmt.Println("---------------mfield-----------------:", input, value)
		if input == "no" || input == "no_if" {
			continue
		}
		for _, rr := range value {
			if rr == '?' || rr == '？' {
				continue
			}
			if unicode.Is(unicode.Han, rr) {
				fieldCharacter = fieldCharacter + 2
			} else {
				fieldCharacter = fieldCharacter + 1
			}
		}
	}
	return fieldCharacter
}
