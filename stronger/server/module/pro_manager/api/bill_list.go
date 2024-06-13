/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 10:34 上午
 */

package api

import (
	"errors"
	"fmt"
	xj "github.com/basgys/goxml2json"
	"io/ioutil"
	"os"
	"server/global"
	"server/global/response"
	"server/module/export"
	model4 "server/module/export/model"
	"server/module/export/project"
	service2 "server/module/export/service"
	"server/module/load"
	model2 "server/module/load/model"
	model3 "server/module/pro_conf/model"
	proconf4 "server/module/pro_conf/model/response"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"

	//proconf4 "server/module/pro_conf/model/response"
	"server/module/pro_manager/const_data"
	"server/module/pro_manager/model"
	"server/module/pro_manager/service"
	responseSysBase "server/module/sys_base/model/response"
	"server/module/upload"
	"server/utils"
	"strconv"
	"strings"
	"sync"

	"github.com/axgle/mahonia"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GetBillByPage
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--获取案件列表
// @Auth xingqiyi
// @Date 2021/11/3 3:54 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目代码"
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param timeStart      query   string   false   "开始时间今天开始时间格式'2021-11-05 15:04:05'"
// @Param timeEnd        query   string   false   "结束时间今天结束时间格式'2021-11-05 15:04:05'"
// @Param billCode       query   string   false   "案件号"
// @Param status         query   int      false   "案件状态"
// @Param saleChannel    query   string   false   "案件状态"
// @Param batchNum       query   string   false   "批次号"
// @Param agency         query   string   false   "机构号"
// @Param insuranceType  query   string      false   "医保类型"
// @Param claimType      query   int      false   "理赔类型"
// @Param stickLevel     query   int      false   "加急件"
// @Param minCountMoney  query   float32  false   "最小账单金额"
// @Param maxCountMoney  query   float32  false   "最大账单金额"
// @Param isQuestion     query   int      false   "是否是问题件0:不筛选，1：true，2：false"
// @Param invoiceNum     query   int      false   "发票数量"
// @Param qualityUser    query   string   false   "质检人"
// @Param stage          query   string   false   "录入状态"
// @Param orderBy        query   string   false   "排序JSON.stringify([["CreatedAt","desc"]])"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/page [get]
func GetBillByPage(c *gin.Context) {
	var billListSearch model.BillListSearch
	err := c.ShouldBindQuery(&billListSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, total, bills := service.GetBillByPage(billListSearch)
	//utils.DictConvert(bills)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      bills,
			Total:     total,
			PageIndex: billListSearch.PageIndex,
			PageSize:  billListSearch.PageSize,
		}, c)
	}
}

// DelByIdsAndProCode
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--根据id删除（软删除）
// @Auth xingqiyi
// @Date 2021年11月08日16:27:32
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.DelByIdAndProCode true "根据id删除"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pro-manager/bill-list/delete [delete]
func DelByIdsAndProCode(c *gin.Context) {
	var reqParam model.DelByIdAndProCode
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	//删除内存库相关的表信息  将数据保存到历史库恢复时候要用

	//查询内存库返回obj
	err, billObj := service.GetTaskById(reqParam)
	if err != nil {
		response.FailWithMessage("查询任务库失败", c)
		return
	}

	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	//更新内存库的数据到历史库
	err = service.UpdateBillObjInDel(reqParam, billObj, customClaims)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}

	//更新历史库
	//err, rows := service.DelByIdsAndProCode(reqParam)
	//if err != nil || rows == 0 {
	//	response.FailWithMessage("更新历史库失败", c)
	//	return
	//}

	//删除内存库obj
	err, rmRow := service.DelTaskById(reqParam)
	global.GLog.Error("", zap.Int64("", rmRow))
	if err != nil {
		response.FailWithMessage("删除失败", c)
	} else {
		backEndPort := global.GProConf[reqParam.ProCode].BackEndPort
		innerIp := global.GProConf[reqParam.ProCode].InnerIp
		err, res := utils.HttpRequest("http://"+innerIp+":"+strconv.Itoa(backEndPort)+"/task/releaseBill", map[string]interface{}{
			"id": reqParam.ID,
		})
		if err != nil {
			global.GLog.Error("通知失败::" + ":::" + err.Error())
			global.GLog.Error("", zap.Any("", res))
			response.OkWithMessage("释放成功，通知失败", c)
			return
		}
		//插入修改结果数据日志
		err = SetLog(c, billObj.ProjectBill.ID, reqParam.ProCode, "案件状态::删除", const_data.BillStatus[billObj.ProjectBill.Status], "已删除")
		if err != nil {
			global.GLog.Error("删除::" + global.LogErr.Error() + ":::" + err.Error())
		}
		response.OkDetailed(responseSysBase.RowResult{
			Row: rmRow,
		}, "删除成功", c)
	}
}

