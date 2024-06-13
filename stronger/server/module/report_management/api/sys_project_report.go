package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
	"reflect"
	"server/global"
	"server/global/response"
	model2 "server/module/homepage/model"
	request2 "server/module/homepage/model/request"
	service2 "server/module/homepage/service"
	"server/module/report_management/model/request"
	"server/module/report_management/project"
	"server/module/report_management/service"
	"server/module/sys_base/model"
	responseSysBase "server/module/sys_base/model/response"
	"server/utils"
	"time"
)

// GetBusinessDetails
// @Tags Project Report (项目报表)
// @Summary 项目报表--查询业务明细表的数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"返回成功"}"
// @Router /report-management/project-report/business-details/list [get]
func GetBusinessDetails(c *gin.Context) {
	var search request.BusinessDetailsSearch
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	proVerify := utils.Rules{
		"ProCode":   {utils.NotEmpty()},
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
		"pageIndex": {utils.Gt("0")},
		"pageSize":  {utils.Gt("0")},
	}
	proVerifyErr := utils.Verify(search, proVerify)
	if proVerifyErr != nil {
		response.FailWithMessage(proVerifyErr.Error(), c)
		return
	}
	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", search.StartTime, time.Local)
	endTime, _ := time.ParseInLocation("2006-01-02 15:04:05", search.EndTime, time.Local)
	_, startDate, _ := startTime.Date()
	_, endDate, _ := endTime.Date()
	if endDate-startDate > 3 {
		response.FailWithMessage(fmt.Sprintf("查询失败，每次最长只能查询三个月的数据"), c)
		return
	}

	err, list, total := project.GetBusinessDetails(search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("查询失败，%v", err), c)
	} else {
		response.OkWithData(responseSysBase.PageResult{
			List:      list,
			Total:     total,
			PageIndex: search.PageInfo.PageIndex,
			PageSize:  search.PageInfo.PageSize,
		}, c)
	}
}

