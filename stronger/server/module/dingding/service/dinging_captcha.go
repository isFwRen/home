package service

import (
	"errors"
	"server/global"
	"server/module/sys_base/model"
)

func GetUserInformation(Phone string) (err error, userInter *model.SysUser) {
	var user model.SysUser
	err = global.GDb.Where("phone = ? AND status = ?", Phone, true).First(&user).Error
	if !user.Status {
		return errors.New("This user is not exists or resignation "), userInter
	}
	return err, &user
}