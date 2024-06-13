package request

import (
	modelBase "server/module/sys_base/model"
)

type UserManagement struct {
	PageInfo  modelBase.BasePageInfo `json:"page_info"`
	Code      string                 `json:"code" form:"code"`           //工号
	Name      string                 `json:"name" form:"name"`           //姓名
	Role      string                 `json:"role" form:"role"`           //角色
	Phone     string                 `json:"phone" form:"phone"`         //手机号
	Status    string                 `json:"status" form:"status"`       //状态
	StartTime string                 `json:"startTime" form:"startTime"` //上岗时间start
	EndTime   string                 `json:"endTime" form:"endTime"`     //上岗时间end
}

type UserAdd struct {
	Code             string `json:"code"`             //工号
	NickName         string `json:"nickName"`         //姓名
	Staff            string `json:"staff"`            //角色
	Phone            string `json:"phone"`            //手机号
	IsMobileTerminal bool   `json:"isMobileTerminal"` //是否手机端
	Referees         string `json:"referees"`         //推荐人
	EntryDate        string `json:"entryDate"`        //入职日期
	MountGuardDate   string `json:"mountGuardDate"`   //上岗日期
	LeaveDate        string `json:"leaveDate"`        //离职日期
}

type ProPermissionInUserInformation struct {
	SysProPermission []modelBase.SysProPermission `json:"sysProPermission"`
}

type ProPermission struct {
	ProCode   string `json:"proCode"`   //项目代码
	ProName   string `json:"proName"`   //项目名字
	HasOp0    bool   `json:"hasOp0"`    //初审
	HasOp1    bool   `json:"hasOp1"`    //一码
	HasOp2    bool   `json:"hasOp2"`    //二码
	HasOpq    bool   `json:"hasOpq"`    //问题件
	HasInNet  bool   `json:"hasInNet"`  //内网
	HasOutNet bool   `json:"hasOutNet"` //外网
	UserCode  string `json:"userCode"`  //用户工号
	HasPm     bool   `json:"hasPm"`     //管理员
}

type GetPermissionInformation struct {
	UserId string `json:"userId" form:"userId"`
}