// ExportBusinessDetails
// @Tags Project Report (项目报表)
// @Summary	项目报表--导出业务明细表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间"
// @Param endTime query string true "结束时间"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"保存成功"}"
// @Router /report-management/project-report/business-details/export [get]
func ExportBusinessDetails(c *gin.Context) {
	var search request.BusinessDetailsSearch
	if err := c.ShouldBindWith(&search, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	proVerify := utils.Rules{
		"ProCode":   {utils.NotEmpty()},
		"StartTime": {utils.NotEmpty()},
		"EndTime":   {utils.NotEmpty()},
	}
	proVerifyErr := utils.Verify(search, proVerify)
	if proVerifyErr != nil {
		response.FailWithMessage(proVerifyErr.Error(), c)
		return
	}
	//uid := c.Request.Header.Get("x-user-id")
	//if uid == "" {
	//	response.FailWithMessage("x-user-id is empty", c)
	//	return
	//}
	//err, a := utils.GetRedisExport("projectReport", uid)
	//if err == nil && a == "true" {
	//	response.FailWithMessage("正在导出!", c)
	//	return
	//}
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(fmt.Sprintf("项目报表--导出业务明细表导出失败, err: %s", r))
	//			global.GLog.Error(fmt.Sprintf("项目报表--导出业务明细表导出失败, err: %s", r))
	//			err = utils.DelRedisExport("projectReport", uid)
	//			if err != nil {
	//				fmt.Println(fmt.Sprintf("项目报表--导出业务明细表删除导出缓存失败, err: %s", err.Error()))
	//				global.GLog.Error(fmt.Sprintf("项目报表--导出业务明细表删除导出缓存失败, err: %s", err.Error()))
	//			}
	//		}
	//	}()
	//	//可以广播同一个登录人的客户端的写法
	//	//global.GSocketIo.BroadcastToRoom("/global-export", uid, "outputStatistics", "所要发送的消息")
	//	err := utils.SetRedisExport("projectReport", uid)
	//	if err != nil {
	//		global.GSocketConnMap[uid].Emit("projectReport", response.ExportResponse{
	//			Code: 400,
	//			Data: "",
	//			Msg:  fmt.Sprintf("projectReport SetRedisExport err: %s", err.Error()),
	//		})
	//		err = utils.DelRedisExport("projectReport", uid)
	//		return
	//	}
	//	path := ""
	//	switch Search.ProCode {
	//	case "B0118":
	//		err, path = B0118.ExportBusinessDetails(Search)
	//		if err != nil {
	//			global.GSocketConnMap[uid].Emit("projectReport", response.ExportResponse{
	//				Code: 400,
	//				Data: "",
	//				Msg:  fmt.Sprintf("ExportBusinessDetails err: %s", err.Error()),
	//			})
	//			err = utils.DelRedisExport("projectReport", uid)
	//			return
	//		}
	//	default:
	//		fmt.Println("123")
	//	}
	//	global.GSocketConnMap[uid].Emit("projectReport", response.ExportResponse{
	//		Code: 200,
	//		Data: path,
	//		Msg:  "导出完成!",
	//	})
	//	err = utils.DelRedisExport("projectReport", uid)
	//	return
	//}()

	//if err != nil {
	//	global.GLog.Error(err.Error())
	//	response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	//} else {
	//
	//}

	err, path, name := project.ExportBusinessDetails(search)
	if err != nil {
		global.GLog.Error(err.Error())
		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	} else {
		c.FileAttachment(path, name)
	}
	//switch Search.ProCode {
	//case "B0118":
	//	err, path, name := B0118.ExportBusinessDetails(Search)
	//	if err != nil {
	//		global.GLog.Error(err.Error())
	//		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	//	} else {
	//		c.FileAttachment(path, name)
	//	}
	//case "B0108":
	//	err, path, name := B0108.ExportBusinessDetails(Search)
	//	if err != nil {
	//		global.GLog.Error(err.Error())
	//		response.FailWithMessage(fmt.Sprintf("导出失败，%v", err), c)
	//	} else {
	//		c.FileAttachment(path, name)
	//	}
	//
	//default:
	//	response.FailWithMessage(fmt.Sprintf("没有该项目"), c)
	//}
	//response.Ok(c)
}

// GetBusinessReport
// @Tags project-report(报表管理--项目报表)
// @Summary 报表管理--获取项目业务量报表(日、周、月业务量)
// @Auth xingqiyi
// @Date 2022/7/28 14:03
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param 	proCode     		query   string   true    "all：全部，B0118：B0118等"
// @Param 	type     			query   int   	 true    "1：日，2：周，3：月，4：年"
// @Param 	startTime      		query   string   true   "开始时间格式2022-07-06T16:00:00.000Z"
// @Param 	endTime        		query   string   true   "结束时间格式2022-07-06T16:00:00.000Z"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/project-report/business-report [get]
func GetBusinessReport(c *gin.Context) {
	var businessReportReq request.BusinessReportReq
	err := c.ShouldBindQuery(&businessReportReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	verify := utils.Rules{
		"Type": {utils.Gt("0")},
	}
	verifyErr := utils.Verify(businessReportReq, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	err, obj := service.GetProReport(businessReportReq, 1)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithoutCode(err.Error(), fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List: obj,
		}, c)
	}
}

// GetAgingReport
// @Tags project-report(报表管理--项目报表)
// @Summary 报表管理--获取项目时效报表(日、周、月)
// @Auth xingqiyi
// @Date 2022/7/29 09:22
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param 	proCode     		query   string   true    "all：全部，B0118：B0118等"
// @Param 	type     			query   int   	 true    "1：日，2：周，3：月，4：年"
// @Param 	startTime      		query   string   true   "开始时间格式2022-07-06T16:00:00.000Z"
// @Param 	endTime        		query   string   true   "结束时间格式2022-07-06T16:00:00.000Z"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/project-report/aging-report [get]
func GetAgingReport(c *gin.Context) {
	var agingReportReq request.BusinessReportReq
	err := c.ShouldBindQuery(&agingReportReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	verify := utils.Rules{
		"Type": {utils.Gt("0")},
	}
	verifyErr := utils.Verify(agingReportReq, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	err, obj := service.GetProReport(agingReportReq, 2)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithoutCode(err.Error(), fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List: obj,
		}, c)
	}
}

// GetDealTimeReport
// @Tags project-report(报表管理--项目报表)
// @Summary 报表管理--获取项目处理时长报表(日、周、月)
// @Auth xingqiyi
// @Date 2022/7/29 10:13
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param 	proCode     		query   string   true    "all：全部，B0118：B0118等"
// @Param 	type     			query   int   	 true    "1：日，2：周，3：月，4：年"
// @Param 	startTime      		query   string   true   "开始时间格式2022-07-06T16:00:00.000Z"
// @Param 	endTime        		query   string   true   "结束时间格式2022-07-06T16:00:00.000Z"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/project-report/deal-time-report [get]
func GetDealTimeReport(c *gin.Context) {
	var dealTimeReportReq request.BusinessReportReq
	err := c.ShouldBindQuery(&dealTimeReportReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	verify := utils.Rules{
		"Type": {utils.Gt("0")},
	}
	verifyErr := utils.Verify(dealTimeReportReq, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	err, obj := service.GetProReport(dealTimeReportReq, 3)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithoutCode(err.Error(), fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List: obj,
		}, c)
	}
}

// ExportReport
// @Tags project-report(报表管理--项目报表)
// @Summary 报表管理--导出报表
// @Auth xingqiyi
// @Date 2022/8/1 11:23
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param 	proCode     		query   string   true    "all：全部，B0118：B0118等"
// @Param 	type     			query   int   	 true    "1：日，2：周，3：月，4：年"
// @Param 	startTime      		query   string   true   "开始时间格式2022-07-06T16:00:00.000Z"
// @Param 	endTime        		query   string   true   "结束时间格式2022-07-06T16:00:00.000Z"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/project-report/report-export [get]
func ExportReport(c *gin.Context) {
	var reportReq request.BusinessReportReq
	err := c.ShouldBindQuery(&reportReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	verify := utils.Rules{
		"Type": {utils.Gt("0")},
	}
	verifyErr := utils.Verify(reportReq, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	err, businessObj := service.GetProReport(reportReq, 1)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, agingObj := service.GetProReport(reportReq, 2)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err, dealTimeObj := service.GetProReport(reportReq, 3)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	path := fmt.Sprintf("%v/%v/%v/", global.GConfig.LocalUpload.FilePath, global.PathProReportExport, reportReq.ProCode)
	name := "报表.xlsx"
	objArr := []map[string]interface{}{businessObj, agingObj, dealTimeObj}
	err = ExportReportExcel(path, name, objArr)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.FileAttachment(path+name, name)
}

func ExportReportExcel(path string, bookName string, objArr []map[string]interface{}) (err error) {
	//文件夹是否存在
	exists, err := utils.PathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		err = utils.CreateDir(path)
		if err != nil {
			return err
		}
	}

	file := excelize.NewFile()

	styleID, err := file.NewStyle(`{"font":{"color":"#000000"}}`)
	if err != nil {
		return err
	}
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}

	tableMap := []string{"业务量", "时效保障率", "处理时长"}
	rowId := 1
	isNewTable := true
	for tableIndex, obj := range objArr {
		for proCode, proData := range obj {
			if len(obj) < 1 {
				return errors.New("没有数据啊")
			}
			dataVal := make([]interface{}, 0)
			headVal := make([]interface{}, 0)
			headVal = append(headVal, "项目/"+tableMap[tableIndex])
			dataVal = append(dataVal, proCode)
			for _, proDayData := range proData.([]interface{}) {
				proDayStruct := proDayData.(map[string]interface{})
				//fmt.Println(proDayStruct["countDate"], proDayStruct["count"])
				headVal = append(headVal, proDayStruct["countDate"])
				dataVal = append(dataVal, proDayStruct["count"])
			}
			if isNewTable {
				headRow := make([]interface{}, len(headVal))
				for i := 0; i < len(headVal); i++ {
					headRow[i] = excelize.Cell{StyleID: styleID, Value: headVal[i]}
				}
				headCell, _ := excelize.CoordinatesToCellName(1, rowId)
				if err = streamWriter.SetRow(headCell, headRow); err != nil {
					return err
				}
				isNewTable = false
			}
			rowId++
			dataRow := make([]interface{}, len(dataVal))
			for i := 0; i < len(dataRow); i++ {
				dataRow[i] = excelize.Cell{StyleID: styleID, Value: dataVal[i]}
			}
			dataCell, _ := excelize.CoordinatesToCellName(1, rowId)
			if err = streamWriter.SetRow(dataCell, dataRow); err != nil {
				return err
			}
		}
		rowId = rowId + 5
		isNewTable = true
	}

	if err = streamWriter.Flush(); err != nil {
		return err
	}

	if err = file.SaveAs(path + bookName); err != nil {
		return err
	}
	return nil
}

// SetOtherReportInfo
// @Tags project-report(报表管理--项目日报)
// @Summary 报表管理--设置项目日报其他信息
// @Auth xingqiyi
// @Date 2022/8/5 10:38
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.ProReportOtherInfoReq true "项目日报其他信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/project-report/set-report-info [post]
func SetOtherReportInfo(c *gin.Context) {
	var proReportOtherInfo request2.ProReportOtherInfoReq
	err := c.ShouldBindJSON(&proReportOtherInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"ReportDate": {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(proReportOtherInfo, verify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}

	err = service.SetOtherReportInfo(proReportOtherInfo)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("创建失败，%v", err), c)
	} else {
		response.OkWithData("创建成功", c)
	}
}

// ReportInfoExport
// @Tags project-report(报表管理--项目日报)
// @Summary 报表管理--导出日报
// @Auth xingqiyi
// @Date 2022/8/5 17:43
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param reportDay       query   string   true    "项目日报日期格式2022-07-06T16:00:00.000Z"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/project-report/report-info-export [get]
func ReportInfoExport(c *gin.Context) {
	var proReportReq request2.ProReportReq
	err := c.ShouldBindQuery(&proReportReq)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list, otherInfo := service2.GetProReport(proReportReq)
	path := fmt.Sprintf("%v/%v/%v/", global.GConfig.LocalUpload.FilePath, global.PathProReportExport, "all")
	name := fmt.Sprintf("%v日报.xlsx", proReportReq.ReportDay.Format("2006-01-02"))
	err = exportProjectReportExcel(path, name, list, otherInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	c.FileAttachment(path+name, name)
}

// exportProjectReportExcel 导出日报
func exportProjectReportExcel(path, name string, list []model2.ProReport, otherInfo model2.ProReportOtherInfo) (err error) {

	//文件夹是否存在
	exists, err := pathExists(path)
	if err != nil {
		return err
	}
	if !exists {
		err = utils.CreateDir(path)
		if err != nil {
			return err
		}
	}

	if reflect.ValueOf(list).IsNil() || reflect.ValueOf(list).IsZero() {
		return errors.New("数据为空")
	}

	file := excelize.NewFile()
	//file.NewSheet("hello") 新建sheet
	streamWriter, err := file.NewStreamWriter("Sheet1")
	if err != nil {
		return err
	}
	styleID, err := file.NewStyle(`{"font":{"color":"#000000"}}`)
	if err != nil {
		return err
	}

	//写入数据
	switch reflect.TypeOf(list).Kind() {
	case reflect.Slice, reflect.Array:
		dataLen := reflect.ValueOf(list).Len()
		fmt.Println("数据长度", dataLen)
		for rowID := 0; rowID < dataLen; rowID++ {
			ele := reflect.ValueOf(list).Index(rowID)
			eleLens := ele.NumField()
			subLen := 0
			for i := 0; i < eleLens; i++ {
				//获取list的元素的第i个字段的tag的excel的值
				val := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
				if val == "" {
					subLen++
				}
			}
			//fmt.Println("1", subLen)
			eleLen := eleLens - subLen
			//fmt.Println("2", eleLen)
			//设置头
			headRow := make([]interface{}, eleLen)
			if rowID == 0 {
				fmt.Println("元素长度", eleLen)
				j := 0
				for i := 0; i < eleLens; i++ {
					//获取list的元素的第i个字段的tag的excel的值
					val := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
					if val == "" {
						continue
					}
					headRow[j] = excelize.Cell{StyleID: styleID, Value: val}
					j++
				}
				headCell, _ := excelize.CoordinatesToCellName(1, 1)
				if err = streamWriter.SetRow(headCell, headRow); err != nil {
					return err
				}
			}

			//写其他数据
			row := make([]interface{}, eleLen)
			k := 0
			for i := 0; i < eleLens; i++ {
				vals := reflect.TypeOf(list).Elem().Field(i).Tag.Get("excel")
				if vals == "" {
					continue
				}
				val := ele.Field(i)
				row[k] = val
				k++
			}
			cellData, _ := excelize.CoordinatesToCellName(1, rowID+2)
			if err = streamWriter.SetRow(cellData, row); err != nil {
				return err
			}
		}
	default:
		return errors.New("数据不是数组或切片")
	}

	otherInfoCell, _ := excelize.CoordinatesToCellName(1, len(list)+3)
	otherInfoRow := make([]interface{}, 7)
	valArr := []interface{}{"人员情况", "编制人数", otherInfo.UserCount, "实到人数", otherInfo.ActiveUserCount, "下班时间", otherInfo.ClosingTime}
	for i, _ := range otherInfoRow {
		otherInfoRow[i] = excelize.Cell{StyleID: styleID, Value: valArr[i]}
	}

	if err = streamWriter.SetRow(otherInfoCell, otherInfoRow); err != nil {
		return err
	}

	otherInfoCell, _ = excelize.CoordinatesToCellName(1, len(list)+4)
	otherInfoRow = make([]interface{}, 2)
	valArr = []interface{}{"其他运行情况", otherInfo.OtherMess}
	for i, _ := range otherInfoRow {
		otherInfoRow[i] = excelize.Cell{StyleID: styleID, Value: valArr[i]}
	}
	if err = streamWriter.SetRow(otherInfoCell, otherInfoRow); err != nil {
		return err
	}

	if err = streamWriter.Flush(); err != nil {
		return err
	}
	if err = file.SaveAs(path + name); err != nil {
		return err
	}
	return nil
}

// GetCharSum
// @Tags project-report(报表管理--项目字符统计表)
// @Summary 报表管理--获取字符统计
// @Auth xingqiyi
// @Date 2022/8/11 17:29
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex      query   int      true    "页码"
// @Param pageSize       query   int      true    "数量"
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间日期格式2022-07-06T16:00:00.000Z"
// @Param endTime query string true "结束时间日期格式2022-07-06T16:00:00.000Z"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/project-report/get-char-sum [get]
func GetCharSum(c *gin.Context) {
	var baseTimePageCode model.BaseTimePageCode
	err := c.ShouldBindQuery(&baseTimePageCode)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}

	err, list, row := service.GetCharSumReportByPage(baseTimePageCode)
	if err != nil {
		global.GLog.Error("", zap.Error(err))
		response.FailWithoutCode(err.Error(), fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:  list,
			Total: row,
		}, c)
	}
}

