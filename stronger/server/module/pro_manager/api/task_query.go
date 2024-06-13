package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
	"net/http"
	"server/global"
	"server/global/response"
	m "server/module/pro_manager/model"
	res "server/module/pro_manager/model/request"
	resp "server/module/pro_manager/model/response"
	server "server/module/pro_manager/service"
	"server/utils"
)

// GetTaskListDetail
// @Tags Task List (任务管理)
// @Summary 任务管理--查询待分配,已分配,缓存区,紧急件,优先件的详细信息
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param op query string true "工序"
// @Param opStage query string true "工序阶段"
// @Param isExpenseAccount query string false "报销单 1 or 非报销单 2 or 初审问题件不要传"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-config/task/detail/List [get]
func GetTaskListDetail(c *gin.Context) {
	var cols res.GetTaskListDetail
	if err := c.ShouldBindWith(&cols, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetTaskListVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
		"Op":      {utils.NotEmpty()},
		"OpStage": {utils.NotEmpty()},
	}
	GetTaskListVerifyErr := utils.Verify(cols, GetTaskListVerify)
	if GetTaskListVerifyErr != nil {
		response.FailWithMessage(GetTaskListVerifyErr.Error(), c)
		return
	}
	err, taskList, total := server.GetTaskListDetail(cols)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("获取待分配-已分配-缓存区-紧急件-优先件 详细信息失败, err : %v", err), c)
	} else {
		response.GetOkWithData(resp.GetTaskList{List: taskList, Total: total}, c)
	}
}

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// GetTaskList
// @Tags Task List (任务管理)
// @Summary 任务管理--刷新
// @accept application/json
// @Produce application/json
// @Param proCode query string true "查询待分配-已分配-缓存区-紧急件-优先件接口"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-config/task/List [get]
func GetTaskList(c *gin.Context) {
	//查询待分配-已分配-缓存区-紧急件-优先件-单子数量
	var cols res.GetTaskList
	if err := c.ShouldBindWith(&cols, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetTaskListVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	GetTaskListVerifyErr := utils.Verify(cols, GetTaskListVerify)
	if GetTaskListVerifyErr != nil {
		response.FailWithMessage(GetTaskListVerifyErr.Error(), c)
		return
	}
	err, taskList, _ := server.GetTaskList(cols)
	if err != nil {
		if err == global.ProDbErr {
			var taskLists m.TaskList
			c.JSON(http.StatusOK, Response{
				200,
				taskLists,
				err.Error(),
			})
			return
		}
		response.GetFailWithData(fmt.Sprintf("刷新失败，%v", err), c)
	} else {
		response.GetOkWithData(taskList, c)
	}
}

// GetUrgencyBillOrPriorityBill
// @Tags Task List (任务管理)
// @Summary 任务管理--查询紧急件,优先件
// @accept application/json
// @Produce application/json
// @Param data body res.GetVariousStateBill true "查询紧急件/优先件接口"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-config/task/UrgencyBillOrPriorityBill/list [post]
func GetUrgencyBillOrPriorityBill(c *gin.Context) {
	var cols res.GetVariousStateBill
	_ = c.ShouldBindJSON(&cols)
	GetUrgencyBillOrPriorityBillVerify := utils.Rules{
		"ProCode":    {utils.NotEmpty()},
		"StickLevel": {utils.Gt("-1"), utils.Lt("3")},
	}
	GetUrgencyBillOrPriorityBillVerifyErr := utils.Verify(cols, GetUrgencyBillOrPriorityBillVerify)
	if GetUrgencyBillOrPriorityBillVerifyErr != nil {
		response.FailWithMessage(GetUrgencyBillOrPriorityBillVerifyErr.Error(), c)
		return
	}
	err, UrgencyBillOrPriorityBillNumber, total := server.GetUrgencyBillOrPriorityBill(cols)
	if err != nil {
		response.GetFailWithData(fmt.Sprintf("获取优先/紧急单失败，%v", err), c)
	} else {
		response.GetOkWithData(resp.GetVariousStateBillResponse{List: UrgencyBillOrPriorityBillNumber, Total: total}, c)
	}
}

// SetUrgencyBillOrPriorityBill
// @Tags Task List (任务管理)
// @Summary 任务管理--设置紧急件,优先件
// @accept application/json
// @Produce application/json
// @Param data body res.SetVariousStateBill true "设置紧急件/优先件接口"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"设置优先/紧急单成功"}"
// @Router /pro-config/task/UrgencyBillOrPriorityBill/edit [post]
func SetUrgencyBillOrPriorityBill(c *gin.Context) {
	var cols res.SetVariousStateBill
	_ = c.ShouldBindJSON(&cols)
	SetUrgencyBillOrPriorityBillVerify := utils.Rules{
		"ProCode":     {utils.NotEmpty()},
		"CaseNumbers": {utils.NotEmpty()},
	}
	SetUrgencyBillOrPriorityBillVerifyErr := utils.Verify(cols, SetUrgencyBillOrPriorityBillVerify)
	if SetUrgencyBillOrPriorityBillVerifyErr != nil {
		response.FailWithMessage(SetUrgencyBillOrPriorityBillVerifyErr.Error(), c)
		return
	}
	err := server.SetUrgencyBillOrPriorityBill(cols)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("设置优先/紧急单失败, err : %v", err), c)
	} else {
		response.OkWithData(fmt.Sprintf("设置优先/紧急单成功,%v", err), c)
	}
}

