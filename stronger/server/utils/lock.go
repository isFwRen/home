package utils

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"time"
)

var locked = "locked"

func Lock(key string, expiry time.Duration) (bool, error) {
	// 使用 SETNX 命令尝试获取锁，只有当键不存在时才设置
	setRes, err := global.GRedis.SetNX(key, locked, expiry).Result()
	if err != nil {
		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	if setRes {
		// 成功获取锁
		return true, nil
	}
	return false, errors.New("请稍后重试")
}

func Unlock(key string) error {
	// 使用 Lua 脚本原子性地释放锁，只有当锁的值与设置时的值相同时才删除
	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		end
		return 0
	`
	res, err := global.GRedis.Eval(luaScript, []string{key}, locked).Result()
	if err != nil || res.(int64) == 0 {
		global.GLog.Error(fmt.Sprintf("failed to release lock，res:%v", res), zap.Error(err))
		return fmt.Errorf("failed to release lock: %w", err)
	}
	return nil
}
