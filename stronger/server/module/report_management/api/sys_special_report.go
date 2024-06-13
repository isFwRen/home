/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/5/16 14:18
 */

package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"server/global"
	"server/global/response"
	model3 "server/module/export/model"
	model2 "server/module/pro_manager/model"
	"server/module/report_management/model/request"
	service2 "server/module/report_management/service"
	"server/module/sys_base/model"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
	"strings"
	"time"
)

// PageNewHospitalAndCatalogue
// @Tags special-report(特殊报表)
// @Summary 特殊报表--目录外数据列表
// @Auth xingqiyi
// @Date 2021/11/3 3:54 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param startTime      query   string   true   "开始时间今天开始时间格式'2022-07-06T16:00:00.000Z'"
// @Param endTime        query   string   true   "结束时间今天结束时间格式'2022-07-06T16:00:00.000Z'"
// @Param type      query   int      false   "类型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /special-report/new-hospital-catalogue/page [get]
func PageNewHospitalAndCatalogue(c *gin.Context) {
	var newHospitalAndCatalogueSearch request.NewHospitalAndCatalogueSearch
	err := c.ShouldBindQuery(&newHospitalAndCatalogueSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, list := service2.PageNewHospitalAndCatalogue(newHospitalAndCatalogueSearch)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     total,
			PageIndex: newHospitalAndCatalogueSearch.PageIndex,
			PageSize:  newHospitalAndCatalogueSearch.PageSize,
		}, c)
	}
}

// ExportNewHospitalAndCatalogue
// @Tags special-report(特殊报表)
// @Summary 特殊报表--目录外数据导出
// @Auth xingqiyi
// @Date 2021/11/3 3:54 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param startTime      query   string   true   "开始时间今天开始时间格式'2022-07-06T16:00:00.000Z'"
// @Param endTime        query   string   true   "结束时间今天结束时间格式'2022-07-06T16:00:00.000Z'"
// @Param type      query   int      false   "类型"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /special-report/new-hospital-catalogue/export [get]
func ExportNewHospitalAndCatalogue(c *gin.Context) {
	var search request.NewHospitalAndCatalogueExportSearch
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list := service2.ExportNewHospitalAndCatalogue(search)
	path := global.GConfig.LocalUpload.FilePath + global.PathProReportExport + search.ProCode + "/"
	// 尝试创建此路径
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		global.GLog.Error("upload file fail:", zap.Any("err", err))
		return
	}
	name := fmt.Sprintf("理赔2.0目录外数据-%v.xlsx", time.Now().Format("20060102"))
	err = utils.ExportBigExcel(path, name, "sheet", list)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.FileAttachment(path+name, name)
}

// PageExtractAgency
// @Tags special-report(特殊报表)
// @Summary 特殊报表--机构抽取列表
// @Auth xingqiyi
// @Date 2021/11/3 3:54 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param startTime      query   string   true   "开始时间今天开始时间格式'2022-07-06T16:00:00.000Z'"
// @Param endTime        query   string   true   "结束时间今天结束时间格式'2022-07-06T16:00:00.000Z'"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /special-report/extract-agency/page [get]
func PageExtractAgency(c *gin.Context) {
	var search model.BaseTimePageCode
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, list := service2.PageExtractAgency(search)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      list,
			Total:     total,
			PageIndex: search.PageIndex,
			PageSize:  search.PageSize,
		}, c)
	}
}

// ExportExtractAgency
// @Tags special-report(特殊报表)
// @Summary 特殊报表--机构抽取列表导出
// @Auth xingqiyi
// @Date 2021/11/3 3:54 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param startTime      query   string   true   "开始时间今天开始时间格式'2022-07-06T16:00:00.000Z'"
// @Param endTime        query   string   true   "结束时间今天结束时间格式'2022-07-06T16:00:00.000Z'"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /special-report/extract-agency/export [get]
func ExportExtractAgency(c *gin.Context) {
	var search model.BaseTimeRangeWithCode
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list := service2.ExportExtractAgency(search)
	path := global.GConfig.LocalUpload.FilePath + global.PathProReportExport + search.ProCode + "/"
	// 尝试创建此路径
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		global.GLog.Error("upload file fail:", zap.Any("err", err))
		return
	}
	name := fmt.Sprintf("理赔2.0机构抽取-%v.xlsx", time.Now().Format("20060102"))
	err = utils.ExportBigExcel(path, name, "sheet", list)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.FileAttachment(path+name, name)
}