// GetDictConst
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--获取列表的常量
// @Auth xingqiyi
// @Date 2021年11月08日16:27:27
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/dict-const [get]
func GetDictConst(c *gin.Context) {
	obj := map[string]interface{}{
		"billStatus":         const_data.BillStatus,
		"billInsuranceType":  const_data.BillInsuranceType,
		"billClaimType":      const_data.BillClaimType,
		"billStage":          const_data.BillStage,
		"billStickLevel":     const_data.BillStickLevel,
		"announcementType":   const_data.AnnouncementType,
		"announcementStatus": const_data.AnnouncementStatus,
		"constType":          const_data.ConstType,
		"billType":           const_data.BillType,
		"businessMsgType":    const_data.BusinessMsgType,
		"contractClaimType":  const_data.ContractClaimType,
	}
	response.OkDetailed(obj, "查询成功", c)
}

// RecoverBill
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--恢复单据
// @Auth xingqiyi
// @Date 2021年11月08日16:27:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProCodeAndId true "项目和单id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"恢复成功"}"
// @Router /pro-manager/bill-list/recover [post]
func RecoverBill(c *gin.Context) {
	var reqParam model.ProCodeAndId
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	//将数据从历史库搬到内存库
	err, obj := service.GetBillObj(reqParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("恢复失败%v", err), c)
	}
	fmt.Println(obj)

	//新增obj到内存库
	err = service.InsertTaskBillObj(reqParam, obj)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}

	//更新历史bill库
	//删除历史库block field信息
	err = service.RecoverBill(reqParam)

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("恢复失败%v", err), c)
	} else {
		//插入修改结果数据日志
		err = SetLog(c, reqParam.ID, reqParam.ProCode, "案件状态::恢复", "已删除", "恢复")
		if err != nil {
			global.GLog.Error("恢复::" + global.LogErr.Error())
		}
		response.OkWithMessage("恢复成功", c)
	}
}

// ExportErrBill
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--导出异常
// @Auth xingqiyi
// @Date 2021年11月08日16:27:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProCodeAndId true "项目和单id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"导出异常成功"}"
// @Router /pro-manager/bill-list/export-err-bill [post]
func ExportErrBill(c *gin.Context) {
	var reqParam model.ProCodeAndId
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err = export.WrongBill(reqParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("导出异常失败%v", err), c)
	} else {
		//插入修改结果数据日志
		err = SetLog(c, reqParam.ID, reqParam.ProCode, "案件状态::导出异常", "正常", "异常")
		if err != nil {
			global.GLog.Error("导出异常::" + global.LogErr.Error())
		}
		response.OkWithMessage("导出异常成功", c)
	}
}

// ForceExport
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--强制导出
// @Auth xingqiyi
// @Date 2021年11月08日16:27:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.DelByIdAndProCode true "项目和单id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"强制导出成功"}"
// @Router /pro-manager/bill-list/force-export [post]
func ForceExport(c *gin.Context) {
	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	var reqParam model.DelByIdAndProCode
	err = c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	//获取history单据信息
	err, b := service.GetBillInfo(reqParam)
	if err != nil {
		response.FailWithMessage("查询单据失败", c)
		return
	}
	if b.Stage != 1 && b.Stage != 2 {
		response.FailWithMessage("单据不是待加载、录入中状态", c)
		return
	}

	//录入状态为待加载、录入中时，显示此按钮，
	//点击此按钮后，录入状态更改为待审核，案件状态更改为正常，手动回传更改为是，
	//强制导出数据时不进行校验直接根据最新的录入数据进行导出，
	//并增加导出校验：强制导出信息(工号姓名)：该案件为强制导出案件，请检查并修改数据。

	//查询内存库返回obj
	err, billObj := service.GetTaskById(reqParam)
	if err != nil {
		response.FailWithMessage("查询任务库失败", c)
		return
	}

	//计算产量
	err = export.AbnormalityExport(reqParam.ProCode, billObj.ProjectBill.ID, "1")
	if err != nil {
		response.FailWithMessage("计算产量失败", c)
		return
	}
	//将数据从内存库搬到历史库
	//更新内存库的数据到历史库
	billObj.ProjectBill.WrongNote = customClaims.NickName + "该案件为强制导出案件，请检查并修改数据。"
	err = service.UpdateBillObjInForceExport(reqParam, billObj)
	if err != nil {
		response.FailWithMessage("删除失败", c)
		return
	}

	//删除内存库obj
	err, _ = service.DelTaskById(reqParam)

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("强制导出失败%v", err), c)
	} else {
		//插入修改结果数据日志
		err = SetLog(c, reqParam.ID, reqParam.ProCode, "案件状态::强制导出", const_data.BillStage[b.Stage], "已导出")
		if err != nil {
			global.GLog.Error("强制导出::" + global.LogErr.Error())
		}
		response.OkWithMessage("强制导出成功", c)
	}
}

