/**
 * @Author: xingqiyi
 * @Description: 问题件配置
 * @Date: 2023/2/1 11:42
 */

package model

import "server/module/sys_base/model"

type SysIssue struct {
	model.Model
	FId       string `json:"fId" form:"fId"  gorm:"comment:字段id" binding:"required"` //字段id
	InputVal  string `json:"inputVal" form:"inputVal" gorm:"comment:录入值"`            //录入值
	ChangeVal string `json:"changeVal" form:"changeVal" gorm:"comment:转换值"`          //转换值
	Code      string `json:"code" form:"code" gorm:"comment:问题件编码"`                  //问题件编码
	Desc      string `json:"desc" form:"desc" gorm:"comment:问题件描述"`                  //问题件描述
	IssueType string `json:"issueType" form:"issueType" gorm:"comment:问题件类型"`        //问题件类型
	ProId     string `json:"proId" form:"proId" gorm:"comment:项目id"`                 //项目id
}

type SysIssueEditReq struct {
	List []SysIssue `json:"list" form:"list"`
	FId  string     `json:"fId" form:"fId" binding:"required"` //字段id
}
