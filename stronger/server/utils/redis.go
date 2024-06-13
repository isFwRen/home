package utils

import (
	"server/global"
	"time"
)

func GetRedisExport(path, userId string) (err error, val string) {
	val, err = global.GRedis.Get("user_export:" + path + ":" + userId).Result()
	return err, val
}

func SetRedisExport(path, userId string) (err error) {
	err = global.GRedis.Set("user_export:"+path+":"+userId, "true", 5*time.Minute).Err()
	return err
}

func DelRedisExport(path, userId string) (err error) {
	err = global.GRedis.Del("user_export:" + path + ":" + userId).Err()
	return err
}

func GetRedisCaptcha(phone string) (err error, val string) {
	val, err = global.GRedis.Get("user_captcha:" + phone).Result()
	return err, val
}

func SetRedisCaptcha(Captcha, phone string) (err error) {
	err = global.GRedis.Set("user_captcha:"+phone, Captcha, 1*time.Minute).Err()
	return err
}

func DelRedisCaptcha(phone string) (err error) {
	err = global.GRedis.Del("user_captcha:" + phone).Err()
	return err
}

func GetRedisAgency(proCode string) (err error, val []string) {
	val, err = global.GRedis.SMembers(proCode + ":download:set_agency").Result()
	return err, val
}

func SetRedisAgency(proCode, agency string) (err error) {
	err = global.GRedis.SAdd(proCode+":download:set_agency", agency).Err()
	return err
}

func DelRedisAgency(proCode string) (err error) {
	err = global.GRedis.Del(proCode + ":download:set_agency").Err()
	return err
}
