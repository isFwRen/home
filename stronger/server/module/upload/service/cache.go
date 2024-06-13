package service

import (
	"server/global"
	"time"
)

func GetGuoShouCacheSecret() (token string, err error) {
	return global.GRedis.Get("guoshou:access_token:").Result()
}

func SetGuoShouCacheSecret(token string, expiresAt int) (err error) {
	err = global.GRedis.Set("guoshou:access_token:", token, time.Duration(expiresAt)).Err()
	return err
}