// SaveXml
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--保存xml
// @Auth xingqiyi
// @Date 2021年11月08日16:27:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.XmlDataAndPath true "xml路径和数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存xml成功"}"
// @Router /pro-manager/bill-list/save-xml [post]
func SaveXml(c *gin.Context) {
	var reqParam model.XmlDataAndPath
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	path := reqParam.Url[:strings.LastIndexAny(reqParam.Url, "/")+1]
	global.GLog.Info(path)
	exists, err := utils.PathExists(path)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("保存xml失败%v", err), c)
		return
	}
	if !exists {
		err = utils.CreateDir(path)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("保存xml失败%v", err), c)
			return
		}
	}

	err, exportConfig := service2.GetExportConf(reqParam.ProCode)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("保存xml失败%v", err), c)
		return
	}
	xmlType := exportConfig.XmlType
	if utils.RegIsMatch("^(B0113|B0121)$", reqParam.ProCode) {
		xmlType = "utf-8"
	}
	enc := mahonia.NewEncoder(xmlType)
	encStr := enc.ConvertString(reqParam.Data)
	err = os.WriteFile(reqParam.Url, []byte(encStr), 0666)
	if err != nil {
		global.GLog.Error("写入失败", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("保存xml失败%v", err), c)
		return
	}

	if xmlType == "utf-8" && utils.RegIsMatch("^(B0103|B0110|B0106|B0122)$", reqParam.ProCode) {
		global.GLog.Info("", zap.Any("", "写入json"))
		xmlReader := strings.NewReader(encStr)
		result, err := xj.Convert(xmlReader)
		if err != nil {
			global.GLog.Error("写入json失败", zap.Error(err))
			response.FailWithMessage(fmt.Sprintf("写入json失败%v", err), c)
			return
		}
		data := []byte("")

		err, data = project.BillExportJsonDealAdapter(reqParam.ProCode, nil, result.Bytes())
		if err != nil {
			global.GLog.Error("写入json失败", zap.Error(err))
			response.FailWithMessage(fmt.Sprintf("json定制处理失败%v", err), c)
			return
		}

		err = os.WriteFile(strings.Replace(reqParam.Url, ".xml", ".json", -1), data, 0666)
	} else {
		global.GLog.Error("该项目不导出json")
	}

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("保存xml失败%v", err), c)
	} else {
		response.OkWithMessage("保存xml成功", c)
	}
}

// Upload
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--回传
// @Auth xingqiyi
// @Date 2021年11月08日16:27:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProCodeAndId true "项目和单id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"回传成功"}"
// @Router /pro-manager/bill-list/upload [post]
func Upload(c *gin.Context) {
	var reqParam model.ProCodeAndId
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, uploadPaths := service.GetUploadPath(reqParam)
	if err != nil || uploadPaths.Upload == "" {
		response.FailWithMessage(fmt.Sprintf("回传失败,回传配置为空%v", err), c)
		return
	}
	err = upload.BillUploadAdapter(reqParam, uploadPaths)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("回传失败，%v", err), c)
	} else {
		//插入修改结果数据日志
		err = SetLog(c, reqParam.ID, reqParam.ProCode, "回传", "", "已回传")
		if err != nil {
			global.GLog.Error("回传::" + global.LogErr.Error())
		}
		response.OkWithMessage("回传成功", c)
	}
}

