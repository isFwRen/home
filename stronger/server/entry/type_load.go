package main

import (
	"fmt"
	core2 "server/core"
	"server/global"
	"server/module"
	"server/module/load/project/B0102"
	"server/module/load/project/B0103"
	"server/module/load/project/B0106"
	"server/module/load/project/B0108"
	"server/module/load/project/B0110"
	"server/module/load/project/B0113"
	"server/module/load/project/B0114"
	"server/module/load/project/B0116"
	"server/module/load/project/B0118"
	"server/module/load/project/B0121"
	"server/module/load/project/B0122"
	"server/module/load/service"
	"server/module/pro_conf/model"
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/core"
	xingqiyiGlobal "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	module2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module"

	"go.uber.org/zap"
	// "server/module/load/"
)

func main() {
	fmt.Println("Hello, type_load!")
	core.InitConfig()
	core2.InitConfig()
	global.GConfig.System.Process = "type_load"
	core.InitZap()
	global.GLog = xingqiyiGlobal.GLog

	module.Base()
	// 程序结束前关闭数据库链接
	defer module2.CloseDatabase()
	typeLoadProcess(global.GConfig.System.ProCode)
	timeAfterTrigger := time.After(30 * time.Second)
	<-timeAfterTrigger
	fmt.Println("程序结束")

}

func typeLoadProcess(proCode string) {
	proID := global.ProCodeId[proCode]
	err, temp := service.GetSysProTempByProIdAndName(proID, "MB002")
	if err != nil {
		global.GLog.Error("temp", zap.Error(err))
	}
	global.GLog.Info("temp", zap.Any("temp", temp))
	err, tempBlocks := service.GetSysProTempBlockByTempId(temp.ID)
	if err != nil {
		global.GLog.Error("GetSysProTempBlockByTempId", zap.Error(err))
		return
	}
	err, cBlocks := service.SelectCropBlocks(proCode + "_task")
	if err != nil {
		global.GLog.Error("SelectCropBlocks", zap.Error(err))
		return
	}
	var CacheFieldConf = make(map[string]model.SysProField)
	err, fields := service.GetSysFields(proID)
	for _, field := range fields {
		CacheFieldConf[field.Code] = field
	}
	if err != nil {
		global.GLog.Error("GetSysFields", zap.Error(err))
		return
	}
	global.GLog.Info("typeLoadProcess", zap.Any("cBlocks", cBlocks))
	for ii, block := range cBlocks {
		global.GLog.Info("typeLoadProcess", zap.Any("block", block))
		global.GLog.Info("typeLoadProcess", zap.Any("ii", ii))
		switch proCode {
		case "B0118":
			err := B0118.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0108":
			err := B0108.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0114":
			err := B0114.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0113":
			err := B0113.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0121":
			err := B0121.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0103":
			err := B0103.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0106":
			err := B0106.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0110":
			err := B0110.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0122":
			err := B0122.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0116":
			err := B0116.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		case "B0102":
			err := B0102.TypeCrop(proCode, block, tempBlocks, CacheFieldConf)
			if err != nil {
				global.GLog.Error("typeLoad", zap.Error(err))
			}
		default:
			fmt.Println("不存在的加载进程项目编号:", proCode)
		}
	}

}
