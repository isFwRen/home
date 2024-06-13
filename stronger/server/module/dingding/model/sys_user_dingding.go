/**
 * @Author: 杨沛基
 * @Description:
 * @Date: 2020/11/9 17:30
 */

package model

import (
	"server/module/sys_base/model"
)

type DingAccessTaken struct {
	model.Model
	ErrCode     string `json:"err_code"`
	AccessToken string `json:"access_token"`
	ErrMsg      string `json:"err_msg"`
	ExpiresIn   string `json:"expires_in"`
}

type DingUserMessage struct {
	Mobile string `json:"mobile"`
}


