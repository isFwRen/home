/**
 * @Author: xingqiyi
 * @Description: 录入定时重启
 * @Date: 2022/3/10 9:57 上午
 */

package schedule

import (
	"server/global"
	"server/module/sys_base/model"
	"server/module/sys_base/service"
	"server/utils"
)

//RestartTaskByCode 根据项目编码个端口重启录入
func RestartTaskByCode(proCode, port string) {
	global.GLog.Info(proCode + "  " + port)
	getPidShell := "netstat -lnpt | grep " + port + " | awk '{print $7}' | awk -F '/' '{print $1}' | xargs kill -9"
	err, stdout, stderr := utils.ShellOut(getPidShell)
	if err != nil {
		global.GLog.Error(err.Error())
		global.GLog.Error(stderr)
		global.GLog.Info(stdout)
		return
	}

	err = service.AddSysLogger(model.SysLogger{
		ProCode:         proCode,
		FunctionModule:  "系统重启",
		ModuleOperation: "系统重启",
		OperationCode:   "0",
		OperationName:   "系统",
		Content:         "",
		Api:             "",
		LogType:         3,
	})
	if err != nil {
		global.GLog.Error("插入日志:::" + err.Error())
	}
	return
}