// SetPriorityOrganizationNumber
// @Tags Task List (任务管理)
// @Summary 任务管理--设置优先处理的机构号
// @accept application/json
// @Produce application/json
// @Param data body res.SetPriorityOrganizationNumber true "设置优先处理的机构号接口"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"设置成功"}"
// @Router /pro-config/task/PriorityOrganizationNumber/edit [post]
func SetPriorityOrganizationNumber(c *gin.Context) {
	var cols res.SetPriorityOrganizationNumber
	_ = c.ShouldBindJSON(&cols)
	SetPriorityOrganizationNumberVerify := utils.Rules{
		"ProCode":            {utils.NotEmpty()},
		"StickLevel":         {utils.Gt("-1"), utils.Lt("100")},
		"OrganizationNumber": {utils.NotEmpty()},
	}
	SetPriorityOrganizationNumberVerifyErr := utils.Verify(cols, SetPriorityOrganizationNumberVerify)
	if SetPriorityOrganizationNumberVerifyErr != nil {
		response.FailWithMessage(SetPriorityOrganizationNumberVerifyErr.Error(), c)
		return
	}
	err := server.SetPriorityOrganizationNumber(cols)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("设置优先处理的机构号失败, err : %v", err), c)
	} else {
		response.OkWithData(fmt.Sprintf("设置优先处理的机构号成功,%v", err), c)
	}
}

// GetCaseDetails
// @Tags Case Details (案件明细)
// @Summary 任务管理--案件明细查询
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-config/CaseDetails/list [get]
func GetCaseDetails(c *gin.Context) {
	//获取案件明细, 主要查询在录入中的单子
	var cols res.GetCaseDetails
	if err := c.ShouldBindWith(&cols, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetCaseDetailsVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	GetCaseDetailsVerifyErr := utils.Verify(cols, GetCaseDetailsVerify)
	if GetCaseDetailsVerifyErr != nil {
		response.FailWithMessage(GetCaseDetailsVerifyErr.Error(), c)
		return
	}
	err, result, total := server.GetCaseDetails(cols)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("查询失败, %v", err), c)
	} else {
		response.OkWithData(resp.GetCaseDetailsResponse{
			List:  result,
			Total: total,
		}, c)
	}
}

