/**
 * @Author: xingqiyi
 * @Description: 迁移用户数据
 * @Date: 2021/12/30 12:54 下午
 */

package main

import (
	core2 "server/core"
	"server/global"
	"server/module"
	"server/module/sys_base"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"
)

func main() {
	core.InitConfig()
	core2.InitConfig()
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog
	module.Base()

	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()

	sys_base.Move("", true)
}