// SetUploadType
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--设置回传方式，自动:true,手动:false
// @Auth xingqiyi
// @Date 2021年11月08日16:27:21
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.AutoUpload true "AutoUpload"
// @Success 200 {string} string "{"success":true,"data":1,"msg":"修改成功"}"
// @Router /pro-manager/bill-list/set-upload-type [post]
func SetUploadType(c *gin.Context) {
	var reqParam model.AutoUpload
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, b := service.SetUploadType(reqParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("修改失败%v", err), c)
	} else {
		//插入修改结果数据日志
		after := "手动回传"
		if reqParam.IsAutoUpload {
			after = "自动回传"
		}
		before := "手动回传"
		if b.IsAutoUpload {
			before = "自动回传"
		}
		err = SetLog(c, reqParam.ID, reqParam.ProCode, "修改回传方式", before, after)
		if err != nil {
			global.GLog.Error("修改回传方式::" + global.LogErr.Error())
		}
		response.OkWithMessage("修改成功", c)
	}
}

// IsQuality
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--是否有人在质检
// @Auth xingqiyi
// @Date 2022/3/25 14:26
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.EditBillResultDataManyFields true "EditBillResultDataManyFields"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/is-quality [post]
func IsQuality(c *gin.Context) {
	var reqParam model.EditBillResultDataManyFields
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	param := model.ProCodeAndId{
		ID:      reqParam.BillId,
		ProCode: reqParam.ProCode,
	}
	err, bill := service.GetProBillById(param)
	if err != nil {
		response.FailWithMessage("获取单据失败", c)
		return
	}
	//获取登录者
	customClaims, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取登录者失败，%v", err), c)
		return
	}
	//里面有一个人进去了，后面的人拦截一下后，点击确认还是可以给他进去
	//回传不用执行以上需求，直接点进去就可以看，不用提示
	if bill.Stage != 5 {
		if (bill.QualityUserCode != "" || bill.QualityUserName != "") &&
			bill.QualityUserCode != customClaims.Code {
			response.OkWithData(false, c)
			return
		}
	}
	response.OkWithData(true, c)
}

// EditBillResultData
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--编辑结果数据
// @Auth xingqiyi
// @Date 2021年11月30日09:40:52
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.EditBillResultDataManyFields true "EditBillResultDataManyFields"
// @Success 200 {object} model4.ResultDataBill "{"success":true,"data":model4.ResultDataBill,"msg":"编辑成功"}"
// @Router /pro-manager/bill-list/edit-bill-result-data [post]
func EditBillResultData(c *gin.Context) {
	var reqParam model.EditBillResultDataManyFields
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	err2, resultDataBill := BillResultDataFunc(reqParam, c, false)
	if err2 != "" {
		response.FailWithMessage(err2, c)
		return
	} else {
		response.OkWithData(resultDataBill, c)
		return
	}
	// return

	//var f model2.ProjectField
	//var val string
	//if len(reqParam.Fields) > 0 {
	//获取更新前字段的值
	//err, f = service.GetField(reqParam)
	//if err != nil {
	//	response.FailWithMessage(fmt.Sprintf("获取field失败%v", err), c)
	//	return
	//}

	//更新fields
	//err = service.UpdateField(reqParam)
	//if err != nil {
	//	response.FailWithMessage(fmt.Sprintf("更新fields%v", err), c)
	//	return
	//}
	//val = reqParam.Fields[0].FieldValue
	//}
	//获取单据信息
	//param := model.EditBillResultData{
	//	BillId:  reqParam.BillId,
	//	ProCode: reqParam.ProCode,
	//}
	// err, bill, blocks, fields, fieldsLen := service.GetBillAndBlocks(reqParam)
	// if err != nil || bill.ID == "" {
	// 	response.FailWithMessage(fmt.Sprintf("获取单据、分块、字段%v", err), c)
	// 	return
	// }

	// if reqParam.EditType == 2 {
	// 	if fieldsLen != int64(len(reqParam.Fields)) {
	// 		response.FailWithMessage(fmt.Sprintf("前端处理数据有问题,正常字段个数：%d，传到后端的个数%d", fieldsLen, len(reqParam.Fields)), c)
	// 		return
	// 	}
	// 	fields = reqParam.Fields
	// }
	// //if bill.Stage != 5 {
	// //添加质检人
	// customClaims, err := api.GetUserByToken(c)
	// if err != nil {
	// 	response.FailWithMessage(fmt.Sprintf("获取登录者失败，%v", err), c)
	// 	return
	// }
	// //if (bill.QualityUserCode != "" || bill.QualityUserName != "") &&
	// //	bill.QualityUserCode != customClaims.Code {
	// //	response.FailWithMessage(bill.QualityUserCode+bill.QualityUserName+"正在质检，请注意", c)
	// //	return
	// //}
	// if bill.QualityUserCode == "" && bill.QualityUserName == "" {
	// 	err = service.UpdateQualityUser(customClaims, reqParam.ProCode, reqParam.BillId)
	// 	if err != nil {
	// 		response.FailWithMessage(fmt.Sprintf("更新质检人失败，%v", err), c)
	// 		return
	// 	}
	// }
	// //}
	// //获取字段配置
	// err, sysProFields := service.GetFieldConf(reqParam.ProCode)
	// if err != nil {
	// 	response.FailWithMessage(fmt.Sprintf("获取字段配置%v", err), c)
	// 	return
	// }

	// //数据转换处理，生成结果数据
	// resultDataBill, err := export.ChangeData(bill, blocks, sysProFields, fields)
	// resultDataBill.Bill.WrongNote = bill.WrongNote
	// if err != nil {
	// 	response.FailWithMessage(fmt.Sprintf("生成结果数据%v", err), c)
	// 	return
	// }

	// if err != nil {
	// 	response.FailWithMessage(fmt.Sprintf("编辑失败%v", err), c)
	// } else {
	// 	//插入修改结果数据日志
	// 	for _, f := range reqParam.EditFields {
	// 		err = SetLog(c, reqParam.BillId, reqParam.ProCode, f.Code+f.Name, f.BeforeVal, f.EndVal)
	// 		if err != nil {
	// 			global.GLog.Error("修改字段值::" + global.LogErr.Error() + ":::" + err.Error())
	// 		}
	// 	}
	// 	response.OkWithData(resultDataBill, c)
	// 	//response.OkWithData(rep.PageResult{
	// 	//	PageSize:  1000,
	// 	//	PageIndex: reqParam.BasePageInfo.PageIndex,
	// 	//	Total:     int64(len(resultDataBill.Invoice)),
	// 	//	List: model4.ResultDataBill{
	// 	//		Bill:    resultDataBill.Bill,
	// 	//		Invoice: []model4.InvoiceMap{resultDataBill.Invoice[reqParam.BasePageInfo.PageIndex]},
	// 	//	},
	// 	//}, c)
	// }
}

