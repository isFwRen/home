package api

import (
	"github.com/gin-gonic/gin"
	"reflect"
	"server/global"
	"server/global/response"
	reqModel "server/module/assessment_management/model/request"
	respModel "server/module/assessment_management/model/response"
	"server/module/assessment_management/service"
	trainStage "server/module/training_guide/service"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
	//"server/module/sys_base/model/request"
	// "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// StartExam
// @Tags assess-problem(考核管理-打单考核)
// @Summary 打单考核--开始考核 获取考核题目
// @Auth
// @Date 2023/6/1 10:30
// @Security XCode
// @Security ApiKeyAuth
// @Security XUserId
// @accept application/json
// @Produce application/json
// @Param data body reqModel.ReqStartExam true "考核开始时间及考核项目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /assessment-exam/test-procedure/start-exam  [post]
func StartExam(c *gin.Context) {
	var startExam reqModel.ReqStartExam
	err := c.ShouldBindJSON(&startExam)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	exam, err := service.StartExam(user.ID, startExam) //*****
	//exam, testList, err := service.StartExam(user.ID, startExam)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if reflect.DeepEqual(exam, respModel.RespStartExam{}) {
		response.OkWithMessage("没有题目可分配", c)
		return
	}
	response.OkWithData(exam, c) //正式代码，测试暂时移除
	//response.OkWithData(map[string]interface{}{
	//	"old": exam.AssessmentBlockList,
	//	"new": testList,
	//}, c)
	return

}

// EndExam
// @Tags assess-problem(考核管理-打单考核)
// @Summary 打单考核--结束考核 计算考核结果
// @Auth
// @Date 2023/6/1 10:30
// @Security XCode
// @Security ApiKeyAuth
// @Security XUserId
// @accept application/json
// @Produce application/json
// @Param data body reqModel.ReqEndExam true "考核结束时间及考核题目"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /assessment-exam/test-procedure/end-exam  [post]
func EndExam(c *gin.Context) {
	var endExam reqModel.ReqEndExam
	err := c.ShouldBindJSON(&endExam)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	user, err := api2.GetUserByToken(c)
	global.GLog.Info(">>>>>>>>>>>>>>>>>>>>>>>>>>" + user.Name)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//user01 := request.CustomClaims{
	//	ID:       "62e8bf716b254cd6a6cdb5be426d1aaf",
	//	NickName: "罗彩良",
	//	Code:     "7504",
	//}

	isPass, pointStr, err := service.EndExamScore(user, endExam)
	//isPass, pointStr, err := service.EndExamScore(&user01, endExam)

	//判断考核结果是否达到标准，并修改培训流程状态
	stage, err := trainStage.GetTrainingStage("", struct {
		UserID      string
		ProjectCode string
	}{UserID: user.ID, ProjectCode: global.GConfig.System.ProCode})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err = stage.NextStage(isPass, 4, global.GDb); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//----------------------------------------

	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(map[string]interface{}{
		"isPass": isPass,
		"point":  pointStr,
	}, c)
	return

}

// GetExamList
// @Tags assess-problem(考核管理-打单考核)
// @Summary 打单考核--获取项目 考核标准列表
// @Auth
// @Date 2023/6/1 10:30
// @Security XCode
// @Security ApiKeyAuth
// @Security XUserId
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /assessment-exam/test-procedure/get-project-list  [get]
func GetExamList(c *gin.Context) {

	user, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//user := &request.CustomClaims{
	//	ID: "6a61f4b179284b088cd35c01a314c299",
	//}

	list, err := service.GetAssessmentCriterion(user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(list) == 0 {
		list = []respModel.RespAssessCriterion{}
	}
	response.OkWithData(map[string]interface{}{
		"list":  list,
		"total": len(list),
	}, c)
	return
}