// GetCaseDetailsBlock
// @Tags Case Details (案件明细)
// @Summary 任务管理--分块明细查询
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param id query string true "单号ID"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-config/CaseDetails/block/list [get]
func GetCaseDetailsBlock(c *gin.Context) {
	//获取案件明细-分块明细, 主要查询在录入中的单子
	var cols res.GetCaseDetailsBlock
	if err := c.ShouldBindWith(&cols, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetCaseDetailsBlockVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	GetCaseDetailsBlockVerifyErr := utils.Verify(cols, GetCaseDetailsBlockVerify)
	if GetCaseDetailsBlockVerifyErr != nil {
		response.FailWithMessage(GetCaseDetailsBlockVerifyErr.Error(), c)
		return
	}
	err, result := server.GetCaseDetailsBlock(cols)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取案件明细-分块明细失败, %v", err), c)
	} else {
		response.OkWithData(resp.GetCaseDetailsResponse{
			List: result,
		}, c)
	}
}

// GetCaseDetailsField
// @Tags Case Details (案件明细)
// @Summary 任务管理--字段明细查询
// @accept application/json
// @Produce application/json
// @Param proCode query string true "项目编码"
// @Param billId query string true "单号ID"
// @Param blockId query string false "分块ID"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /pro-config/CaseDetails/field/list [get]
func GetCaseDetailsField(c *gin.Context) {
	//获取案件明细-字段明细, 主要查询在录入中的单子
	var cols res.GetCaseDetailsField
	if err := c.ShouldBindWith(&cols, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetCaseDetailsFieldVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
	}
	GetCaseDetailsFieldVerifyErr := utils.Verify(cols, GetCaseDetailsFieldVerify)
	if GetCaseDetailsFieldVerifyErr != nil {
		response.FailWithMessage(GetCaseDetailsFieldVerifyErr.Error(), c)
		return
	}
	err, result := server.GetCaseDetailsField(cols)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取案件明细-字段明细失败, %v", err), c)
	} else {
		response.OkWithData(resp.GetCaseDetailsResponse{
			List: result,
		}, c)
	}
}

// SetUrgencyPriorityBill
// @Tags Case Details (案件明细)
// @Summary 任务查询--将单据紧急1,优先2,正常0状态
// @Auth xingqiyi
// @Date 2022/2/23 11:26 上午
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body res.SetVariousStateBill true "紧急单结构体"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/task/urgency-priority-bill/edit [post]
func SetUrgencyPriorityBill(c *gin.Context) {
	var cols res.SetVariousStateBill
	_ = c.ShouldBindJSON(&cols)
	paramVerify := utils.Rules{
		"ProCode":     {utils.NotEmpty()},
		"StickLevel":  {utils.Gt("-1"), utils.Lt("3")},
		"CaseNumbers": {utils.NotEmpty()},
	}
	verifyErr := utils.Verify(cols, paramVerify)
	if verifyErr != nil {
		response.FailWithMessage(verifyErr.Error(), c)
		return
	}
	err := server.SetUrgencyPriorityBill(cols, "")
	if err != nil {
		global.GLog.Error("更新历史库失败", zap.Error(err))
	}
	err = server.SetUrgencyPriorityBill(cols, "_task")
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("设置优先/紧急单失败, err : %v", err), c)
	} else {
		response.OkWithData("设置优先/紧急单成功", c)
	}
}

// GetProPermissionPeople
// @Tags Case Details (任务管理)
// @Summary 任务查询--获取拥有该项目各个工序的用户
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param pageIndex query string true "页码"
// @Param pageSize query string true "数量"
// @Param proCode query string true "项目编码"
// @Param op query string true "工序"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pro-config/task/getProPermissionPeople [get]
func GetProPermissionPeople(c *gin.Context) {
	var cols res.Search
	if err := c.ShouldBindWith(&cols, binding.Query); err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败，%v", err), c)
		return
	}
	GetTaskListVerify := utils.Rules{
		"ProCode": {utils.NotEmpty()},
		"op":      {utils.NotEmpty()},
	}
	GetTaskListVerifyErr := utils.Verify(cols, GetTaskListVerify)
	if GetTaskListVerifyErr != nil {
		response.FailWithMessage(GetTaskListVerifyErr.Error(), c)
		return
	}
	err, list, total := server.GetProPermissionPeople(cols)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("获取失败"), c)
	} else {
		response.OkWithData(resp.GetCaseDetailsResponse{
			List:  list,
			Total: total,
		}, c)
	}
}