func BillResultDataFunc(reqParam model.EditBillResultDataManyFields, c *gin.Context, check bool) (string, model4.ResultDataBill) {
	err, bill := service.GetBill(reqParam)
	if err != nil || bill.ID == "" {
		return fmt.Sprintf("获取单据失败%v", err), model4.ResultDataBill{}
	}

	if bill.Stage != 5 && bill.Stage != 7 {
		//添加质检人
		customClaims, err := api2.GetUserByToken(c)
		if err != nil {
			return fmt.Sprintf("获取登录者失败，%v", err), model4.ResultDataBill{}
		}

		if bill.QualityUserCode == "" && bill.QualityUserName == "" {
			err = service.UpdateQualityUser(customClaims, reqParam.ProCode, reqParam.BillId, bill.Stage, bill.ExportAt)
			if err != nil {
				return fmt.Sprintf("更新质检人失败，%v", err), model4.ResultDataBill{}
			}
		}
	}

	err, blocks, fields, fieldsLen := service.GetBlockAndFields(reqParam)
	if err != nil {
		return fmt.Sprintf("获取分块、字段%v", err), model4.ResultDataBill{}
	}

	if reqParam.EditType == 2 {
		if fieldsLen != int64(len(reqParam.Fields)) {
			return fmt.Sprintf("前端处理数据有问题,正常字段个数：%d，传到后端的个数%d", fieldsLen, len(reqParam.Fields)), model4.ResultDataBill{}
		}
		fields = reqParam.Fields
	}
	//获取字段配置
	err, sysProFields := service.GetFieldConf(reqParam.ProCode)
	if err != nil {
		return fmt.Sprintf("获取字段配置%v", err), model4.ResultDataBill{}
	}

	if check && !strings.Contains(bill.Remark, "<<发票查验>>") {
		bill.Remark += "<<发票查验>>"
	}

	//数据转换处理，生成结果数据
	resultDataBill, err := export.ChangeData(bill, blocks, sysProFields, fields)
	resultDataBill.Bill.WrongNote += bill.WrongNote
	if err != nil {
		return fmt.Sprintf("生成结果数据%v", err), model4.ResultDataBill{}
	}

	if err != nil {
		return fmt.Sprintf("编辑失败%v", err), model4.ResultDataBill{}
	} else {
		//插入修改结果数据日志
		for _, f := range reqParam.EditFields {
			err = SetLog(c, reqParam.BillId, reqParam.ProCode, f.Code+f.Name, f.BeforeVal, f.EndVal)
			if err != nil {
				global.GLog.Error("修改字段值::" + global.LogErr.Error() + ":::" + err.Error())
			}
		}
		// response.OkWithData(resultDataBill, c)
		return "", resultDataBill
	}
}

