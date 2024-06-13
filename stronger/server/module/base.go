/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/10/27 10:03 上午
 */

package module

import (
	"context"
	"fmt"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
	"os"
	"runtime"
	global2 "server/global"
	initialize2 "server/initialize"
	"server/middleware"
	"server/utils"
	"time"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/global"
	"xingqiyi.com/gitlab-instance-09305a81/ums_server.git/initialize"
)

func Base() {
	println(`系统类型：`, runtime.GOOS)
	println(`系统架构：`, runtime.GOARCH)
	println(`CPU核数：`, runtime.GOMAXPROCS(0))
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	println(`电脑名称：`, name)

	//初始化常量表库
	fmt.Println("初始化常量表库-couchdb, admin-HL2018dbZHHLYXGS")
	initialize.InitCouchdb()
	global2.GCouchdbClient = global.GCouchdbClient

	//初始化数据库连接
	//initialize.GormPostgreSql()

	//初始化数据库连接
	global.GDb = initialize.GormPostgreSql(global.GConfig.Postgresql)
	global2.GDb = global.GDb
	global.GUserDb = initialize.GormPostgreSql(global.GConfig.UserDB)
	global2.GUserDb = global.GUserDb
	//初始化Redis
	initialize.Redis()
	global2.GRedis = global.GRedisMap[global.GConfig.System.Name]

	//初始化项目数据库连接
	initialize2.InitProGorm()

	//初始化项目task数据库连接
	initialize2.InitProTaskGorm()

	//初始化id生成器
	utils.IdWorker()

	global.WhiteList = middleware.WhiteList
	global.LogWhitelists = middleware.LogWhitelists
	global2.GStore = base64Captcha.NewMemoryStore(global.GConfig.Captcha.GCLimitNumber, global.GConfig.Captcha.Expiration*time.Minute)
	err = initialize2.SetGlobalAPICache()
	if err != nil {
		global.GLog.Error(fmt.Sprintf("Base:::SetGlobalAPIMap:::%v", err))
		panic(err)
	}
	err = initialize2.SetGlobalUserMap()
	if err != nil {
		global.GLog.Error(fmt.Sprintf("Base:::SetGlobalUserMap:::%v", err))
		panic(err)
	}
}

func CloseDatabase() {
	db, err := global.GDb.DB()
	if err != nil {
		global.GLog.Error("CloseDatabase0", zap.Any("错误❎", err))
	}
	//if err1 := recover(); err1 != nil {
	//	global.GLog.Error("main", zap.Any("错误❎", err1))
	//关闭配置数据库连接
	err = db.Close()
	if err != nil {
		global.GLog.Error("CloseDatabase1", zap.Any("错误❎", err))
	}
	//关闭项目数据库连接
	for _, v := range global.ProDbMap {
		proDb, _ := v.DB()
		errTemp := proDb.Close()
		if errTemp != nil {
			global.GLog.Error("CloseDatabase2", zap.Any("错误❎", err))
			continue
		}
	}
	//}
	err = global.GCouchdbClient.Close(context.Background())
	if err != nil {
		global.GLog.Error("CloseDatabase3", zap.Any("错误❎", err))
		return
	}
}
