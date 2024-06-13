/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/6 15:22
 */

package model

import modelBase "server/module/sys_base/model"

type TaskNotice struct {
	modelBase.Model
	ProCode    string `json:"proCode"`
	Msg        string `json:"msg"`
	SendCode   string `json:"sendCode"`
	SendName   string `json:"sendName"`
	SendStatus int    `json:"sendStatus"`
	FailReason string `json:"failReason"`
}
