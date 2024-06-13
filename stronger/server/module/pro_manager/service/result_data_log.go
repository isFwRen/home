/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/12/3 2:39 下午
 */

package service

import (
	"server/global"
	"server/module/pro_manager/model"
)

//InsertLog 插入一个日志
func InsertLog(proCode string, r model.ResultDataLog) error {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}
	err := db.Model(&model.ResultDataLog{}).Create(&r).Error
	return err
}
