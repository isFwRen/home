/**
 * @Author: xingqiyi
 * @Description: 自定义单据导出中转站
 * @Date: 2021/11/24 10:43 上午
 */

package project

import (
	"errors"
	"server/global"
	model2 "server/module/export/model"
	"server/module/export/project/B0102"
	"server/module/export/project/B0103"
	"server/module/export/project/B0106"
	"server/module/export/project/B0108"
	"server/module/export/project/B0110"
	"server/module/export/project/B0113"
	"server/module/export/project/B0114"
	"server/module/export/project/B0116"
	"server/module/export/project/B0118"
	"server/module/export/project/B0121"
	"server/module/export/project/B0122"
	model3 "server/module/load/model"
	model4 "server/module/pro_conf/model"
	"server/module/pro_manager/model"

	"github.com/flosch/pongo2/v4"
)

// BillExportResultDataAdapter 结果数据
func BillExportResultDataAdapter(bill model.ProjectBill, blocks []model3.ProjectBlock, fieldMap map[string][]model3.ProjectField) (err error, obj model2.ResultDataBill) {
	if bill.ProCode == "" || global.ProCodeId[bill.ProCode] == "" {
		return errors.New("该单据项目编码错误,id::" + bill.ID), obj
	}
	switch bill.ProCode {
	case "B0118":
		return B0118.ResultData(bill, blocks, fieldMap)
	case "B0108":
		return B0108.ResultData(bill, blocks, fieldMap)
	case "B0114":
		return B0114.ResultData(bill, blocks, fieldMap)
	case "B0113":
		return B0113.ResultData(bill, blocks, fieldMap)
	case "B0121":
		return B0121.ResultData(bill, blocks, fieldMap)
	case "B0106":
		return B0106.ResultData(bill, blocks, fieldMap)
	case "B0103":
		return B0103.ResultData(bill, blocks, fieldMap)
	case "B0110":
		return B0110.ResultData(bill, blocks, fieldMap)
	case "B0122":
		return B0122.ResultData(bill, blocks, fieldMap)
	case "B0116":
		return B0116.ResultData(bill, blocks, fieldMap)
	case "B0102":
		return B0102.ResultData(bill, blocks, fieldMap)
	default:
		return errors.New("该项目没有自定义结果数据"), obj
	}
}

// BillExportTempFilterAdapter 模板过滤
func BillExportTempFilterAdapter(proCode string, e model4.SysExport) (err error, tpl *pongo2.Template) {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,id::" + proCode), tpl
	}
	switch proCode {
	case "B0118":
		return B0118.TempRender(e)
	case "B0108":
		return B0108.TempRender(e)
	case "B0114":
		return B0114.TempRender(e)
	case "B0113":
		return B0113.TempRender(e)
	case "B0121":
		return B0121.TempRender(e)
	case "B0106":
		return B0106.TempRender(e)
	case "B0103":
		return B0103.TempRender(e)
	case "B0110":
		return B0110.TempRender(e)
	case "B0122":
		return B0122.TempRender(e)
	case "B0116":
		return B0116.TempRender(e)
	case "B0102":
		return B0102.TempRender(e)
	default:
		return errors.New("该项目没有自定义导出模板过滤"), tpl
	}
}

// BillExportCheckXmlAdapter 导出校验
func BillExportCheckXmlAdapter(proCode string, obj interface{}, xmlValue string) (err error, wrongNote string) {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,id::" + proCode), wrongNote
	}
	switch proCode {
	case "B0118":
		return B0118.CheckXml(obj, xmlValue)
	case "B0108":
		return B0108.CheckXml(obj, xmlValue)
	case "B0114":
		return B0114.CheckXml(obj, xmlValue)
	case "B0113":
		return B0113.CheckXml(obj, xmlValue)
	case "B0121":
		return B0121.CheckXml(obj, xmlValue)
	case "B0106":
		return B0106.CheckXml(obj, xmlValue)
	case "B0103":
		return B0103.CheckXml(obj, xmlValue)
	case "B0110":
		return B0110.CheckXml(obj, xmlValue)
	case "B0122":
		return B0122.CheckXml(obj, xmlValue)
	case "B0116":
		return B0116.CheckXml(obj, xmlValue)
	case "B0102":
		return B0102.CheckXml(obj, xmlValue)
	default:
		return errors.New("该项目没有自定义导出导出校验"), wrongNote
	}
}

// BillExportXmlDealAdapter 自定义xml处理
func BillExportXmlDealAdapter(proCode string, obj interface{}, xmlValue string) (err error, newXmlValue string) {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,id::" + proCode), ""
	}
	switch proCode {
	case "B0118":
		return B0118.XmlDeal(obj, xmlValue)
	case "B0108":
		return B0108.XmlDeal(obj, xmlValue)
	case "B0114":
		return B0114.XmlDeal(obj, xmlValue)
	case "B0113":
		return B0113.XmlDeal(obj, xmlValue)
	case "B0121":
		return B0121.XmlDeal(obj, xmlValue)
	case "B0106":
		return B0106.XmlDeal(obj, xmlValue)
	case "B0103":
		return B0103.XmlDeal(obj, xmlValue)
	case "B0110":
		return B0110.XmlDeal(obj, xmlValue)
	case "B0122":
		return B0122.XmlDeal(obj, xmlValue)
	case "B0116":
		return B0116.XmlDeal(obj, xmlValue)
	case "B0102":
		return B0102.XmlDeal(obj, xmlValue)
	default:
		return errors.New("该项目没有自定义导出自定义xml处理"), ""
	}
}

