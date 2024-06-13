package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"server/global"
	"server/module/sys_base/model"
	sys_base2 "server/module/sys_base/model/request"
)

// @title    Register
// @description   register, 用户注册
// @auth                     （2020/04/05  20:22）
// @param     u               model.SysUser
// @return    err             error
// @return    userInter       *SysUser

func Register(u sys_base2.RegisterStruct) (err error, userInter sys_base2.RegisterStruct) {
	var user sys_base2.RegisterStruct
	//var a model.AuthorityManager
	//判断用户名是否注册
	err1 := global.GDb.Where("id_card = ?", u.IDCard).First(&u).Error
	fmt.Print(err1)
	// notRegister为false表明读取到了 不能注册
	if !errors.Is(err1, gorm.ErrRecordNotFound) {
		fmt.Println("注册失败")
		return errors.New("该用户已注册"), userInter
	} else {
		// 否则 附加uuid 密码md5简单加密 注册
		//u.UUID = uuid.NewV4()
		//_ =global.GDb.Where("authority_name=?","员工").First(&a)
		//c := a.UserCount+1
		//_ =global.GDb.Model(model.AuthorityManager{}).Where("authority_name=?","员工").Update("user_count",c)
		fmt.Println("注册成功")
		err = global.GDb.Create(&u).Error
	}
	return err, user
}

func DingDing(u *model.SysUser) (err error, userInter *model.SysUser) {
	var user model.SysUser
	err = global.GDb.Where("username = ?", u.Code).First(&user).Error
	return err, &user
}

func Resignation(u sys_base2.ResignationStruct) (err error, userInter model.SysUser) {
	var user1, user2 model.SysUser
	//var user2 sys_base2.ResignationStruct
	fmt.Println(u)
	err1 := global.GDb.Where("code = ? AND nick_name = ? AND status = 't'", u.Username, u.Nickname).First(&user1).Error
	fmt.Println(err1)
	if !errors.Is(err1, gorm.ErrRecordNotFound) {
		fmt.Println("已查询到该用户")
		//err2:= global.GDb.Where("username = ? And nickname = ?", u.Code,u.Nickname).First(&user2).Error
		//fmt.Println(err2)
		//if !errors.Is(err2, gorm.ErrRecordNotFound) {
		//	fmt.Println("提交失败，该用户已提交过离职申请")
		//	return errors.New("该用户已提交过离职申请"), userInter
		//} else {
		//	fmt.Println("提交成功")
		//	err = global.GDb.Create(&u).Error
		//
		//}
		if user1.Status == false {
			return errors.New("该用户已提交过离职申请"), user1
		} else {
			err = global.GDb.Model(&user2).Where("code = ? And nick_name = ?", u.Username, u.Nickname).Update("reason", u.Reason).Update("status", false).Error
			return err, user2
		}
	} else {
		fmt.Println("查询失败，该用户尚未注册")
		return errors.New("未查询到该用户，请查询工号姓名是否填写正确"), userInter
	}
}