func ZbjBillResultData(c *gin.Context) {
	var reqParam model.EditBillResultDataManyFields
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	err2, resultDataBill := BillResultDataFunc(reqParam, c, true)
	if err2 != "" {
		response.FailWithMessage(err2, c)
		return
	} else {
		err := ""
		wrongNets := strings.Split(resultDataBill.Bill.WrongNote, ";")
		for _, wrongNet := range wrongNets {
			if strings.Index(wrongNet, "查验失败：") != -1 {
				err += wrongNet
			}
		}

		//去除飘窗前的wrongNet
		if indexErrStr := strings.Index(err, "查验失败："); indexErrStr != -1 {
			err = err[indexErrStr:]
		}

		if err != "" {
			response.FailWithoutCode(resultDataBill, err, c)
			return
		}

		response.OkWithData(resultDataBill, c)
		return
	}
}

// SaveBillResultData
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--保存结果数据后的单据
// @Auth xingqiyi
// @Date 2021年11月30日09:40:52
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.EditBillResultDataManyFields true "EditBillResultData"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /pro-manager/bill-list/save-bill-result-data [post]
func SaveBillResultData(c *gin.Context) {
	var reqParam model.EditBillResultDataManyFields
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	reqParam.EditType = 2
	//获取分块
	err, bill, blocks, fieldsLen := service.GetBillAndBlocks(reqParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取分块%v", err), c)
		return
	}
	fields := reqParam.Fields
	if fieldsLen != int64(len(fields)) {
		response.FailWithMessage(fmt.Sprintf("前端处理数据有问题，正常字段个数：%d，传到后端的个数%d", fieldsLen, len(reqParam.Fields)), c)
		return
	}
	//获取字段配置
	err, sysProFields := service.GetFieldConf(reqParam.ProCode)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取字段配置%v", err), c)
		return
	}
	flag := "1"
	obj, err := export.BillExport(bill, blocks, sysProFields, fields, flag)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("导出失败%v", err), c)
		return
	}
	//更新字段到数据库
	//service.UpdateField(reqParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("保存成功%v", err), c)
	} else {
		response.OkWithData(obj.Bill.WrongNote, c)
	}
}

// SeeBillResultData
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--查看格式化后的单据
// @Auth xingqiyi
// @Date 2021年12月08日13:55:15
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode   query   string   true    "项目代码"
// @Param id 		query 	string   true    "单据id"
// @Success 200 {object} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/see-bill-result-data [get]
func SeeBillResultData(c *gin.Context) {
	var proCodeAndId model.ProCodeAndId
	err := c.ShouldBindQuery(&proCodeAndId)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	//获取项目质检配置
	err, s := service.GetQualitiesByPro(proCodeAndId.ProCode)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取质检配置失败，%v", err), c)
		return
	}

	//获取单据xml文件
	err, b := service.GetBillInfo(model.DelByIdAndProCode{ProCode: proCodeAndId.ProCode, ID: proCodeAndId.ID})
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取单据信息失败，%v", err), c)
		return
	}

	//格式化查看结果数据xml质检
	err, data, amount := formatBill(s, b)

	t := make(map[string]interface{}, 2)
	t["amount"] = amount
	t["data"] = data

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkWithData(t, c)
	}
}

// GetFieldInfo
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--获取字段、分块和字段配置信息
// @Auth xingqiyi
// @Date 2021年12月02日09:22:43
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode   query   string   true    "项目代码"
// @Param id 		query 	string   true    "字段id"
// @Success 200 {object} model.FieldObj "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/get-field-info/ [get]
func GetFieldInfo(c *gin.Context) {
	var proCodeAndId model.ProCodeAndId
	err := c.ShouldBindQuery(&proCodeAndId)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, obj := service.GetFieldInfo(proCodeAndId)
	obj.Block.Op0Code += global.UserCodeName[obj.Block.Op0Code]
	obj.Block.Op1Code += global.UserCodeName[obj.Block.Op1Code]
	obj.Block.Op2Code += global.UserCodeName[obj.Block.Op2Code]
	obj.Block.OpqCode += global.UserCodeName[obj.Block.OpqCode]
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkWithData(obj, c)
	}
}