// BillGetTotalFeeAdapter 自定义获取项目账单总额
func BillGetTotalFeeAdapter(proCode string, formatQuality map[int][]model4.SysQuality, xmlValue string) (err error, formatQuality2 map[int][]model4.SysQuality, obj model.BillAmount) {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,id::" + proCode), formatQuality, obj
	}
	switch proCode {
	case "B0118":
		return B0118.TotalFee(formatQuality, xmlValue)
	case "B0114":
		return B0114.TotalFee(formatQuality, xmlValue)
	default:
		global.GLog.Info("该项目没有自定义获取项目账单总额处理")
		return nil, formatQuality, obj
	}
}

// GetQingDan 自定义获取单据清单
func GetQingDan(form model.QingDanForm, obj model2.ResultDataBill) (err error, dan []model2.QingDan) {
	if form.ProCode == "" || global.ProCodeId[form.ProCode] == "" {
		return errors.New("该单据项目编码错误,id::" + form.ProCode), dan
	}
	switch form.ProCode {
	case "B0118":
		return B0118.GetQingDan(form, obj)
	default:
		return errors.New("该项目没有自定义获取单据清单处理"), dan
	}
}

// FormatRenderObjTempFilterAdapter 格式化渲染xml的对象
func FormatRenderObjTempFilterAdapter(proCode string, obj model2.ResultDataBill) (err error, formatObj interface{}) {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,id::" + proCode), formatObj
	}
	switch proCode {
	case "B0118":
		return B0118.FormatRenderObj(obj)
	case "B0108":
		return B0108.FormatRenderObj(obj)
	case "B0114":
		return B0114.FormatRenderObj(obj)
	case "B0113":
		return B0113.FormatRenderObj(obj)
	case "B0121":
		return B0121.FormatRenderObj(obj)
	case "B0106":
		return B0106.FormatRenderObj(obj)
	case "B0103":
		return B0103.FormatRenderObj(obj)
	case "B0110":
		return B0110.FormatRenderObj(obj)
	case "B0122":
		return B0122.FormatRenderObj(obj)
	case "B0116":
		return B0116.FormatRenderObj(obj)
	case "B0102":
		return B0102.FormatRenderObj(obj)
	default:
		return errors.New("该项目没有自定义格式化渲染xml的对象"), formatObj
	}
}

// BillExportJsonDealAdapter 自定义json处理
func BillExportJsonDealAdapter(proCode string, obj interface{}, jsonValue []byte) (err error, newJsonValue []byte) {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,code::" + proCode), nil
	}
	switch proCode {
	case "B0106":
		return B0106.JsonDeal(obj, jsonValue)
	case "B0103":
		return B0103.JsonDeal(obj, jsonValue)
	case "B0110":
		return B0110.JsonDeal(obj, jsonValue)
	case "B0122":
		return B0122.JsonDeal(obj, jsonValue)
	default:
		return nil, nil
	}
}

// BillExportDealErrJsonAdapter 自定义导出异常json处理
func BillExportDealErrJsonAdapter(proCode string, obj interface{}, jsonValue []byte) (newJsonValue []byte, err error) {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return nil, errors.New("该单据项目编码错误,code::" + proCode)
	}
	switch proCode {
	case "B0106":
		return B0106.DealErrJson(obj, jsonValue)
	case "B0103":
		return B0103.DealErrJson(obj, jsonValue)
	case "B0110":
		return B0110.DealErrJson(obj, jsonValue)
	default:
		return nil, nil
	}
}

// BillExportFetchNewHospitalAndCatalogueAdapter 提取目录外数据
func BillExportFetchNewHospitalAndCatalogueAdapter(proCode string, obj model2.ResultDataBill) error {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,code::" + proCode)
	}
	switch proCode {
	case "B0106":
		return B0106.FetchNewHospitalAndCatalogue(obj)
	case "B0103":
		return B0103.FetchNewHospitalAndCatalogue(obj)
	case "B0110":
		return B0110.FetchNewHospitalAndCatalogue(obj)
	default:
		return nil
	}
}

// BillExportFetchAgencyAdapter 机构抽取
func BillExportFetchAgencyAdapter(proCode string, obj model2.ResultDataBill, xmlValue string) error {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,code::" + proCode)
	}
	switch proCode {
	case "B0106":
		return B0106.FetchAgency(obj, xmlValue)
	case "B0103":
		return B0103.FetchAgency(obj, xmlValue)
	case "B0110":
		return B0110.FetchAgency(obj, xmlValue)
	default:
		return nil
	}
}

func BillExportDeductionDetails(proCode string, obj model2.ResultDataBill, xmlValue string) error {
	if proCode == "" || global.ProCodeId[proCode] == "" {
		return errors.New("该单据项目编码错误,code::" + proCode)
	}
	switch proCode {
	case "B0118":
		return B0118.DeductionDetails(obj, xmlValue)
	// case "B0103":
	// 	return B0103.FetchAgency(obj, xmlValue)
	// case "B0110":
	// 	return B0110.FetchAgency(obj, xmlValue)
	default:
		return nil
	}
}
