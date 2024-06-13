/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/9/1 5:35 下午
 */

package utils

import (
	"server/global"
)

var (
	GWorker *SnowflakeIdWorker
)

func IdWorker() {
	worker, err := CreateWorker(0, 0)
	if err != nil {
		global.GLog.Error(err.Error())
		return
	}
	GWorker = worker
}
