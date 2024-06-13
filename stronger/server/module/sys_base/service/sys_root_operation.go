package service

import (
	uuid "github.com/satori/go.uuid"
	"server/global"
	"server/module/sys_base/model"
	sys_base2 "server/module/sys_base/model/request"
)

// @title    QueryUser
// @description   query User by id_card, 通过身份证号查询工号
// @auth                      （2020/10/29  17：53）
// @param     info             request.QueryUserId
// @return    err              error

func QueryUser(u sys_base2.RegisterStruct) (err error, userInter *sys_base2.RegisterStruct) {
	var user sys_base2.RegisterStruct
	//fmt.Println(u.IDCard)
	err = global.GDb.Where("id_card = ?", u.IDCard).First(&user).Error
	//fmt.Println(user.Code)
	return err, &user
}

// @title    SetUserAuthority
// @description   set the authority of a certain user, 设置一个用户的权限
// @auth                     （2020/04/05  20:22）
// @param     uuid            UUID
// @param     authorityId     string
// @return    err             error

func SetUserAuthority(uuid uuid.UUID, authorityId string) (err error) {
	err = global.GDb.Where("uuid = ?", uuid).First(&model.SysUser{}).Update("authority_id", authorityId).Error
	return err
}

// @title    SetUserAuthority
// @description   set the authority of a certain user, 设置一个用户的权限
// @auth                     （2020/04/05  20:22）
// @param     uuid            UUID
// @param     authorityId     string
// @return    err             error

func DeleteUser(id string) (err error) {
	var user model.SysUser
	err = global.GDb.Where("id = ?", id).Delete(&user).Error
	return err
}
