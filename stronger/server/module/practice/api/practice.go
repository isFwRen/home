package api

import (
	"errors"
	"fmt"
	"regexp"
	"server/global"
	"server/global/response"
	"server/module/practice/model"
	"server/module/practice/service"

	responseSysBase "server/module/sys_base/model/response"
	"server/module/task/request"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/wxnacy/wgo/arrays"
	"gorm.io/gorm"
)

// GetPracticeTask
// @Tags practice (练习录入)
// @Summary 领取练习分块
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param code        query   string   true    "工号"
// @Param name        query   string   true    "姓名"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"领取成功"}"
// @Router /practice/get/ [get]
func GetPracticeTask(c *gin.Context) {
	var taskGet model.TaskGet
	_ = c.ShouldBindQuery(&taskGet)
	// op := taskGet.Op
	code := taskGet.Code
	name := taskGet.Name

	err, block, user, pacticeSum := service.GetTaskBlock(global.GConfig.System.ProCode, code, name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("---------查询可领取新分块失败------------:", err)
		response.FailWithMessage(fmt.Sprintf("没有分块"), c)
		return
	} else if err != nil {
		fmt.Println("---------查询可领取异常------------:", err)
		response.FailWithMessage(fmt.Sprintf("领取异常"), c)
		return
	}
	err, data := getBillData(block, code)
	data.ApplyAt = user.ApplyAt
	data.FieldCharacter = pacticeSum.SummaryFieldCharacter
	data.AccuracyRate = pacticeSum.SummaryAccuracyRate

	// fmt.Println("---------blockblock------------:", block.Code)
	if arrays.Contains(user.Bcode, block.Code) == -1 {
		data.Videos = service.GetBlockVideo(block.Code)
	}

	if err != nil {
		fmt.Println("---------领取新分块失败------------:", err)
	} else {
		response.OkDetailed(data, "领取成功", c)
		return
	}

}

// SubmitTask
// @Tags practice (练习录入)
// @Summary 提交练习分块
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param bill        body   string   true    "单证"
// @Param block       body   string   true    "分块"
// @Param fields      body   string   true    "字段"
// @Param code        body   string   true    "工号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"提交成功"}"
// @Router /practice/submit [post]
func SubmitTask(c *gin.Context) {
	var taskSubmit model.TaskSubmit
	_ = c.ShouldBindJSON(&taskSubmit)
	fmt.Println("Bill:", taskSubmit.Bill)
	fmt.Println("Fields:", taskSubmit.Fields)
	err := service.SubmitTask(global.GConfig.System.ProCode, taskSubmit, taskSubmit.Code, taskSubmit.Block.ID, taskSubmit.Block.Code)
	if err != nil {
		fmt.Println("SubmitOpTask  err:", err)
		response.FailWithMessage(fmt.Sprintf("更新失败"), c)
	} else {
		// response.OkDetailed(responseOp.OpData{}, "更新成功", c)
		response.OkWithMessage("更新成功", c)
	}
}

// ExitPractice
// @Tags practice (练习录入)
// @Summary 退出练习
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param code        body   string   true    "工号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"提交成功"}"
// @Router /practice/exit [post]
func ExitPractice(c *gin.Context) {
	var taskGet request.TaskGet
	_ = c.ShouldBindJSON(&taskGet)
	err, msg := service.Exit_func(global.GConfig.System.ProCode, taskGet.Code)

	if err != nil {
		fmt.Println("ExitPractice  err:", err)
		response.FailWithMessage(fmt.Sprintf("退出练习失败"), c)
	} else {
		response.OkWithMessage(msg, c)
	}
}

func getBillData(block model.PracticeProjectBlock, code string) (error, model.OpData) {
	// if isUpdate {
	err, bill := service.SelectBillByID(global.GConfig.System.ProCode, block.BillID)
	if err != nil {
		return err, model.OpData{}
	}
	err, fields := service.SelectOpFieldsByBlockID(global.GConfig.System.ProCode, block.ID)
	opFields := service.FieldsFormat(fields)
	codeValues := make(map[string]interface{})
	// codeValues := project.BillValue(strings.Replace(global.GConfig.System.ProCode, "_task", "", -1), block, fields, op)
	return err, model.OpData{
		Bill:       bill,
		Block:      block,
		Fields:     opFields,
		CodeValues: codeValues,
		Videos:     []string{},
	}
}

