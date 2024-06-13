/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/31 10:54 上午
 */

package model

type SysProPermission struct {
	Model
	ProCode   string `json:"proCode"`   //项目代码
	ProId     string `json:"proId"`     //项目id
	ProName   string `json:"proName"`   //项目名字
	HasOp0    bool   `json:"hasOp0"`    //初审
	HasOp1    bool   `json:"hasOp1"`    //一码
	HasOp2    bool   `json:"hasOp2"`    //二码
	HasOpq    bool   `json:"hasOpq"`    //问题件
	HasInNet  bool   `json:"hasInNet"`  //内网
	HasOutNet bool   `json:"hasOutNet"` //外网
	UserCode  string `json:"userCode"`  //用户工号
	UserId    string `json:"userId"`    //用户id
	ObjectId  string `json:"objectId"`  //mongodb的useid
	HasPm     bool   `json:"hasPm"`     //管理该项目的权限
}

type SysProPermissionExport struct {
	ProCode   string `json:"proCode" excel:"项目编码"` //项目代码
	UserCode  string `json:"userCode" excel:"工号"`  //用户工号
	UserName  string `json:"userName" excel:"姓名"`  //用户姓名
	HasOp0    bool   `json:"hasOp0" excel:"初审"`    //初审
	HasOp1    bool   `json:"hasOp1" excel:"一码"`    //一码
	HasOp2    bool   `json:"hasOp2" excel:"二码"`    //二码
	HasOpq    bool   `json:"hasOpq" excel:"问题件"`   //问题件
	HasPm     bool   `json:"hasPm" excel:"PM"`     //管理该项目的权限
	HasInNet  bool   `json:"hasInNet" excel:"内网"`  //内网
	HasOutNet bool   `json:"hasOutNet" excel:"外网"` //外网
}