// EditFeedbackVal
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--修改反馈值
// @Auth xingqiyi
// @Date 2021年12月02日15:08:10
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.EditFeedbackVal true "反馈结构体"
// @Success 200 {object} model.EditFeedbackVal "{"success":true,"data":model.EditFeedbackVal,"msg":"编辑成功"}"
// @Router /pro-manager/bill-list/edit-feedback-val [post]
func EditFeedbackVal(c *gin.Context) {
	var reqParam model.EditFeedbackVal
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	err = service.UpdateFeedback(reqParam)

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("保存成功%v", err), c)
	} else {
		response.OkWithMessage("保存成功", c)
	}
}

// SetPractice
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--设置分块为练习
// @Auth xingqiyi
// @Date 2021年12月02日15:08:10
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SetPracticeForm true "设置练习form"
// @Success 200 {string} string "{"success":true,"data":"","msg":"编辑成功"}"
// @Router /pro-manager/bill-list/set-practice [post]
func SetPractice(c *gin.Context) {
	var reqParam model.SetPracticeForm
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	err = service.SetPractice(reqParam)

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("设置成功%v", err), c)
	} else {
		response.OkWithMessage("设置成功", c)
	}
}

// GetLog
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--获取单据的修改日志
// @Auth xingqiyi
// @Date 2021年12月02日15:08:10
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id query string true "单据id"
// @Param proCode query string true "项目代码"
// @Success 200 {object} []model.ResultDataLog "{"success":true,"data":"","msg":"获取成功"}"
// @Router /pro-manager/bill-list/get-log [get]
func GetLog(c *gin.Context) {
	var reqParam model.ProCodeAndId
	err := c.ShouldBindQuery(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	err, logs := service.GetLog(reqParam)

	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取成功%v", err), c)
	} else {
		response.OkWithData(logs, c)
	}
}

// GetQingDan
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--获取单据的清单信息
// @Auth xingqiyi
// @Date 2021年12月02日15:08:10
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id query string true "单据id"
// @Param proCode query string true "项目代码"
// @Param itemName query string false "项目名称"
// @Param invoiceNum query string false "账单号"
// @Success 200 {object} string "{"success":true,"data":"","msg":"获取成功"}"
// @Router /pro-manager/bill-list/qing-dan/get [get]
func GetQingDan(c *gin.Context) {
	var reqParam model.QingDanForm
	err := c.ShouldBindQuery(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, obj := getResultData(reqParam)

	err, dan := project.GetQingDan(reqParam, obj)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取成功%v", err), c)
	} else {
		response.OkWithData(dan, c)
	}
}

func getResultData(reqParam model.QingDanForm) (err error, obj model4.ResultDataBill) {
	var wg sync.WaitGroup
	var fieldConfList []model3.SysProField
	var bill model.ProjectBill
	var blocks []model2.ProjectBlock
	var fields []model2.ProjectField

	//获取分块
	wg.Add(1)
	go service2.GetBillBlocksAndFields(reqParam.ID, reqParam.ProCode, &wg, &bill, &blocks, &fields)

	//获取字段配置
	wg.Add(1)
	go service2.GetFieldConf(reqParam.ProCode, &wg, &fieldConfList)
	wg.Wait()

	if bill.ID == "" {
		return errors.New("没有找到该单据"), obj
	}

	//生成结果数据
	obj, err = export.ChangeData(bill, blocks, fieldConfList, fields)

	return err, obj
}

// ExportQingDan
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--导出清单
// @Auth xingqiyi
// @Date 2021年12月21日16:35:33
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param id query string true "单据id"
// @Param proCode query string true "项目代码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/qing-dan/export [get]
func ExportQingDan(c *gin.Context) {
	var reqParam model.QingDanForm
	err := c.ShouldBindQuery(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, obj := getResultData(reqParam)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, dan := project.GetQingDan(reqParam, obj)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	path := fmt.Sprintf("%v/%v/%v/", global.GConfig.LocalUpload.FilePath, reqParam.ProCode, global.PathQingDan)
	name := obj.Bill.BillNum + ".xlsx"
	err = utils.ExportBigExcel(path, name, "sheet", dan)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.FileAttachment(path+name, name)
}

// Reload
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--重加载
// @Auth xingqiyi
// @Date 2022/2/24 9:20 上午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ProCodeAndId true "单据id和项目编码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/reload [post]
func Reload(c *gin.Context) {
	var reqParam model.ProCodeAndId
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	//获取单据
	err, b := service.GetProBillById(reqParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取单据失败%v", err), c)
		return
	}
	//清空内存库
	err = service2.DelBill(reqParam.ID, reqParam.ProCode+"_task")
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("清空内存库失败%v", err), c)
		return
	}

	//重加载
	b.Status = 1
	err = load.ProLoadFunc(reqParam.ProCode, b)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("重加载失败%v", err), c)
		return
	}
	//修改状态
	err = service.UploadState(reqParam, 2)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("重加载失败%v", err), c)
	} else {
		response.OkWithData("重加载成功", c)
	}
}

