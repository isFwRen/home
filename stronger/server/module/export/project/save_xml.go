/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/29 5:18 下午
 */

package project

import (
	"errors"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"github.com/djimenez/iconv-go"
	"go.uber.org/zap"
	"io/ioutil"
	"regexp"
	"server/global"
	"server/module/pro_manager/model"
	"server/utils"
	"strings"
)

func SaveXml(bill *model.ProjectBill, obj interface{}, xml, xmlType string, isErr bool) error {
	matched := regexp.MustCompile("\n\\s*\n")
	xml = matched.ReplaceAllString(xml, "\n")
	xmlFileDir := global.GConfig.LocalUpload.FilePath + bill.ProCode + "/upload_xml/" +
		fmt.Sprintf("%v/%v/%v/",
			bill.CreatedAt.Year(), int(bill.CreatedAt.Month()),
			bill.CreatedAt.Day())
	exists, err := utils.PathExists(xmlFileDir)
	if err != nil {
		return err
	}
	if !exists {
		err = utils.CreateDir(xmlFileDir)
		if err != nil {
			return err
		}
	}
	if utils.RegIsMatch("^(B0113|B0121)$", bill.ProCode) {
		xmlType = "utf-8"
	}
	global.GLog.Info("saveXml", zap.Any("", xmlType))
	//global.GLog.Info("saveXml", zap.Any("", xml))
	if xmlType != "utf-8" {
		xml, err = iconv.ConvertString(xml, "utf-8", xmlType)
		if err != nil {
			line := strings.Split(xml, "\n")
			bill.WrongNote = fmt.Sprintf("报文第%v行存在特殊字符，请确认并修改；", len(line))
			global.GLog.Error("xml转码失败", zap.Error(err))
		}
	}
	//enc := mahonia.NewEncoder(xmlType)
	//encStr := enc.ConvertString(xml)
	err = ioutil.WriteFile(xmlFileDir+"/"+bill.BillNum+".xml", []byte(xml), 0666)
	if err != nil {
		global.GLog.Error("写入xml失败", zap.Error(err))
		return errors.New("写入xml失败")
	}

	if xmlType == "utf-8" && utils.RegIsMatch("^(B0103|B0110|B0106|B0122)$", bill.ProCode) {
		global.GLog.Info("", zap.Any("", "写入json"))
		xmlReader := strings.NewReader(xml)
		result, err := xj.Convert(xmlReader)
		//fmt.Println(result)
		if err != nil {
			global.GLog.Error("写入json失败", zap.Error(err))
			return err
		}
		data := []byte("")
		if isErr {
			data, err = BillExportDealErrJsonAdapter(bill.ProCode, obj, result.Bytes())
			if err != nil {
				return err
			}
		} else {
			err, data = BillExportJsonDealAdapter(bill.ProCode, obj, result.Bytes())
			if err != nil {
				return err
			}
		}
		err = ioutil.WriteFile(xmlFileDir+"/"+bill.BillNum+".json", data, 0666)
	} else {
		global.GLog.Error("导出模板非utf-8编码  不导出json  请重修修改模板类型重新保存")
	}
	return err
}
