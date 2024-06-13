/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/3 1:51 下午
 */

package model

import "server/module/sys_base/model"

type ResultDataLog struct {
	model.Model
	Name      string `json:"name"`      //名称
	BillId    string `json:"billId"`    //单据id
	BeforeVal string `json:"beforeVal"` //修改前的内容
	AfterVal  string `json:"afterVal"`  //修改后的内容
	EditCode  string `json:"editCode"`  //编辑的工号
	EditName  string `json:"editName"`  //编辑的名字
}
