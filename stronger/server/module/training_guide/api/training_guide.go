package api

import (
	"github.com/gin-gonic/gin"
	"server/global/response"
	api2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"

	//sys_base2 "server/module/sys_base/model/request"
	//"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
	"server/module/training_guide/service"
	//"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/api"
	//"server/module/sys_base/model/request"
	//base "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
)

// AddAndFindGuide
// @Tags training-guide(培训流程-培训指引)
// @Summary 培训指引-获取培训阶段
// @Auth xingqiyi
// @Date 2021/11/3 3:54 下午
// @Security XCode
// @Security ApiKeyAuth
// @Security XUserId
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /training-guide/training-stage/find [post]
func AddAndFindGuide(c *gin.Context) {

	user, err := api2.GetUserByToken(c)
	//user1 := user.(request.CustomClaims)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	respSysUser, err := service.AddAndFindTrainingStage(user.ID, user.NickName)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(map[string]int{
		"userTrainingStage": respSysUser.TrainingStage,
	}, c)
	return

}

// UpdateTrainingStage
// @Tags training-guide(培训流程-培训指引)
// @Summary 培训指引-更新培训阶段
// @Auth xingqiyi
// @Date 2021/11/3 3:54 下午
// @Security XCode
// @Security ApiKeyAuth
// @Security XUserId
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"成功"}"
// @Router /training-guide/training-stage/update [put]
func UpdateTrainingStage(c *gin.Context) {
	user, err := api2.GetUserByToken(c)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	//userID:=
	if err := service.UpdateTrainingStage(user.ID, 1); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新进度成功", c)
}
