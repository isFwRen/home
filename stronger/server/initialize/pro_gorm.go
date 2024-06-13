/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/1 10:59 上午
 */

package initialize

import (
	"fmt"
	"os"
	"regexp"
	"server/global"
	"server/module/pro_conf/model"
	"server/utils"
	"strings"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func InitProGorm() {
	//global.GProConf = "pro_conf"
	//将项目配置加载到Redis，后序放到刷配置
	//proConf := utils.UpdateProjectConfToRedis(global.GProConf)
	var sysProjectList []model.SysProjectCache
	db := global.GDb.Model(&model.SysProject{})
	//非管理只连接当前项目的数据库 2023年04月12日15:08:16 xingqiyi
	if global.GConfig.System.Process != "common" {
		db = db.Where("code in (?)", strings.Split(global.GConfig.System.ProCode, "-"))
	}
	db = db.Where("code in (?)", global.GConfig.System.ProArr)
	//Joins("DownloadPaths", global.GDb.Where(&model.SysProDownloadPaths{IsDownload: true})).
	//Joins("UploadPaths", global.GDb.Where(&model.SysProDownloadPaths{IsUpload: true})).
	//Preload("DownloadPaths", "is_download = true").
	//Preload("UploadPaths", "is_download = true").
	err = db.Order("code").Find(&sysProjectList).Error
	if err != nil {
		panic(err)
	}
	global.ProDbMap = make(map[string]*gorm.DB)
	global.GProConf = make(map[string]model.SysProjectCache, 0)
	global.ProCodeId = make(map[string]string, 0)
	for _, conf := range sysProjectList {
		err = global.GDb.Where("pro_id = ? and is_download = true", conf.ID).First(&conf.DownloadPaths).Error
		err = global.GDb.Where("pro_id = ? and is_upload = true", conf.ID).First(&conf.UploadPaths).Error
		timeStart := time.Now()
		reg, err := regexp.Compile("^(task|export|type_load)$")
		if err != nil {
			global.GLog.Error("regexp Compile err::" + err.Error())
			continue
		}

		if global.GConfig.System.Process == "common" || (reg.MatchString(global.GConfig.System.Process) && conf.Code == global.GConfig.System.ProCode) {
			//(global.GConfig.System.Process != "common" && utils.HasItem(strings.Split(global.GConfig.System.ProCode, "-"), conf.Code)) {
			//非管理只加载当前项目的常量 2023年12月15日11:33:21 xingqiyi
			err, constMap := utils.LoadConst(conf.Code)
			global.GLog.Info("加载常量耗时", zap.Any(conf.Code, time.Since(timeStart).Milliseconds()))
			if err != nil {
				global.GLog.Error(conf.Code + "::" + conf.Name + "::加载常量错误" + err.Error())
				continue
			}
			conf.ConstTable = constMap
		}
		global.GProConf[conf.Code] = conf
		global.ProCodeId[conf.Code] = conf.ID
		global.ProCodeId[conf.ID] = conf.Code
		// 取出每个项目的数据库连接
		//var conf model.SysProject
		//_ = json.Unmarshal([]byte(item), &conf)
		if strings.Index(conf.DbHistory, "113.106") != -1 {
			continue
		}
		//本地测试专用
		if global.GConfig.System.Env == "develop" {
			conf.DbHistory = strings.Replace(conf.DbHistory, "192.168.202.18", "192.168.202.18", -1)
			conf.DbHistory = strings.Replace(conf.DbHistory, "port=5432", "port=5432", -1)
		}

		postgresConfig := postgres.Config{
			DSN:                  conf.DbHistory, // DSN data source name
			PreferSimpleProtocol: true,           // 禁用隐式 prepared statement
		}
		gormConfig := config(global.GConfig.Postgresql.Logger)

		if conf.DbHistory == "" {
			global.GLog.Error(conf.Code + "::" + conf.Name + "::历史数据库没配置")
			continue
		}
		dbConn, errPro := gorm.Open(postgres.New(postgresConfig), gormConfig)
		if errPro != nil {
			global.GLog.Error(conf.Code+"::"+conf.Name+"::历史数据库异常", zap.Any("err", errPro))
			continue
		}
		global.ProDbMap[conf.Code] = dbConn
		global.GLog.Info(conf.DbHistory)
		global.GLog.Info(conf.Code + "::" + conf.Name + "::历史数据库启动连接正常")
		GormsDBHistory(global.ProDbMap[conf.Code])
		sqlDB, _ := global.ProDbMap[conf.Code].DB()
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(10)
		_ = global.ProDbMap[conf.Code].Use(dbresolver.Register(dbresolver.Config{
			//读
			Replicas: []gorm.Dialector{postgres.Open(conf.DbHistory)},
			//写
			Sources: []gorm.Dialector{postgres.Open(conf.DbHistory)},
			//Sources: []gorm.Dialector{postgres.Open("user=postgres password=Change.Postgres host=14.29.64.91 dbname=stronger port=15432 sslmode=disable TimeZone=Asia/Shanghai")},
		}))
		err = global.ProDbMap[conf.Code].Callback().Create().Replace("gorm:before_create", beforeCreate)
		if err != nil {
			global.GLog.Error(fmt.Sprintf("GormPostgreSql:::Callback:::%v", err))
			panic(err)
		}
	}
	//InitProGormMove()
}

func InitProGormMove() {
	p := global.GConfig.Postgresql
	// dsn := "user=postgres password=Change.Postgres host=113.106.108.93 dbname=stronger port=54321 sslmode=disable TimeZone=Asia/Shanghai"
	dsn := "user=postgres password=Change.Postgres host=192.168.0.50 dbname=stronger port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	global.GLog.Info(dsn)
	postgresConfig := postgres.Config{
		DSN:                  dsn,                    // DSN data source name
		PreferSimpleProtocol: p.PreferSimpleProtocol, // 禁用隐式 prepared statement
	}
	gormConfig := config(p.Logger)
	if global.GDbMove, err = gorm.Open(postgres.New(postgresConfig), gormConfig); err != nil {
		global.GLog.Error("PostgreSql启动异常", zap.Any("err", err))
		os.Exit(0)
	} else {
		GormDBTables(global.GDbMove)
		sqlDB, _ := global.GDbMove.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
	}

	err := global.GDbMove.Use(dbresolver.Register(dbresolver.Config{
		//读
		Replicas: []gorm.Dialector{postgres.Open(dsn)},
		//写
		Sources: []gorm.Dialector{postgres.Open(dsn)},
		//Sources: []gorm.Dialector{postgres.Open("user=postgres password=Change.Postgres host=14.29.64.91 dbname=stronger port=15432 sslmode=disable TimeZone=Asia/Shanghai")},
	}))
	if err != nil {
		global.GLog.Error(fmt.Sprintf("GormPostgreSql:::Use:::%v", err))
		panic(err)
	}

	err = global.GDbMove.Callback().Create().Replace("gorm:before_create", beforeCreate)
	if err != nil {
		global.GLog.Error(fmt.Sprintf("GormPostgreSql:::Callback:::%v", err))
		panic(err)
	}

}
