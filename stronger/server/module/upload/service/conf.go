/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/30 4:24 下午
 */

package service

import (
	"server/global"
	"server/module/pro_conf/model"
)

func GetProIsAutoReturn(proCode string) (err error, pro model.SysProject) {
	err = global.GDb.Model(&model.SysProject{}).
		Where("code = ?", proCode).
		First(&pro).Error
	return err, pro
}
