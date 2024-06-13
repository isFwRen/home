package service

import (
	"server/global"
	"server/module/sys_base/model"
)

func GetUsersInformation(code string) (err error, list interface{}, total int64) {
	var user []string
	err = global.GDb.Model(&model.SysUser{}).Where("code = ? AND status = 'true'", code).Pluck("nick_name", &user).Count(&total).Error
	return err, user, total
}
