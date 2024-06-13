/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2023/4/26 14:54
 */

package model

import "server/module/sys_base/model"

type InvoiceNum struct {
	model.Model
	Num     string `json:"num" form:"num"`
	BillNum string `json:"billNum" form:"billNum"`
}