// ExportExtractAgencyRealTime
// @Tags special-report(特殊报表)
// @Summary 特殊报表--实时机构抽取列表导出
// @Auth xingqiyi
// @Date 2021/11/3 3:54 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param startTime      query   string   true   "开始时间今天开始时间格式'2022-07-06T16:00:00.000Z'"
// @Param endTime        query   string   true   "结束时间今天结束时间格式'2022-07-06T16:00:00.000Z'"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /special-report/extract-agency/export-real-time [get]
func ExportExtractAgencyRealTime(c *gin.Context) {
	var search model.BaseTimeRangeWithCode
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list := service2.FetchBills(search)

	var agencyArr []model3.Agency
	for _, bill := range list {
		//获取xml
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		dir = strings.Replace(dir, "bin", "", 1)
		xmlFile := dir + global.GConfig.LocalUpload.FilePath + bill.ProCode + "/upload_xml/" +
			fmt.Sprintf("%v/%v/%v/%v.xml",
				bill.CreatedAt.Year(), int(bill.CreatedAt.Month()),
				bill.CreatedAt.Day(),
				bill.BillNum)
		global.GLog.Info(xmlFile)
		file, err := os.ReadFile(xmlFile)
		if err != nil {
			global.GLog.Error("ReadFile fail:", zap.Any("err", err))
			//return
		}

		//机构抽取 fetchAgency
		agency := fetchAgency(bill, string(file))
		if err != nil {
			global.GLog.Error("FetchAgency fail:", zap.Any("err", err))
			//return
		}
		agencyArr = append(agencyArr, agency)
	}

	path := global.GConfig.LocalUpload.FilePath + global.PathProReportExport + search.ProCode + "/"
	// 尝试创建此路径
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		global.GLog.Error("upload file fail:", zap.Any("err", err))
		return
	}
	name := fmt.Sprintf("理赔2.0机构抽取-%v.xlsx", time.Now().Format("20060102"))
	err = utils.ExportBigExcel(path, name, "sheet", agencyArr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.FileAttachment(path+name, name)
}

var expendModeMap = map[string]string{"0020": "门诊", "0040": "住院"}
var isMatchMap = map[bool]string{true: "是", false: "否"}
var constNameMap = map[string]string{"B0103": "B0103_广西贵州国寿理赔_ICD10疾病编码", "B0106": "B0106_陕西国寿理赔_ICD10疾病编码", "B0110": "B0110_新疆国寿理赔_ICD10疾病编码"}

func fetchAgency(bill model2.ProjectBill, xmlValue string) model3.Agency {
	//常量
	//constMap := constDeal(obj.Bill.ProCode)
	mIcd10CodeName := ""
	isMatch := ""
	mIcd10Codes := utils.GetNodeData(xmlValue, "mIcd10Code")
	for _, code := range mIcd10Codes {
		//name, ok := constMap["jiBingBianMaCodeToNameMap"][code]
		name, total := utils.FetchConst(bill.ProCode, constNameMap[bill.ProCode], "疾病名称", map[string]string{"疾病代码": code})
		mIcd10CodeName += name + ";"
		//isMatch += isMatchMap[ok] + ";"
		isMatch += isMatchMap[total != 0] + ";"
	}

	expenModeName := ""
	expenModes := utils.GetNodeData(xmlValue, "expenMode")
	for _, mode := range expenModes {
		if val, ok := expendModeMap[mode]; ok {
			expenModeName += val + ";"
		}
	}

	rcptAmnt := strings.Join(utils.GetNodeData(xmlValue, "rcptAmnt"), ";")
	socialPayAmnt := strings.Join(utils.GetNodeData(xmlValue, "socialPayAmnt"), ";")

	agency := model3.Agency{
		BillId:         bill.ID,
		BillNum:        bill.BillNum,
		MIcd10Code:     mIcd10CodeName,
		Agency:         bill.Agency,
		IsMatch:        isMatch,
		ExpenMode:      expenModeName,
		CountMoney:     rcptAmnt,
		SocialPayMoney: socialPayAmnt,
	}

	return agency
}