// Remark
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--修改备注
// @Auth xingqiyi
// @Date 2022/2/24 9:20 上午
// @Security ApiKeyAuth
// @Security UserID
// @accept application/json
// @Produce application/json
// @Param data body model.Remark true "修改备注"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/remark [post]
func Remark(c *gin.Context) {
	var reqParam model.Remark
	err := c.ShouldBindJSON(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	//修改备注
	err = service.UploadRemark(reqParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("修改失败%v", err), c)
	} else {
		response.OkWithData("修改成功", c)
	}
}

// GetBlockImg
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--过去字段图片
// @Auth xingqiyi
// @Date 2022/2/25 3:44 下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode   query   string   true    "项目代码"
// @Param id 		query 	string   true    "分块id"
// @Param fieldId 		query 	string   false    "字段id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/get-block-img/ [get]
func GetBlockImg(c *gin.Context) {
	var reqParam model.ProCodeAndId
	err := c.ShouldBindQuery(&reqParam)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, blocks, bill := service.GetBlockById(reqParam)
	for i, block := range blocks {
		if block.Status == 1 {
			//初审要自定义获取
			err, fields := service.GetFieldById(reqParam)
			if err != nil {
				response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
				return
			}
			//先范围再页码 field_index = 2是范围  1是类型  0是页码  有自定义再说
			name := fields[2].ResultValue
			if name == "" {
				index, err2 := strconv.Atoi(fields[0].ResultValue)
				if err2 != nil {
					response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
					return
				}
				name = bill.Pictures[index]
			}
			blocks[i].Picture = strings.Replace(global.GConfig.LocalUpload.FilePath, ".", "", -1) + bill.DownloadPath + name

		} else {
			blocks[i].Picture = strings.Replace(global.GConfig.LocalUpload.FilePath, ".", "", -1) + bill.DownloadPath + block.Picture
		}
	}
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
	} else {
		response.OkWithData(blocks, c)
	}
}

// QQQ
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--获取单据的修改日志
// @Auth xingqiyi
// @Date 2021年12月02日15:08:10
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} []model.ResultDataLog "{"success":true,"data":"","msg":"获取成功"}"
// @Router /pro-manager/bill-list/qqq [get]

func QQQ(c *gin.Context) {
	//var a = service.QQQ{}
	//err := c.ShouldBindQuery(&a)
	//if err != nil {
	//	return
	//}
	//err, logs := service.QQ()
	contentByte, err := ioutil.ReadFile("/Users/mjl/project/stronger/trunk/server/files/B0118/download/2020/04-21/WechatIMG1.JPG")
	//contentByte1 := bufio.NewReader(bytes.NewReader(contentByte))
	//if err != nil {
	//	fmt.Println(err)
	//	response.FailWithMessage(fmt.Sprintf("获取成功%v", err), c)
	//	return
	//}
	//buf := bytes.NewBuffer([]byte{})
	//for _, v := range contentByte {
	//	buf.WriteString(fmt.Sprintf("%08b", v))
	//}
	//encrypt, err := utils.AesEncrypt(contentByte, []byte("xingqiyistronger"))
	//fmt.Println(contentByte1.ReadString('\n'))
	fmt.Println(contentByte)
	//res := base64.StdEncoding.EncodeToString(contentByte)
	//fmt.Println(res)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取成功%v", err), c)
	} else {
		response.OkWithData(contentByte, c)
	}
}

// GetTimeLinessBriefing
// @Tags pro-manager/bill-list(案件列表)
// @Summary 案件列表--时效简报
// @Auth lhc
// @Date 2023/9/7 10:44 上午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode   query   string   true    "项目代码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-manager/bill-list/get-time-liness-briefing/ [get]
func GetTimeLinessBriefing(c *gin.Context) {
	proCode := c.Query("proCode")
	if proCode == "" {
		response.FailWithMessage(fmt.Sprintf("项目编码不能为空."), c)
		return
	}
	err := service.SetTimeLinessBriefing(proCode)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取数据失败, %v", err), c)
		return
	}
	response.OkWithData(proconf4.ProjectConfigAgingResponse{
		List: proCode,
	}, c)
}