// ExportCharSum
// @Tags project-report(报表管理--项目字符统计表)
// @Summary 报表管理--导出字符统计
// @Auth xingqiyi
// @Date 2022/8/11 17:30
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param startTime query string true "开始时间日期格式2022-07-06T16:00:00.000Z"
// @Param endTime query string true "结束时间日期格式2022-07-06T16:00:00.000Z"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /report-management/project-report/export-char-sum [get]
func ExportCharSum(c *gin.Context) {
	var baseTimeRangeWithCode model.BaseTimeRangeWithCode
	err := c.ShouldBindQuery(&baseTimeRangeWithCode)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	err, list, _ := service.GetCharSumReportByPage(model.BaseTimePageCode{
		BaseTimeRangeWithCode: baseTimeRangeWithCode,
	})
	path := fmt.Sprintf("%v/%v/%v/", global.GConfig.LocalUpload.FilePath, global.PathProReportExport, baseTimeRangeWithCode.ProCode)
	//name := fmt.Sprintf("字符统计表%v-%v.xlsx", baseTimeRangeWithCode.StartTime.Format("20060102"), baseTimeRangeWithCode.EndTime.Format("20060102"))
	name := fmt.Sprintf("%v字符统计表.xlsx", baseTimeRangeWithCode.ProCode)
	err = utils.ExportBigExcel(path, name, "", list)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	global.GLog.Info("导出字符统计", zap.Any("path", path+name))
	c.FileAttachment(path+name, name)
}
