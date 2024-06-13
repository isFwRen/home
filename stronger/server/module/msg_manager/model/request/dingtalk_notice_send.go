/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/6 14:57
 */

package model

import "server/module/sys_base/model"

type DingtalkNoticeSendReq struct {
	GroupId []string `json:"groupId"`
	Msg     string   `json:"msg"`
}

type DingtalkNoticeMsgReq struct {
	model.BasePageInfo
}
