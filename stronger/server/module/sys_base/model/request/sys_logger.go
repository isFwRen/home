/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/2/9 1:58 下午
 */

package request

import (
	"server/module/sys_base/model"
	"time"
)

type SysLogger struct {
	model.BasePageInfo
	ProCode         string    `json:"proCode" form:"proCode"`
	StartTime       time.Time `json:"startTime" form:"startTime" binding:"required"`
	EndTime         time.Time `json:"endTime" form:"endTime" binding:"required"`
	FunctionModule  string    `json:"functionModule" form:"functionModule"`
	ModuleOperation string    `json:"moduleOperation" form:"moduleOperation"`
	OperationPeople string    `json:"operationPeople" form:"operationPeople"`
	LogType         int       `json:"logType" form:"logType" binding:"required"`
}
