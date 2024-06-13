/**
 * @Author: 星期一
 * @Description:
 * @Date: 2020/11/3 4:20 下午
 */

package model

import (
	"server/module/sys_base/model"
)

type SysProFieldCheck struct {
	model.Model
	CheckType string `json:"checkType" form:"checkType" gorm:"comment:'校验类型'"`
	Code      string `json:"code" form:"code" gorm:"comment:'字段编码'"`
	Value     string `json:"value" form:"value" gorm:"comment:'校验内容'"`
	Mark      string `json:"mark" form:"mark" gorm:"comment:'备注'"`
	ProId     string `json:"proId" form:"proId" gorm:"comment:项目id"`                //项目id
	FId       string `json:"fId" form:"fId" gorm:"comment:字段id" binding:"required"` //字段id
}

type SysFieldCheckEditReq struct {
	List []SysProFieldCheck `json:"list" form:"list"`
	FId  string             `json:"fId" form:"fId" binding:"required"` //字段id
}
