/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/18 3:38 下午
 */

package model

import (
	"server/module/load/model"
	model2 "server/module/pro_conf/model"
)

type BillObj struct {
	ProjectBill      ProjectBill          `json:"projectBill"`
	ProjectBlockList []model.ProjectBlock `json:"projectBlockList"`
	ProjectFieldList []model.ProjectField `json:"projectFieldList"`
}

type FieldObj struct {
	Field     model.ProjectField `json:"field"`
	Block     model.ProjectBlock `json:"block"`
	FieldConf model2.SysProField `json:"fieldConf"`
}
