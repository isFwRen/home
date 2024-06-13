package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"server/global/response"
	sys_base3 "server/module/sys_base/model/request"
	sys_base4 "server/module/sys_base/model/response"
	"server/module/sys_base/service"
	"server/utils"
)

// Register
// @Tags SysBase
// @Summary 用户注册账号
// @Produce  application/json
// @Param data body sys_base3.RegisterStruct true "用户注册接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"用户注册成功"}"
// @Router /sys-base/sys-register/register [post]
func Register(c *gin.Context) {
	form, _ := c.MultipartForm()
	var R sys_base3.RegisterStruct
	_ = c.ShouldBindJSON(&R)
	text := form.Value["body"]
	text1 := []byte(text[0])
	json.Unmarshal(text1, &R)
	UserVerify := utils.Rules{
		"IDCard": {utils.NotEmpty()},
		"BankId": {utils.NotEmpty()},
	}
	UserVerifyErr := utils.Verify(R, UserVerify)
	if UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}
	err, userReturn := service.Register(R)
	if err != nil {
		response.FailWithDetailed(response.ERROR, sys_base4.RegisterResponse{RegisterUser: userReturn}, fmt.Sprintf("%v", err), c)
	} else {
		response.Ok(c)
	}
}

// Resignation
// @Tags SysUser
// @Summary 用户提交离职申请
// @Produce  application/json
// @Param data body sys_base3.ResignationStruct true "提交离职申请接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"用户提交离职申请成功"}"
// @Router /sys-user-leave/user-leave/resignation [post]
func Resignation(c *gin.Context) {
	var R sys_base3.ResignationStruct
	_ = c.ShouldBindJSON(&R)
	fmt.Println(R)
	err, _ := service.Resignation(R)
	fmt.Println(err)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%v", err), c)
	} else {
		response.OkWithMessage("已成功提交离职申请", c)
	}
}
