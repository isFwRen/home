/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/6 15:22
 */

package model

import modelBase "server/module/sys_base/model"

type DingtalkNotice struct {
	modelBase.Model
	GroupName  string `json:"groupName"`
	GroupId    string `json:"groupId"`
	Msg        string `json:"msg"`
	SendCode   string `json:"sendCode"`
	SendName   string `json:"sendName"`
	SendStatus int    `json:"sendStatus"`
	FailReason string `json:"failReason"`
}
