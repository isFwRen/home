package api

import (
	"fmt"
	"regexp"
	"server/global"
	"server/global/response"
	sys_base2 "server/module/sys_base/model"
	sys_base3 "server/module/sys_base/model/request"
	"server/module/sys_base/service"
	"server/utils"

	"github.com/gin-gonic/gin"
)

// ChangePassword
// @Tags SysUser
// @Summary 用户更改密码
// @Produce  application/json
// @Param data body sys_base3.ChangePasswordStruct true "更改密码接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"用户更改密码成功"}"
// @Router /sys-user/sys-user-operation-pw/changePassword [post]
func ChangePassword(c *gin.Context) {
	var params sys_base3.ChangePasswordStruct
	_ = c.ShouldBindJSON(&params)
	compile, _ := regexp.Compile("[0-9]+")
	compile1, _ := regexp.Compile("[a-zA-Z]+")
	compile2, _ := regexp.Compile("[!@#$%^&*()`~.,;'！￥。，；‘——_|+-/《》、]+")
	if !params.IsIntranet {
		//外网
		UserVerify := utils.Rules{
			"Phone": {utils.NotEmpty()},
			// "Captcha":     {utils.NotEmpty()},
			"Password":    {utils.NotEmpty()},
			"NewPassword": {utils.NotEmpty()},
		}
		UserVerifyErr := utils.Verify(params, UserVerify)
		if UserVerifyErr != nil {
			response.FailWithMessage(UserVerifyErr.Error(), c)
			return
		}
		if !(compile.MatchString(params.NewPassword) && compile1.MatchString(params.NewPassword) && compile2.MatchString(params.NewPassword) && 16 > len(params.NewPassword) && len(params.NewPassword) > 8) {
			response.FailWithMessage("新密码不符合密码为数字，字母和符号的组合，长度为9-16位要求的密码", c)
			return
		}
		U := &sys_base2.SysUser{Phone: params.Phone, Password: params.Password}

		// err, a := utils.GetRedisCaptcha(params.Phone)
		// if err == nil && a == "" {
		// 	response.FailWithMessage("验证码已过期!", c)
		// 	return
		// }

		// if a == params.Captcha {
		// fmt.Println("global.GStore.Verify", global.GStore.Verify(params.CaptchaId, params.Captcha, true))
		if err, _ := service.ChangePassword(U, params.Password, params.NewPassword, params.IsIntranet); err != nil {
			response.FailWithMessage(fmt.Sprintf("修改失败, 请确认工号、原密码是否正确"), c)
		} else {
			err = utils.DelRedisCaptcha(params.Phone)
			if err != nil {
				response.FailWithMessage("验证码缓存删除失败!", c)
			}
			response.OkWithMessage("修改成功", c)
		}
		// } else {
		// 	response.FailWithMessage("验证码错误", c)
		// }
	} else {
		UserVerify := utils.Rules{
			"Code":        {utils.NotEmpty()},
			"Password":    {utils.NotEmpty()},
			"NewPassword": {utils.NotEmpty()},
		}
		UserVerifyErr := utils.Verify(params, UserVerify)
		if UserVerifyErr != nil {
			response.FailWithMessage(UserVerifyErr.Error(), c)
			return
		}
		if !(compile.MatchString(params.NewPassword) && compile1.MatchString(params.NewPassword) && compile2.MatchString(params.NewPassword) && 16 > len(params.NewPassword) && len(params.NewPassword) > 8) {
			response.FailWithMessage("新密码不符合密码为数字，字母和符号的组合，长度为9-16位要求的密码", c)
			return
		}
		U := &sys_base2.SysUser{Code: params.Username, Password: params.Password}
		if err, _ := service.ChangePassword(U, params.Password, params.NewPassword, params.IsIntranet); err != nil {
			response.FailWithMessage(fmt.Sprintf("修改失败, 请确认工号、原密码是否正确"), c)
		} else {
			response.OkWithMessage("修改成功", c)
		}
	}
}

// ResetPassword
// @Tags SysUser
// @Summary 用户重置密码
// @Produce  application/json
// @Param data body sys_base3.ResetPassword true "重置密码接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"用户重置密码成功"}"
// @Router /sys-base/sys-user-operation-pw/resetPassword [post]
func ResetPassword(c *gin.Context) {
	//判断内外网登录, 内网传true, 外网传false
	var resetPassword sys_base3.ResetPassword
	_ = c.ShouldBindJSON(&resetPassword)
	if !resetPassword.IsIntranet {
		UserVerify := utils.Rules{
			"Phone": {utils.NotEmpty()},
		}
		UserVerifyErr := utils.Verify(resetPassword, UserVerify)
		if UserVerifyErr != nil {
			response.FailWithMessage(UserVerifyErr.Error(), c)
			return
		}
	} else {
		UserVerify := utils.Rules{
			"Username": {utils.NotEmpty()},
		}
		UserVerifyErr := utils.Verify(resetPassword, UserVerify)
		if UserVerifyErr != nil {
			response.FailWithMessage(UserVerifyErr.Error(), c)
			return
		}
	}
	err, a := utils.GetRedisCaptcha(resetPassword.Phone)
	if err == nil && a == "" {
		response.FailWithMessage("验证码已过期!", c)
		return
	}
	if a == resetPassword.Captcha {
		fmt.Println("global.GStore:", global.GStore.Verify(resetPassword.CaptchaId, resetPassword.Captcha, true))
		U := &sys_base2.SysUser{Code: resetPassword.Username, Phone: resetPassword.Phone}
		if err, _ := service.ResetPassword(U, resetPassword.IsIntranet); err != nil {
			response.FailWithMessage("重置密码失败，请检查工号", c)
		} else {
			response.OkWithMessage("重置密码成功,初始密码为123456", c)
		}
	} else {
		response.FailWithMessage("验证码错误", c)
	}
}
