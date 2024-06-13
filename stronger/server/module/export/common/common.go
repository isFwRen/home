package common

import (
	"errors"
	"fmt"
	"regexp"
	"server/global"
	"server/module/export/service"
	"server/module/export/utils"
	"server/module/load/model"
	proModel "server/module/pro_conf/model"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
)

func CheckInspectWrongNote(ProCode, xmlValue string) (error, string) {
	wrongNote := ""
	// fields := obj.Fields
	err, inspects := service.GetInspectConf(ProCode)
	if err != nil {
		return err, wrongNote
	}
	for _, item := range inspects {
		if item.XmlNodeCode != "" {
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
				wrongNote += item.XmlNodeName + "节点不存在;"
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
					wrongNote += item.XmlNodeName + "【" + strconv.Itoa(i+1) + "】不能为空;"
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
						wrongNote += item.XmlNodeName + "录入内容错误;"
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
						wrongNote += item.XmlNodeName + "录入内容错误;"
						continue
					}
				}
				maxLen, _ := strconv.Atoi(item.MaxLen)
				minLen, _ := strconv.Atoi(item.MinLen)
				maxVal, _ := strconv.Atoi(item.MaxVal)
				minVal, _ := strconv.Atoi(item.MinVal)
				if maxLen > 0 && len(val) > maxLen {
					wrongNote += item.XmlNodeName + "长度超过限制;"
					continue
				}
				if minLen > 0 && len(val) < minLen {
					wrongNote += item.XmlNodeName + "长度小于限制;"
					continue
				}
				if maxVal > 0 {
					f, err := strconv.ParseFloat(val, 64)
					if err == nil && f > float64(maxVal) {
						wrongNote += item.XmlNodeName + "超过最大值;"
						continue
					}
				}
				if minVal > 0 {
					f, err := strconv.ParseFloat(val, 64)
					if err == nil && f < float64(minVal) {
						wrongNote += item.XmlNodeName + "小于最小值;"
						continue
					}
				}
				wrongNote += utils.Validator(validationMap, val, item.XmlNodeName)
			}
		}

	}
	return nil, wrongNote

}

func CheckWrongNote(pro, xmlValue string, fields []model.ProjectField) (error, string) {
	// fmt.Println("---------CheckWrongNoteCheckWrongNote------------")
	wrongNote := ""
	err, wrongNote := CheckInspectWrongNote(pro, xmlValue)
	fmt.Println("--------wrongNote-err------------", err)
	if err != nil {
		return err, wrongNote
	}

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
	for _, field := range fields {
		// fmt.Println("---------field------------", field.Code, field.ResultValue)
		items, isExit := fieldCheckConfMap[field.Code]
		// fmt.Println("---------items------------", items)
		if isExit {
			for _, item := range items {
				fffs := strings.Split(item.Value, ";")
				// fmt.Println("---------fffs------------", fffs)
				for _, fff := range fffs {
					if item.CheckType == "1" {
						if field.ResultValue == fff {
							wrongNote += item.Mark + ";"
						}
					} else if item.CheckType == "2" {
						if strings.Index(field.ResultValue, fff) != -1 {
							wrongNote += item.Mark + ";"
						}
					} else if item.CheckType == "3" {
						if strings.Index(field.ResultValue, fff) == -1 {
							wrongNote += item.Mark + ";"
						}
					}
				}

			}
		}
	}
	return err, wrongNote

}

func SumFloat(a1, a2 float64, ff string) float64 {
	b1 := decimal.NewFromFloat(a1)
	b2 := decimal.NewFromFloat(a2)
	// b1.StringFixed(2)
	if ff == "+" {
		return ParseFloat(b1.Add(b2).StringFixed(2))
	}
	if ff == "-" {
		return ParseFloat(b1.Sub(b2).StringFixed(2))
	}

	return 0.00
}

func GetOneField(fields []model.ProjectField, code string, finalOrResult bool) (bool, string) {
	for _, field := range fields {
		if field.Code == code {
			if finalOrResult {
				return true, field.FinalValue
			} else {
				return true, field.ResultValue
			}
		}
	}
	return false, ""
}

func ToString(data float64) string {
	// return strconv.FormatFloat(data, 'E', -1, 64)
	return fmt.Sprintf("%v", data)
}

func RegReplace(data string, query string, value string) string {
	reg := regexp.MustCompile(query)
	return reg.ReplaceAllString(data, value)
}

func HasKey(data map[string]string, key string) bool {
	_, isOK := data[key]
	return isOK
}

func ParseFloat(value string) float64 {
	val, _ := strconv.ParseFloat(value, 64)
	return val
}

func RegIsMatch(value string, query string) bool {
	// reg := regexp.MustCompile(query)
	matched, _ := regexp.MatchString(query, value)
	return matched
}

func RegMatchAll(value string, query string) []string {
	reg := regexp.MustCompile(query)
	return reg.FindAllString(value, -1)
}

func SetNodeValue(xmlValue, nodeName, value string) string {
	reg := regexp.MustCompile(`>.*</` + nodeName + `>`)
	return reg.ReplaceAllString(xmlValue, ">"+value+"</"+nodeName+">")
}

func GetNodeValue(xmlValue, nodeName string) string {
	beginNode := strings.Index(xmlValue, "<"+nodeName+">") + len(nodeName) + 2
	endNode := strings.Index(xmlValue, "</"+nodeName+">")
	sValue := ""
	if beginNode != -1 && endNode != -1 {
		sValue = Substr(xmlValue, beginNode, endNode)
	}
	return sValue
}

func Substr(str string, start, end int) string {
	if start == -1 || end == -1 {
		return ""
	}
	return string(str[start:end])
}
