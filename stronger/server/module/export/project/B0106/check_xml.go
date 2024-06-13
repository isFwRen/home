package B0106

import (
	"encoding/json"
	"fmt"
	"reflect"
	"server/global"
	"server/module/export/service"
	"server/module/load/model"
	proModel "server/module/pro_conf/model"
	"server/utils"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/wxnacy/wgo/arrays"
	_ "github.com/wxnacy/wgo/arrays"
	_ "go.uber.org/zap"
)

func CheckXml(o interface{}, xmlValue string) (err error, wrongNote string) {
	global.GLog.Info("B0110:::CheckXml")
	obj := o.(FormatObj)

	global.GLog.Info("B0106:::CheckXml")
	global.GLog.Error("wrongNote：：：" + wrongNote)

	err, fwrongNote := CheckWrongNote(obj.Bill.ProCode, xmlValue, obj)
	wrongNote += fwrongNote

	fmt.Println("------------------------wrongNote-----------------------:", wrongNote)
	// ----------------------------------code---------------------------------------------
	constMap := InitConst(obj.Bill.ProCode)
	wrongNote += CodeCheck(obj, xmlValue, constMap)
	wrongNote += XmlCheck(obj, xmlValue, constMap)

	// ----------------------------------xml---------------------------------------------

	global.GLog.Error("wrongNote：：：" + wrongNote)
	fmt.Println("------------------------wrongNote-----------------------:", wrongNote)

	return err, wrongNote
}

func CheckWrongNote(pro, xmlValue string, obj FormatObj) (error, string) {
	// fmt.Println("---------CheckWrongNoteCheckWrongNote------------")
	wrongNote := ""

	err, fieldCheckConfs := service.GetProFieldCheckConf(pro)
	fmt.Println("--------fieldCheckConfs-err------------", err)
	if err != nil {
		return err, wrongNote
	}
	fieldCheckConfMap := make(map[string][]proModel.SysProFieldCheck)
	// fmt.Println("---------fieldCheckConfMap------------", len(fieldCheckConfMap))
	for _, fieldCheckConf := range fieldCheckConfs {
		fieldCheckConfMap[fieldCheckConf.Code] = fieldCheckConf.SysProFieldChecks
	}
	eleLen := reflect.ValueOf(obj).NumField()
	for j := 0; j < eleLen; j++ {
		if reflect.TypeOf(obj).Field(j).Name != "Bill" && reflect.TypeOf(obj).Field(j).Name != "Fields" {
			//每张发票每种类型的字段
			fmt.Println("---------------------------", reflect.TypeOf(obj).Field(j).Name)
			fieldsMaps := reflect.ValueOf(obj).Field(j).Interface().([]FieldsMap)
			for _, fieldsMap := range fieldsMaps {
				if fieldsMap.Code == "" {
					continue
				}
				for _, field := range fieldsMap.Fields {
					items, isExit := fieldCheckConfMap[field.Code]
					// fmt.Println("---------items------------", items)
					if isExit {
						for _, item := range items {
							fffs := strings.Split(item.Value, ";")
							// fmt.Println("---------fffs------------", fffs)
							for _, fff := range fffs {
								mess := "账单号:" + fieldsMap.Code + item.Mark + ";"
								if strings.Index(wrongNote, mess) != -1 {
									continue
								}
								if item.CheckType == "1" {
									if field.ResultValue == fff {
										wrongNote += mess
									}
								} else if item.CheckType == "2" {
									if strings.Index(field.ResultValue, fff) != -1 {
										wrongNote += mess
									}
								} else if item.CheckType == "3" {
									if strings.Index(field.ResultValue, fff) == -1 {
										wrongNote += mess
									}
								}
							}

						}
					}
				}
			}
		}
	}
	// for
	return err, wrongNote

}

