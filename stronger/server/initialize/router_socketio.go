package initialize

import (
	"log"
	"server/global"
	"strings"
	"time"

	"go.uber.org/zap"

	socketIo "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
)

func InitSocketIO() {
	global.GSocketIo = socketIo.NewServer(&engineio.Options{
		PingTimeout:  time.Minute * 5,
		PingInterval: time.Minute * 2,
	})

	//释放单
	global.GSocketConnMap = make(map[string]socketIo.Conn, 0)
	global.GSocketIo.OnConnect("/global-release", func(s socketIo.Conn) error {
		userCode := getQueryParam(s.URL().RawQuery, "userCode")
		global.GLog.Info("global-release connected", zap.Any(userCode, s.ID()))
		global.GSocketConnMap[userCode] = s
		//可以广播同一个登录人的客户端
		//global.GSocketIo.JoinRoom("/global-export", userId, s)
		return nil
	})

	global.GSocketIo.OnDisconnect("/global-release", func(s socketIo.Conn, reason string) {
		global.GLog.Info("global-release OnDisconnect", zap.Any(reason, s.ID()))
		s.Close()
	})

	//全局消息通知
	global.GSocketConnMsgMap = make(map[string]socketIo.Conn, 0)
	global.GSocketIo.OnConnect("/global-message", func(s socketIo.Conn) error {
		userId := getQueryParam(s.URL().RawQuery, "userId")
		global.GLog.Info("global-message connected", zap.Any(userId, s.ID()))
		global.GSocketConnMsgMap[userId] = s
		return nil
	})

	//客户通知，通知函B0108
	global.GSocketIo.OnConnect("/global-notice", func(s socketIo.Conn) error {
		userId := getQueryParam(s.URL().RawQuery, "userId")
		global.GLog.Info("global-notice connected", zap.Any(userId, s.ID()))
		return nil
	})

	global.GSocketIo.OnError("/", func(s socketIo.Conn, e error) {
		global.GLog.Error("global-socket OnError")
		log.Println("meet error:", e)
	})
}

func getQueryParam(query string, name string) string {
	arr := strings.Split(query, "&")
	for _, obj := range arr {
		keyVal := strings.Split(obj, "=")
		if len(keyVal) > 0 && keyVal[0] == name {
			return keyVal[1]
		}
	}
	return ""
}
