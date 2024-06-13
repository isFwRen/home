/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/2/9 9:21 上午
 */

package model

type SysLogger struct {
	Model
	ProCode         string `json:"proCode"`         //项目编码
	FunctionModule  string `json:"functionModule"`  //功能模块
	ModuleOperation string `json:"moduleOperation"` //名称
	OperationName   string `json:"operationName"`   //修改人名字
	OperationCode   string `json:"operationCode"`   //修改人工号
	Content         string `json:"content"`         //内容,保留字段
	Api             string `json:"api"`             //url
	LogType         int    `json:"logType"`         //日志类型（1：项目开发，2：项目管理，3：系统重启，4：录入管理）
}
