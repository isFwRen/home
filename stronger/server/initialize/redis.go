package initialize

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"server/global"
)

func Redis() {
	redisCfg := global.GConfig.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
		global.GLog.Error("redis connect ping failed, err:", zap.Any("err", err))
		panic(err)
	} else {
		global.GLog.Info("redis connect ping response:", zap.String("pong", pong))
		global.GRedis = client
	}

	///添加
	//ret, err := client.SAdd("set_test", "11", "22", "33", "44").Result()
	//log.Println(ret, err)

}