// CodeCheck 20230516新增
func CodeCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	wrongNote := ""
	fields := obj.Fields
	//sanFang := constMap["sanfang"]
	//muLuBianMa := constMap["mulubianma"]
	// YLMLM622002 := constMap["YLML_M622002"]

	items := map[string]string{
		"fc012": "fc013",
		"fc025": "fc032",
		"fc026": "fc033",
		"fc027": "fc034",
		"fc028": "fc035",
		"fc029": "fc036",
		"fc030": "fc037",
		"fc031": "fc038",
	}
	var fc091Arr []string
	var fc005Arr []string
	for _, invoice := range obj.Invoice {
		var filedArr []string
		_, fc005 := GetOneField(invoice.Fields, "fc005", true)
		var notTranscoding []string
		for _, QingDan := range invoice.QingDan {
			for _, field := range QingDan.Fields {
				code2, is := items[field.Code]
				fmt.Println(is)

				//需求编码 ：CSB0106RC0062000 fc012,fc025,fc026,fc027,fc028,fc029,fc030,fc031
				//以上字段的录入值不为常量表的数据（没有正常转码）时，出导出校验：“发票XXX的项目名称“xxx”不在目录中，请检查；
				//”（第一个XXX为同一发票的fc005的值，第二个xxx为字段的录入值）（以上字段为空或为A或包含?时不执行该校验）
				if RegIsMatch(field.Code, `^(fc012|fc025|fc026|fc027|fc028|fc029|fc030|fc031)$`) {
					// fmt.Println("***************************进不去*****************")
					if field.FinalValue != "" && field.FinalValue != "A" && field.FinalValue != "?" {
						// fmt.Println("***************************进去*********************")
						// if !HasKey(YLMLM622002, field.FinalValue) {
						if !RegIsMatch(field.FinalValue, `^\d+$`) {
							notTranscoding = append(notTranscoding, field.ResultValue)
							//wrongNote += "发票" + fc005 + "的项目名称" + field.ResultValue + "不在目录中，请检查；"
						}
					}
				}
				//需求编码 ：CSB0106RC0059000以下字段录入内容包含以下大项时，出导出校验：发票XXX的“xxx”为大项，请确认；（注：XXX为同一发票的fc005的值，xxx为字段的录入内容）
				//字段：fc012、fc025、fc026、fc027、fc028、fc029、fc030、fc031
				//大项：西药费、民族药费、检查费、化验费、检验费、放射检查费、特检费、治疗费、手术费、输血费、材料费、护理费、抢救费、自制制剂、辩证处方费、其他费、中药费、冲洗溶液、中成药、中草药
				if RegIsMatch(field.Code, `^(fc012|fc025|fc026|fc027|fc028|fc029|fc030|fc031)$`) {
					if RegIsMatch(field.ResultValue, `^(西药费|民族药费|检查费|化验费|检验费|放射检查费|特检费|治疗费|手术费|输血费|材料费|护理费|抢救费|自制制剂|辩证处方费|其他费|中药费|冲洗溶液|中成药|中草药)$`) {
						wrongNote += "发票" + fc005 + "的“" + field.ResultValue + "”为大项，请确认;"
					}
					//需求编码：CSB0106RC0061000 当fc012、fc025、fc026、fc027、fc028、fc029、fc030、fc031录入内容为“0.9%氯化钠注液射”时，出导出校验：0.9%氯化钠注液射录入错误，请修改。(多条只出一条校验)
					if field.ResultValue == "0.9%氯化钠注液射" && strings.Index(wrongNote, "0.9%氯化钠注液射录入错误，请修改;") == -1 {
						wrongNote += "0.9%氯化钠注液射录入错误，请修改;"
					}

				}

				//需求编码 ：CSB0106RC0060000 当以下字段包含英文状态下问号“?”时，出一条导出校验，导出校验提示为：“发票XXX项目名称（项目金额）不能包含问号，请修改；”（XXX为同一发票的fc005的值）(多条只出一条校验)
				if RegIsMatch(field.Code, `^(fc012|fc014|fc025|fc039|fc026|fc040|fc027|fc041|fc028|fc042|fc029|fc043|fc030|fc044|fc031|fc045)$`) && strings.Index(field.ResultValue, "?") != -1 {
					mes := "发票" + fc005 + "项目名称（项目金额）不能包含问号，请修改;"
					if strings.Index(wrongNote, mes) == -1 {
						wrongNote += mes
					}
				}

				//需求编码 ：CSB0106RC0063000 在同一发票中(根据fc018发票属性、fc019清单所属属性判断是否属于同一张发票)对下列分组中的录入内容进行校验(每一行为一组)；
				//当左边字段录入内容包含"床位"、"人间"、"病房"、"陪护床"、"走廊"且不包含"取暖"或"空调"字样时，校验对应右边字段录入内容是否为“1”
				//如是则出导出校验：XXX发票号床位费数量录入有误；（XXX为同一发票的fc005的值）
				// 	if RegIsMatch(field.ResultValue, `(床位|人间|病房|陪护床|走廊)`) && !RegIsMatch(field.Code, `(取暖|空调)`) {
				if RegIsMatch(field.Code, `^(fc012|fc025|fc026|fc027|fc028|fc029|fc030|fc031)$`) {
					if RegIsMatch(field.ResultValue, `(床位|人间|病房|陪护床|走廊)`) && !RegIsMatch(field.ResultValue, `(取暖|空调)`) {
						_, values := GetOneField(QingDan.Fields, code2, false)
						if values == "1" {
							wrongNote += fc005 + "发票号床位费数量录入有误;"
						}
					}
				}

				//需求编码 CSB0106RC0085000 清单
				if RegIsMatch(field.ResultValue, `[?？]`) {
					filedArr = append(filedArr, field.Name)
				}
			}

			//CSB0106RC0118000
			sFields := [][]string{
				//名称 金额 数量
				{"fc012", "fc014", "fc013"},
				{"fc025", "fc039", "fc032"},
				{"fc026", "fc040", "fc033"},
				{"fc027", "fc041", "fc034"},
				{"fc028", "fc042", "fc035"},
				{"fc029", "fc043", "fc036"},
				{"fc030", "fc044", "fc037"},
				{"fc031", "fc045", "fc038"},
			}
			var cSB0106RC0118000Msg []string
			for _, item := range sFields {
				_, name := GetOneField(QingDan.Fields, item[0], false)
				_, money := GetOneField(QingDan.Fields, item[1], false)
				_, number := GetOneField(QingDan.Fields, item[2], false)
				if !(len(name) == 0 && len(money) == 0 && len(number) == 0) &&
					!(len(name) > 0 && len(money) > 0 && len(number) > 0) {
					for _, value := range item {
						_, val := GetOneField(QingDan.Fields, value, true)
						if len(val) == 0 {
							cSB0106RC0118000Msg = append(cSB0106RC0118000Msg, value)
						}
					}
				}
			}
			if len(cSB0106RC0118000Msg) > 0 {
				mes := "【" + fc005 + "】发票【" + strings.Join(cSB0106RC0118000Msg, "、") + "】字段位为空；"
				if strings.Index(wrongNote, mes) == -1 {
					wrongNote += mes
				}
			}
		}

		_, fc005nc := GetOneField(invoice.Fields, "fc005", true)

		for i, field := range invoice.Fields {
			//需求编码 CSB0106RC0085000
			//1.校验所有字段的录入值，当包含?或？时，出导出校验：发票【xxx】的【yyy】存在?号，请核实；
			//2.如同一发票下多个字段包含?或？，则导出校验提示格式为：发票【xxx】的【yyy1、yyy2】存在?号，请核实；（xxx为发票号fc005的值，yyy为包含问号的字段名）
			//3.当最后一页字段包含?或？时，则直接出一条导出校验：最后一页存在?，请检查；
			if RegIsMatch(field.ResultValue, `[?？]`) {
				filedArr = append(filedArr, field.Name)
			}
			if i == len(invoice.Fields)-1 {
				if len(filedArr) > 0 {
					join := strings.Join(filedArr, "、")
					mes := "发票" + fc005nc + "的【" + join + "】存在?号，请核实;"
					if strings.Index(wrongNote, mes) == -1 {
						wrongNote += mes
						filedArr = make([]string, 0)
					}
				}
				if len(notTranscoding) > 0 {
					join := strings.Join(notTranscoding, "、")
					mes := "发票" + fc005 + "的项目名称“" + join + "”不在目录中，请检查；"
					fmt.Println("=======================notTranscoding===========", notTranscoding)
					if strings.Index(wrongNote, mes) == -1 {
						wrongNote += mes
						notTranscoding = make([]string, 0)
					}
				}
			}

			//CSB0106RC0108000
			//同一发票下，当fc090录入值为1时，校验fc092、fc093、fc094、fc095、fc096是否存在结果值，任意一个字段结果值为空时，出导出校验：电子发票五要素未录入齐全，请确认是否为电子发票；
			sFiled := []string{"fc092", "fc093", "fc094", "fc095", "fc096"}
			if field.Code == "fc090" {
				_, fc091Val := GetOneField(invoice.Fields, "fc091", false)
				if fc091Val != "3" && field.ResultValue == "1" {
					for k := 0; k < len(sFiled); k++ {
						_, val := GetOneField(invoice.Fields, sFiled[k], true)
						if val == "" {
							mes := "【" + fc005 + "】" + "电子发票五要素未录入齐全，请确认是否为电子发票；"
							if strings.Index(wrongNote, mes) == -1 {
								wrongNote += mes
							}
						}

					}
				}
			}

			//CSB0106RC0109000
			//"同一发票下，当fc091录入值为1时，（xxx为fc005的结果值）
			//1、校验fc092的结果值是否为12位数，否则出导出校验：【XX】发票增值税票据代码字段应为12位数，请检查；
			//2、校验fc093的结果值是否为8位数，否则出导出校验：【XX】发票增值税票据号码字段应为8位数，请检查；
			//3、校验fc095的结果值是否为6位数，否则出导出校验：【XX】发票校验码字段应为6位数，请检查；
			//4、fc097的结果值不为“01”时，出导出校验，【XX】发票查验结果不通过，请确认五要素是否录入正确；"
			//CSB0106RC0110000
			//"同一发票下，当fc091录入值为2时，（xxx为fc005的结果值）
			//1、校验fc092的结果值是否为8位数，否则出导出校验：XX】发票财政部票据代码字段应为8位数，请检查；
			//2、校验fc093的结果值是否为10位数，否则出导出校验：【XX】发票财政部票据号码字段应为10位数，请检查；"
			//3、fc097的结果值不为“01”时，出导出校验，【XX】发票查验结果不通过，请确认五要素是否录入正确；
			sFiledTo := [][]string{{"fc092", "12", "发票增值税票据代码字段应为12位数，请检查；"}, {"fc093", "8", "发票增值税票据号码字段应为8位数，请检查；"}, {"fc095", "6", "发票校验码字段应为6位数，请检查；"}}
			sFiledThree := [][]string{{"fc092", "8", "发票财政部票据代码字段应为8位数，请检查；"}, {"fc093", "10", "发票财政部票据号码字段应为10位数，请检查；"}}
			if field.Code == "fc091" {
				if field.ResultValue == "1" {
					for k := 0; k < len(sFiledTo); k++ {
						_, val := GetOneField(invoice.Fields, sFiledTo[k][0], true)
						atoi, _ := strconv.Atoi(sFiledTo[k][1])
						if val != "" && len(val) != atoi {
							mes := "【" + fc005 + "】" + sFiledTo[k][2]
							if strings.Index(wrongNote, mes) == -1 {
								wrongNote += mes
							}
						}
					}
					//_, val097 := GetOneField(invoice.Fields, "fc097", true)
					//if val097 != "" && val097 != "01" {
					//	mes := "【" + fc005 + "】发票查验结果不通过，请确认五要素是否录入正确；"
					//	if strings.Index(wrongNote, mes) == -1 {
					//		wrongNote += mes
					//	}
					//}

				} else if field.ResultValue == "2" {
					_, fc094Val := GetOneField(invoice.Fields, "fc094", true)
					fc094Parse, _ := time.Parse("2006-01-02", fc094Val)
					subTime := time.Now().Sub(fc094Parse)
					oneYear, _ := time.ParseDuration("8760h")
					if subTime < oneYear {
						for k := 0; k < len(sFiledThree); k++ {
							_, val := GetOneField(invoice.Fields, sFiledThree[k][0], true)
							atoi, _ := strconv.Atoi(sFiledThree[k][1])
							if val != "" && len(val) != atoi {
								mes := "【" + fc005 + "】" + sFiledThree[k][2]
								if strings.Index(wrongNote, mes) == -1 {
									wrongNote += mes
								}
							}
						}
						_, val := GetOneField(invoice.Fields, "fc097", true)
						if val != "01" {
							fc005Arr = append(fc005Arr, fc005)
						}
					}
				} else if field.ResultValue == "3" {
					_, val := GetOneField(invoice.Fields, "fc005", false)
					if len(val) < 20 {
						fc091Arr = append(fc091Arr, fc005)
					}
				}
			}
		}

		//需求编码 ：CSB0106RC0055000 同一发票属性中，当fc016录入值和fc006录入值不一致时，出导出校验：发票XXX【发票属性xxx】的治疗医院异常，请检查。
		//（XXX为同一发票的fc005的结果值，xxx为同一发票的fc018的结果值）例如：发票123456【发票属性1】的治疗医院异常，请检查。
		_, fc016 := GetOneField(invoice.Fields, "fc016", false)
		_, fc006 := GetOneField(invoice.Fields, "fc006", false)
		_, fc018 := GetOneField(invoice.Fields, "fc018", true)
		if fc016 != fc006 {
			wrongNote += "发票" + fc005 + "【发票属性" + fc018 + "】的治疗医院异常，请检查;"
		}

		//需求编码 ：CSB0106RC0051000 当fc017录入内容为1时，出导出校验：XXX发票号有两个日期，请检查是否为最早日期。（XXX为同一发票的fc005的结果值）
		_, fc005One := GetOneField(invoice.Fields, "fc005", true)
		_, fc017 := GetOneField(invoice.Fields, "fc017", false)
		if fc017 == "1" {
			wrongNote += fc005One + "发票号有两个日期，请检查是否为最早日期；"

		}

		//需求编码：CSB0106RC0056000 同一发票属性下，当fc010和c024录入值不一致（其中一个为空也算不一致，同时为空不校验）时，
		//出导出校验：XXX发票号，报销金额录入不一致。（XXX为同一发票的fc005的结果值）（fc024不存在时不执行该校验）
		_, fc010 := GetOneField(invoice.Fields, "fc010", false)
		isFc024, fc024 := GetOneField(invoice.Fields, "fc024", false)
		if isFc024 && fc010 != fc024 {
			wrongNote += fc005 + "发票号，报销金额录入不一致;"
		}
		_, fc010Value := GetOneField(invoice.Fields, "fc010", true)

		fmt.Println("/------------------", fc010Value)

		_, fc090 := GetOneField(invoice.Fields, "fc090", true)
		_, fc097 := GetOneField(invoice.Fields, "fc097", true)
		if fc090 == "Y" && fc097 == "" {
			wrongNote += fc005 + "发票,发票查验结果为空,请检查;"
		}

	}

	//需求编码：CSB0106RC0053000 fc007、fc008出问题件时,将自动回传改为手动回传，出导出校验：诊断信息异常,请确认；
	for _, field := range fields {
		if strings.Index(wrongNote, "诊断信息异常,请确认；") == -1 {
			if field.Code == "fc008" || field.Code == "fc007" {
				if len(field.Issues) > 0 {
					wrongNote += "诊断信息异常,请确认；"
				}
			}
		}

		//需求编码：CSB0106RC0052000 fc007录入值为空时，出导出校验：疾病诊断不能为空；
		if field.Code == "fc007" {
			fmt.Println("fc007录入", field.ResultValue)
			if field.ResultValue == "" {
				wrongNote += "疾病诊断不能为空;"
			}
		}
		//20230713 取消
		//if strings.Index(wrongNote, "案件中存在?号，请核实；") == -1 {
		//	if RegIsMatch(field.ResultValue, "\\?|？") {
		//		wrongNote += "案件中存在?号，请核实；"
		//	}
		//}

		//需求编码：CSB0106RC0054000 当fc087的录入值不为《B0106_陕西国寿理赔_第三方出具单位》中的“第三方出具单位名称”（第三列）中的内容时，出导出校验：第三方出具单位录入有误，请检查
		if field.Code == "fc087" {
			_, total := utils.FetchConst(obj.Bill.ProCode, "B0106_陕西国寿理赔_第三方出具单位", "第三方出具单位代码", map[string]string{"第三方出具单位名称": field.ResultValue})
			//fmt.Println("fc087录入", field.ResultValue, "=", sanFang[field.FinalValue])
			if total == 0 {
				if strings.Index(wrongNote, "第三方出具单位录入有误，请检查;") == -1 {
					wrongNote += "第三方出具单位录入有误，请检查;"
				}
			}
		}

	}

	if len(fc005Arr) > 0 {
		mes := "【" + strings.Join(fc005Arr, "、") + "】" + "发票查验结果不通过，请确认五要素是否录入正确；"
		if strings.Index(wrongNote, mes) == -1 {
			wrongNote += mes
		}
	}

	if len(fc091Arr) > 0 {
		mes := "【" + strings.Join(fc091Arr, "、") + "】" + "发票票据号不为20位数，请检查；"
		if strings.Index(wrongNote, mes) == -1 {
			wrongNote += mes
		}
	}

	return wrongNote
}

