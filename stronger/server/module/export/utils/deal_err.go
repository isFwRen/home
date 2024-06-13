/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/24 5:41 下午
 */

package utils

import (
	"go.uber.org/zap"
	"server/global"
	"sync"
)

func DealErr(wg *sync.WaitGroup) {
	wg.Done()
	if r := recover(); r != nil {
		global.GLog.Error("", zap.Any("err", r))
	}
}
