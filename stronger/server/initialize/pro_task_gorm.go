/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022年01月19日15:17:09
 */

package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"server/global"
	"strings"
)

func InitProTaskGorm() {
	if global.ProDbMap == nil {
		global.ProDbMap = make(map[string]*gorm.DB)
	}
	for _, conf := range global.GProConf {
		if strings.Index(conf.DbTask, "113.106") != -1 {
			continue
		}
		//本地测试专用
		if global.GConfig.System.Env == "develop" {
			conf.DbTask = strings.Replace(conf.DbTask, "127.0.0.1", "192.168.202.23", -1)
			conf.DbTask = strings.Replace(conf.DbTask, "port=5433", "port=5433", -1)
		}
		postgresConfig := postgres.Config{
			DSN:                  conf.DbTask, // DSN data source name
			PreferSimpleProtocol: true,        // 禁用隐式 prepared statement
		}
		gormConfig := config(global.GConfig.Postgresql.Logger)

		if conf.DbTask == "" {
			global.GLog.Error(conf.Code + "::" + conf.Name + "::task数据库没配置")
			continue
		}
		dbConn, errPro := gorm.Open(postgres.New(postgresConfig), gormConfig)
		if errPro != nil {
			global.GLog.Error(conf.Code+"::"+conf.Name+"::task数据库异常", zap.Any("err", errPro))
			continue
		}
		global.ProDbMap[conf.Code+"_task"] = dbConn
		global.GLog.Info(conf.DbTask)
		global.GLog.Info(conf.Code + "::" + conf.Name + "::task数据库启动连接正常")
		sqlDB, _ := global.ProDbMap[conf.Code+"_task"].DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(10)
		_ = global.ProDbMap[conf.Code+"_task"].Use(dbresolver.Register(dbresolver.Config{
			//读
			Replicas: []gorm.Dialector{postgres.Open(conf.DbTask)},
			//写
			Sources: []gorm.Dialector{postgres.Open(conf.DbTask)},
			//Sources: []gorm.Dialector{postgres.Open("user=postgres password=Change.Postgres host=14.29.64.91 dbname=stronger port=15432 sslmode=disable TimeZone=Asia/Shanghai")},
		}))
		err = global.ProDbMap[conf.Code+"_task"].Callback().Create().Replace("gorm:before_create", beforeCreate)
		if err != nil {
			global.GLog.Error(fmt.Sprintf("InitProTaskGorm:::Callback:::%v", err))
			panic(err)
		}
	}
}