// XmlCheck 20230516新增
func XmlCheck(obj FormatObj, xmlValue string, constMap map[string]map[string]string) string {
	//20230516新增 每个案件可能包含多个<rcptNo>，判断所有的<rcptNo>，如果有重复值，将自动回传改为手动回传，并出导出校验提示：发票号[xxx]重复，请检查；（xxx是rcptNo值）
	wrongNote := ""
	var otherInfo OtherInfo

	err := json.Unmarshal([]byte(obj.Bill.OtherInfo), &otherInfo)
	if err != nil {
		global.GLog.Error("otherInfo", zap.Error(err))
	}

	//需求编码：CSB0106RC0050000 下载文件中的bpoSendRemark中有内容时，需要将bpoSendRemark中的内容，放到对应单号的导出校验中。
	//（若bpoSendRemark中内容有【】时，删除【】及其中的所有信息，如“aaa【bbb】ccc”，出导出校验：aaaccc，且屏蔽bpoSendRemark中的“；”，“；”不显示在导出校验中。）
	bpoSendRemark := otherInfo.BpoSendRemark
	if bpoSendRemark != "" {
		bpoSendRemark = RegReplace(bpoSendRemark, `【.*】`, "")
		bpoSendRemark = RegReplace(bpoSendRemark, `(；|，|；)`, "")
		if bpoSendRemark != "" {
			wrongNote += bpoSendRemark + ";"
		}
	}

	items := []string{}
	//需求编码 : CSB0106RC0049000
	//rcptNos := RegMatchAll(xmlValue, `<rcptNo>.*?<\/rcptNo>`)
	//for _, rcptNo := range rcptNos {
	//	if arrays.Contains(items, rcptNo) != -1 {
	//		wrongNote += "发票号[" + rcptNo + "]重复，请检查；"
	//	} else {
	//		items = append(items, rcptNo)
	//	}
	//}

	rcptInfoLists := RegMatchAll(xmlValue, `<rcptInfoList>[\s\S]*?<\/rcptInfoList>`)
	for _, rcptInfoList := range rcptInfoLists {
		//需求编码 ：CSB0106RC0058000 判断每个<rcptInfoList>发票大节点下的<socialPayAmnt>的值是否“大于或等于”<rcptAmnt>，若是，则将自动回传改为手动回传，并出导出校验，
		//提示为：发票号[xxx]报销金额大于或等于发票总金额，请检查；（xxx为同一个rcptInfoList大节点下的rcptNo的值）
		rcptNo := GetNodeValue(rcptInfoList, "rcptNo")
		indexRcptNo := strings.Index(rcptInfoList, "</rcptNo>")
		//CSB0106RC0049001
		if arrays.Contains(items, rcptNo) != -1 {
			wrongNote += "发票号[" + rcptNo + "]重复，请检查；"
		} else {
			items = append(items, rcptNo)
		}

		socialPayAmnt := GetNodeValue(rcptInfoList, "socialPayAmnt")
		rcptAmnt := GetNodeValue(rcptInfoList, "rcptAmnt")
		if ParseFloat(socialPayAmnt) >= ParseFloat(rcptAmnt) && indexRcptNo != -1 {
			wrongNote += "发票号[" + rcptNo + "]报销金额大于或等于发票总金额，请检查；"
		}
		//需求编码 ：CSB0106RC0057000 判断每个<rcptInfoList>发票大节点下的<endDate>的日期是否早于<beginDate>的日期，若是，则将自动回传改为手动回传，并出导出校验，
		//提示为：发票号[xxx]出院时间早于入院时间，请检查；（xxx为同一个rcptInfoList大节点下的rcptNo的值）
		endDate := GetNodeValue(rcptInfoList, "endDate")
		beginDate := GetNodeValue(rcptInfoList, "beginDate")
		a, _ := time.Parse("2006-01-02", endDate)
		b, _ := time.Parse("2006-01-02", beginDate)
		if a.Before(b) {
			wrongNote += "发票号[" + rcptNo + "]出院时间早于入院时间，请检查；"
		}

		sum := 0.0
		rcptLists := RegMatchAll(rcptInfoList, `<(rcptList|errorRcptList)>[\s\S]*?<\/(rcptList|errorRcptList)>`)

		if len(rcptLists) == 0 {
			continue
		}
		for _, rcptList := range rcptLists {

			quantity := GetNodeValue(rcptList, "quantity")
			price := GetNodeValue(rcptList, "price")
			if quantity != "" && price != "" {
				total := SumFloat(ParseFloat(quantity), ParseFloat(price), "*")
				sum = SumFloat(sum, total, "+")
			}
		}
		chae := SumFloat(ParseFloat(rcptAmnt), sum, "-")
		expenMode := GetNodeValue(rcptInfoList, "expenMode")
		if expenMode == "0020" && (chae >= 0.1 || chae <= -0.1) {
			fmt.Println("*********0020**********chae=", chae)
			wrongNote += "发票号[" + rcptNo + "]清单项目金额与发票金额不一致，差额:" + ToString(chae) + "，请检查；"
		}
		if expenMode == "0040" && (chae >= 2 || chae <= -2) {
			fmt.Println("*********0040**********chae=", chae)
			wrongNote += "发票号[" + rcptNo + "]清单项目金额与发票金额不一致，差额:" + ToString(chae) + "，请检查；"
		}

	}

	return wrongNote
}

