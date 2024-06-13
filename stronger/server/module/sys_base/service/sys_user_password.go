package service

import (
	"errors"
	"fmt"
	"server/global"
	"server/module/sys_base/model"
	"server/utils"
	"time"
)

func ResetPassword(u *model.SysUser, isIntranet bool) (err error, userInter *model.SysUser) {
	var user model.SysUser
	newPassword := "123456"
	u.Password = utils.MD5V([]byte(u.Password))
	if !isIntranet {
		err = global.GDb.Model(&user).Where("phone = ?", u.Phone).Update("password", utils.MD5V([]byte(newPassword))).Error
	} else {
		err = global.GDb.Model(&user).Where("code = ?", u.Code).Update("password", utils.MD5V([]byte(newPassword))).Error
	}
	return err, u
}

func ChangePassword(u *model.SysUser, oldPassword, newPassword string, isIntranet bool) (err error, userInter *model.SysUser) {
	var user model.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	// EntryDate, _ := time.ParseInLocation("2006-01-02 15:04:05", adduser.EntryDate, time.Local)
	if !isIntranet {
		fmt.Println("isIntranet", isIntranet)
		err = global.GDb.Model(&user).Where("phone = ? AND password = ? AND status = ?", u.Phone, u.Password, true).Find(&user).Error
		if err != nil {
			return err, u
		}
		fmt.Println("外网：", user.Password, oldPassword)
		if user.Password != utils.MD5V([]byte(oldPassword)) {
			return errors.New("原密码错误! "), u
		}
		err = global.GDb.Model(&user).Where("phone = ? AND password = ? AND status = ?", u.Phone, u.Password, true).Updates(map[string]interface{}{"password": utils.MD5V([]byte(newPassword)), "password_date": time.Now()}).Error
	} else {
		fmt.Println("isIntranet", isIntranet)
		err = global.GDb.Model(&user).Where("code = ? AND password = ? AND status = ?", u.Code, u.Password, true).Find(&user).Error
		if err != nil {
			return err, u
		}
		fmt.Println("外网：", user.Password, utils.MD5V([]byte(oldPassword)))
		if user.Password != utils.MD5V([]byte(oldPassword)) {
			return errors.New("原密码错误! "), u
		}
		err = global.GDb.Model(&user).Where("code = ? AND password = ? AND status = ?", u.Code, u.Password, true).Update("password", utils.MD5V([]byte(newPassword))).Updates(map[string]interface{}{"password": utils.MD5V([]byte(newPassword)), "password_date": time.Now()}).Error
	}
	return err, u
}