// func UploadImage(c *gin.Context) {
// 	form, _ := c.MultipartForm()
// 	file := form.File["file"][0]
// 	path := form.Value["path"][0]
// 	name := form.Value["name"][0]
// 	// sysProTempId := c.PostForm("sysProTempId")
// 	// pathArr := make([]string, 0)
// 	// sysProTempBIdpath := ""
// 	// for _, file := range files {
// 	err, sysProTempBIdpath := utils.SaveImgFile(c, file, name, global.GConfig.LocalUpload.FilePath+path)
// 	fmt.Println("err:", err)
// 	fmt.Println("sysProTempBIdpath:", sysProTempBIdpath)
// 	if err != nil {
// 		response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
// 	}
// 	// pathArr = append(pathArr, sysProTempBIdpath)
// 	// }
// 	// reRow := service2.UpdateImages(pathArr, sysProTempId)
// 	// if reRow != 1 {
// 	// response.FailWithMessage(fmt.Sprintf("更新失败，%v", err), c)
// 	// } else {
// 	response.OkDetailed(name, "上传成功", c)
// 	// }
// }

func RegReplace(data string, query string, value string) string {
	reg := regexp.MustCompile(query)
	return reg.ReplaceAllString(data, value)
}

// GetPracticeSumByPage
// @Tags practice (练习管理)
// @Summary 练习管理--练习产量列表
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目编码"
// @Param code           query   string   false    "工号"
// @Param name       	 query   string   false    "姓名"
// @Param pageIndex      query   int      false    "页码"
// @Param pageSize       query   int      false    "数量"
// @Param startTime      query   string   true   "开始时间今天开始时间格式'2022-07-06T16:00:00.000Z'"
// @Param endTime        query   string   true   "结束时间今天结束时间格式'2022-07-06T16:00:00.000Z'"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"领取成功"}"
// @Router /practice/sum [get]
func GetPracticeSumByPage(c *gin.Context) {
	var sumSearch model.SumSearch
	fmt.Println("--------------GetPracticeSumByPage-----------")
	err := c.ShouldBindQuery(&sumSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	fmt.Println("--------------sumSearch-----------", sumSearch)
	err, total, bills := service.GetPracticeSumLists(sumSearch)
	for ii, bill := range bills {
		bills[ii].CostTime = time.Unix(bill.SummaryCostTime, 0).UTC().Format("15:04:05")
	}
	if err != nil {
		fmt.Println("--------------err-----------", err)
		// global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      bills,
			Total:     total,
			PageIndex: sumSearch.PageIndex,
			PageSize:  sumSearch.PageSize,
		}, c)
	}
}

// GetPracticeWrongByPage
// @Tags practice (练习管理)
// @Summary 练习管理--练习错误列表
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param proCode        query   string   true    "项目编码"
// @Param name       	 query   string   false    "字段名"
// @Param pageIndex      query   int      false    "页码"
// @Param pageSize       query   int      false    "数量"
// @Param startTime      query   string   true   "开始时间今天开始时间格式'2022-07-06T16:00:00.000Z'"
// @Param endTime        query   string   true   "结束时间今天结束时间格式'2022-07-06T16:00:00.000Z'"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"领取成功"}"
// @Router /practice/wrong [get]
func GetPracticeWrongByPage(c *gin.Context) {
	var sumSearch model.SumSearch
	fmt.Println("--------------GetPracticeWrongByPage-----------")
	err := c.ShouldBindQuery(&sumSearch)
	if err != nil {
		response.FailWithParamErr(err, c)
		return
	}
	fmt.Println("--------------GetPracticeWrongByPage-----------", sumSearch)
	err, total, bills := service.GetPracticeWrongLists(sumSearch)

	if err != nil {
		fmt.Println("--------------err-----------", err)
		// global.GLog.Error("", zap.Error(err))
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(responseSysBase.BasePageResult{
			List:      bills,
			Total:     total,
			PageIndex: sumSearch.PageIndex,
			PageSize:  sumSearch.PageSize,
		}, c)
	}
}

// GetPracticeAskList
// @Tags practice (练习管理)
// @Summary 练习管理--项目练习要求列表
// @Auth sf
// @Date 2023/10/7  下午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"领取成功"}"
// @Router /practice/ask-list [get]
func GetPracticeAskList(c *gin.Context) {
	err, list := service.PracticeAskList()
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败%v", err), c)
	} else {
		response.OkWithData(list, c)
	}
}
