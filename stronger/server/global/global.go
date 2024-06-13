package global

import (
	"database/sql"
	"server/config"
	"server/module/pro_conf/model"

	"github.com/go-kivik/kivik/v3"
	"github.com/go-redis/redis"
	socketio "github.com/googollee/go-socket.io"
	"github.com/mojocn/base64Captcha"
	"github.com/robfig/cron"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	GDb               *gorm.DB
	GUserDb           *gorm.DB
	GDbMove           *gorm.DB
	GRedis            *redis.Client
	GConfig           config.Server
	GVp               *viper.Viper //用来配置yaml
	GLog              *zap.Logger
	ProDbMap          map[string]*gorm.DB
	GCouchdbClient    *kivik.Client
	GSocketIo         *socketio.Server
	GSocketConnMap    map[string]socketio.Conn
	GSocketConnMsgMap map[string]socketio.Conn
	GTaskDb           *sql.DB
	GCron             *cron.Cron
	GTaskCron         *cron.Cron
	GDownloadCron     *cron.Cron
	//GIdWorker    *utils.SnowflakeIdWorker
	GStore base64Captcha.Store
)

// GProConf
var (
	GProConf      = make(map[string]model.SysProjectCache, 0) //项目配置
	ProCodeId     = make(map[string]string, 0)                //项目code转id,id转code
	UserCodeName  = make(map[string]string, 0)                //用户code转name
	UserCodeAndId = make(map[string][]string, 0)              //用户code转id,objectId
	ApiPathTitle  = make(map[string]string, 0)                //api路径转标题
)

var XUserId = "x-user-id"
var XToken = "x-token"
var XProCode = "pro-code"
var XProcess = "process"
var XCode = "x-code"
