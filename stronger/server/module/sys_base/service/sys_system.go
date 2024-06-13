package service

import (
	"server/config"
	"server/global"
	"server/module/sys_base/model"
	"server/utils"
)

// @title    GetSystemConfig
// @description   读取配置文件
// @auth                     （2020/04/05  20:22）
// @return    err             error
// @return    conf            Server

func GetSystemConfig() (err error, conf config.Server) {
	return nil, global.GConfig
}

// @title    SetSystemConfig
// @description   set system config, 设置配置文件
// @auth                    （2020/04/05  20:22）
// @param     system         model.System
// @return    err            error

func SetSystemConfig(system model.System) (err error) {
	cs := utils.StructToMap(system.Config)
	for k, v := range cs {
		global.GVp.Set(k, v)
	}
	err = global.GVp.WriteConfig()
	return err
}
