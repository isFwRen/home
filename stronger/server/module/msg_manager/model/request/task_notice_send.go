/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/6 14:57
 */

package model

type TaskNoticeSendReq struct {
	ProCode []string `json:"proCode"`
	Msg     string   `json:"msg"`
}