func CheckFieldHasIssue(fields []model.ProjectField, code string) bool {
	for _, field := range fields {
		if field.Code == code {
			if len(field.Issues) > 0 {
				return true
			}
		}
	}
	return false
}

func InitConst(proCode string) map[string]map[string]string {
	constObj := make(map[string]map[string]string, 0)
	nameMap := [][]string{
		{"muLuBianMa", "B0106_陕西国寿理赔_数据库编码对应表", "1", "3"},
		{"sanfang", "B0106_陕西国寿理赔_第三方出具单位", "1", "2"},
		{"YLML_M622002", "B0106_陕西国寿理赔_医疗目录M6100002022002", "1", "1"},
	}
	for _, item := range nameMap {
		c, ok := global.GProConf[proCode].ConstTable[item[1]]
		tempMap := make(map[string]string, 0)
		if ok {
			for _, arr := range c {
				k, _ := strconv.Atoi(item[2])
				v, _ := strconv.Atoi(item[3])
				tempMap[strings.TrimSpace(arr[k])] = arr[v]
			}
		}
		constObj[item[0]] = tempMap
	}
	return constObj
}

func GetFieldsFinal(fields []model.ProjectField, code string) string {
	for _, field := range fields {
		if field.Code == code {
			return field.FinalValue
		}
	}
	return ""
}
