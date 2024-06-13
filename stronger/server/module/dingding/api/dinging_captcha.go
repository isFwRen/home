package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"server/global"
	"server/global/response"
	res "server/module/dingding/model/request"
	"server/module/dingding/service"
	sys_base4 "server/module/sys_base/model/response"
	"server/utils"
)

// Captcha
// @Tags base
// @Summary 生成验证码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body res.GetCaptcha true "生成验证码接口"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /dinging/captcha [post]
func Captcha(c *gin.Context) {
	//字符,公式,验证码配置
	//生成默认数字的driver
	var GetCaptcha res.GetCaptcha
	_ = c.ShouldBindJSON(&GetCaptcha)
	UserVerify := utils.Rules{
		"Phone": {utils.NotEmpty()},
	}
	UserVerifyErr := utils.Verify(GetCaptcha, UserVerify)
	if UserVerifyErr != nil {
		response.FailWithMessage(UserVerifyErr.Error(), c)
		return
	}

	driver := base64Captcha.NewDriverDigit(global.GConfig.Captcha.ImgHeight, global.GConfig.Captcha.ImgWidth, global.GConfig.Captcha.KeyLong, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, global.GStore)
	fmt.Println(driver)
	id, _, err := cp.Generate()
	fmt.Println("验证码", id, cp.Store.Get(id, false))
	fmt.Println("验证码", global.GStore.Get(id, false))
	//if err != nil {
	//	response.FailWithMessage(fmt.Sprintf("获取数据失败，%v", err), c)
	//} else {
	//	_, FailMsg := service.DdingTalk(GetCaptcha.Phone, cp.Store.Get(id, false))
	//	if FailMsg == "" {
	//		response.OkDetailed(sys_base4.SysCaptchaResponse{
	//			CaptchaId: id,
	//			Captcha:   cp.Store.Get(id, false),
	//			PicPath:   b64s,
	//		}, "验证码获取成功", c)
	//	} else {
	//		response.FailWithMessage(fmt.Sprintf("发送钉钉验证码失败，%s", FailMsg), c)
	//	}
	//}
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("生成验证码失败，%v", err), c)
	} else {
		CaptchaCode := global.GStore.Get(id, false)
		err = utils.SetRedisCaptcha(CaptchaCode, GetCaptcha.Phone)
		if err != nil {
			response.FailWithMessage(fmt.Sprintf("设置验证码过期时间失败，%v", err), c)
		} else {
			_, FailMsg := service.DdingTalk(GetCaptcha.Phone, cp.Store.Get(id, false))
			if FailMsg == "" {
				response.OkDetailed(sys_base4.SysCaptchaResponse{
					//CaptchaId: id,
					//Captcha: CaptchaCode,
					//PicPath:   b64s,
				}, "验证码获取成功", c)
			} else {
				response.FailWithMessage(fmt.Sprintf("发送钉钉验证码失败，%s", FailMsg), c)
			}
		}

	}

}
