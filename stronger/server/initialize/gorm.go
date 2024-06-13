package initialize

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"os"
	"reflect"
	"server/global"
	"server/module/sys_base/model"
	"strings"
	model2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
)

var err error

// GormDBTables 注册数据库表专用
func GormDBTables(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		global.GLog.Error("register table failed", zap.Any("err", err))
		os.Exit(0)
	}
	global.GLog.Info("register table success")

}

func GormsDBHistory(db *gorm.DB) {
	err := db.AutoMigrate(
	//model.History{},
	)
	if err != nil {
		global.GLog.Error("register history failed", zap.Any("err", err))
	}
	global.GLog.Info("register history success")
}

// GormPostgreSql 初始化PostgreSql数据库
func GormPostgreSql() {
	p := global.GConfig.Postgresql
	dsn := "user=" + p.Username + " password=" + p.Password + " host=" + p.Host + " dbname=" + p.Dbname + " port=" + p.Port + " " + p.Config
	global.GLog.Info(dsn)
	postgresConfig := postgres.Config{
		DSN:                  dsn,                    // DSN data source name
		PreferSimpleProtocol: p.PreferSimpleProtocol, // 禁用隐式 prepared statement
	}
	gormConfig := config(p.Logger)
	if global.GDb, err = gorm.Open(postgres.New(postgresConfig), gormConfig); err != nil {
		global.GLog.Error("PostgreSql启动异常", zap.Any("err", err))
		os.Exit(0)
	} else {
		GormDBTables(global.GDb)
		sqlDB, _ := global.GDb.DB()
		sqlDB.SetMaxIdleConns(p.MaxIdleConns)
		sqlDB.SetMaxOpenConns(p.MaxOpenConns)
	}

	err := global.GDb.Use(dbresolver.Register(dbresolver.Config{
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

	//注册回调
	// Define callbacks for creating
	//func init() {
	//	DefaultCallback.Create().Register("gorm:begin_transaction", beginTransactionCallback)
	//	DefaultCallback.Create().Register("gorm:before_create", beforeCreateCallback)
	//	DefaultCallback.Create().Register("gorm:save_before_associations", saveBeforeAssociationsCallback)
	//	DefaultCallback.Create().Register("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	//	DefaultCallback.Create().Register("gorm:create", createCallback)
	//	DefaultCallback.Create().Register("gorm:force_reload_after_create", forceReloadAfterCreateCallback)
	//	DefaultCallback.Create().Register("gorm:save_after_associations", saveAfterAssociationsCallback)
	//	DefaultCallback.Create().Register("gorm:after_create", afterCreateCallback)
	//	DefaultCallback.Create().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
	//}
	err = global.GDb.Callback().Create().Replace("gorm:before_create", beforeCreate)
	if err != nil {
		global.GLog.Error(fmt.Sprintf("GormPostgreSql:::Callback:::%v", err))
		panic(err)
	}
	err = SetGlobalUserMap()
	if err != nil {
		global.GLog.Error(fmt.Sprintf("GormPostgreSql:::SetGlobalUserMap:::%v", err))
		panic(err)
		return
	}

	err = SetGlobalAPIMap()
	if err != nil {
		global.GLog.Error(fmt.Sprintf("GormPostgreSql:::SetGlobalAPIMap:::%v", err))
		panic(err)
		return
	}
}

func SetGlobalUserMap() error {
	var users []model2.SysUser
	err1 := global.GUserDb.Where("status = 2 and code != '-1'").Select("id,code,name").Find(&users).Error
	if err1 != nil {
		return err1
	}
	for _, user := range users {
		global.UserCodeName[user.Code] = user.Name
		global.UserCodeAndId[user.Code] = []string{user.ID}
		//global.UserCodeAndId[user.ID] = []string{user.Code, user.ObjectId}
	}
	return err1
}

func SetGlobalAPIMap() error {
	var apis []model.SysApi
	err1 := global.GDb.Where("path != ''").Select("path,title").Find(&apis).Error
	if err1 != nil {
		return err1
	}
	for _, api := range apis {
		global.ApiPathTitle[api.Path] = api.Title
	}
	return err1
}

func SetGlobalAPICache() error {
	var apis []model.SysApi
	err1 := global.GUserDb.Where("path != '' and sys_code = ?", global.GConfig.System.Name).Select("path,title").Find(&apis).Error
	if err1 != nil {
		return err1
	}
	for _, api := range apis {
		global.ApiPathTitle[api.Path] = api.Title
	}
	return err1
}

// config 根据配置决定是否开启日志
func config(mod bool) (c *gorm.Config) {
	if mod {
		c = &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Info),
			DisableForeignKeyConstraintWhenMigrating: true,
		}
		//if global.GConfig.System.Process == "common" || global.GConfig.System.Process == "export" {
		//	// 创建日志文件
		//	file, err := os.OpenFile("gorm_"+global.GConfig.System.Process+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		//	if err != nil {
		//		log.Fatal(err)
		//	}
		//	c = &gorm.Config{
		//		Logger: logger.New(
		//			log.New(file, "\r\n", log.LstdFlags),
		//			logger.Config{
		//				SlowThreshold: time.Second,
		//				LogLevel:      logger.Info,
		//				Colorful:      true,
		//			},
		//		),
		//		DisableForeignKeyConstraintWhenMigrating: true,
		//	}
		//}
		//defer file.Close()
	} else {
		c = &gorm.Config{
			Logger:                                   logger.Default.LogMode(logger.Silent),
			DisableForeignKeyConstraintWhenMigrating: true,
		}
	}
	return
}

func beforeCreate(db *gorm.DB) {
	defer func() {
		if r := recover(); r != nil {
			panic(fmt.Sprintf("beforeCreate→→→→→→%v", r))
		}
	}()
	if db.Statement.Schema != nil {
		switch db.Statement.ReflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			//global.GLog.Info("insert多条数据")
			for i := 0; i < db.Statement.ReflectValue.Len(); i++ {
				// 从字段中获取数值
				//nextId := utils.GWorker.NextId()
				nextId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
				//global.GLog.Info("id:::" + strconv.FormatInt(nextId, 10))
				if db.Statement.Schema.PrioritizedPrimaryField == nil {
					panic(fmt.Sprintf("该结构体没有主键id:::%v", db.Statement.Table))
				}
				if _, isZero := db.Statement.Schema.PrioritizedPrimaryField.ValueOf(db.Statement.ReflectValue.Index(i)); !isZero {
					//global.GLog.Warn("存在id新增", zap.Any(db.Statement.Schema.Table, fieldValue))
					continue
				}
				//err = db.Statement.Schema.PrioritizedPrimaryField.Set(db.Statement.ReflectValue.Index(i), strconv.FormatInt(nextId, 10))
				err = db.Statement.Schema.PrioritizedPrimaryField.Set(db.Statement.ReflectValue.Index(i), nextId)
				if err != nil {
					panic(fmt.Sprintf("beforeCreate→→→→→→%v", err))
				}
			}
		case reflect.Struct:
			//global.GLog.Info("insert单条数据")
			//nextId := utils.GWorker.NextId()
			nextId := strings.Replace(uuid.NewV4().String(), "-", "", -1)
			//global.GLog.Info("id:::" + strconv.FormatInt(nextId, 10))
			//global.GLog.Info("id:::" + nextId)
			if db.Statement.Schema.PrioritizedPrimaryField == nil {
				panic(fmt.Sprintf("该结构体没有主键id:::%v", db.Statement.Table))
			}
			if _, isZero := db.Statement.Schema.PrioritizedPrimaryField.ValueOf(db.Statement.ReflectValue); !isZero {
				//global.GLog.Warn("存在id新增", zap.Any(db.Statement.Schema.Table, fieldValue))
				return
			}
			//err = db.Statement.Schema.PrioritizedPrimaryField.Set(db.Statement.ReflectValue, strconv.FormatInt(nextId, 10))
			err = db.Statement.Schema.PrioritizedPrimaryField.Set(db.Statement.ReflectValue, nextId)
			if err != nil {
				panic(fmt.Sprintf("beforeCreate→→→→→→%v", err))
			}
		}
		// 当前实体的所有字段
		//db.Statement.Schema.Fields

		// 当前实体的所有主键字段
		//db.Statement.Schema.PrimaryFields

		// 优先主键字段：带有数据库名称`id`或第一个定义的主键的字段。
		//db.Statement.Schema.PrioritizedPrimaryField

		// 当前模型的所有关系
		//db.Statement.Schema.Relationships

		// 使用字段名或数据库名查找字段
		//field := db.Statement.Schema.LookUpField("Name")

		// processing
	}
}
