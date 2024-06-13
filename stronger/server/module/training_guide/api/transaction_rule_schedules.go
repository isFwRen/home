package api

import (
	"github.com/gin-gonic/gin"
	"server/global"
	"server/global/response"
	"server/module/training_guide/request"
	"server/module/training_guide/service"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
)

// FinishStep
// @Tags  training-guide(培训流程-培训指引)
// @Summary 项目规则--完成一个文件学习
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.FinishParam true "完成一个文件学习"
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /training-guide/finish-read [post]
func FinishStep(c *gin.Context) {
	var params request.FinishParam
	err := c.ShouldBindJSON(&params)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	user, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = service.FinishStep(params, user)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	isFinish, err := service.CheckFinishAll(params.ProjectCode, user)
	if isFinish {
		stage, err := service.GetTrainingStage("", struct {
			UserID      string
			ProjectCode string
		}{UserID: user.ID, ProjectCode: params.ProjectCode})
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		if err = stage.NextStage(isFinish, 2, global.GDb); err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
	}
	response.OkWithMessage("已完成该文件学习", c)
	return
}

// FinishDocReading
// @Tags  training-guide(培训流程-培训指引)
// @Summary 项目规则--完成培训流程指引文阅读
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string	"{"success":200,"data":{},"msg":"获取成功"}"
// @Router /training-guide/finish-read-doc [get]
func FinishDocReading(c *gin.Context) {
	user, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	stage, err := service.GetTrainingStage("", struct {
		UserID      string
		ProjectCode string
	}{UserID: user.ID, ProjectCode: global.GConfig.System.ProCode})
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err = stage.NextStage(true, 1, global.GDb); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("已完成该文件学习", c)
	return
}
