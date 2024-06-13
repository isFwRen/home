package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"server/global"
	"server/module/sys_base/model"
)

// @title    CreateAuthority
// @description   创建一个角色
// @auth                     （2020/04/05  20:22）
// @param     auth            model.SysAuthority
// @return                    error
// @return    authority       model.SysAuthority

func CreateAuthority(auth model.SysAuthority) (err error, authority model.SysAuthority) {
	var authorityBox model.SysAuthority
	if !errors.Is(global.GDb.Where("authority_id = ?", auth.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同角色id"), auth
	}
	err = global.GDb.Create(&auth).Error
	return err, auth
}

// @title    CopyAuthority
// @description   复制一个角色
// @auth                     （2020/04/05  20:22）
// @param     copyInfo        resp.SysAuthorityCopyResponse
// @return                    error
// @return    authority       model.SysAuthority

// @title    UpdateAuthority
// @description   更改一个角色
// @auth                     （2020/04/05  20:22）
// @param     auth            model.SysAuthority
// @return                    error
// @return    authority       model.SysAuthority

//func UpdateAuthority(auth model.SysAuthority) (err error, authority model.SysAuthority) {
//	err = global.G_DB.Where("authority_id = ?", auth.AuthorityId).First(&model.SysAuthority{}).Updates(&auth).Error
//	return err, auth
//}

// @title    GetInfoList
// @description   删除文件切片记录
// @auth                     （2020/04/05  20:22）
// @param     info            request.PaveInfo
// @return                    error
// 分页获取数据

//func GetAuthorityInfoList(info request.PageInfo) (err error, list interface{}, total int) {
//	limit := info.PageSize
//	offset := info.PageSize * (info.Page - 1)
//	db := global.G_DB
//	var authority []model.SysAuthority
//	err = db.Limit(limit).Offset(offset).Preload("DataAuthorityId").Where("parent_id = 0").Find(&authority).Error
//	if len(authority) > 0 {
//		for k := range authority {
//			err = findChildrenAuthority(&authority[k])
//		}
//	}
//	return err, authority, total
//}

// @title    GetAuthorityInfo
// @description   获取所有角色信息
// @auth                     （2020/04/05  20:22）
// @param     auth            model.SysAuthority
// @return                    error
// @param     authority       model.SysAuthority

func GetAuthorityInfo(auth model.SysAuthority) (err error, sa model.SysAuthority) {
	err = global.GDb.Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(&sa).Error
	return err, sa
}

// @title    SetDataAuthority
// @description   设置角色资源权限
// @auth                     （2020/04/05  20:22）
// @param     auth            model.SysAuthority
// @return                    error

//func SetDataAuthority(auth model.SysAuthority) error {
//	var s model.SysAuthority
//	global.G_DB.Preload("DataAuthorityId").First(&s, "authority_id = ?", auth.AuthorityId)
//	err := global.G_DB.Model(&s).Association("DataAuthorityId").Replace(&auth.DataAuthorityId).Error
//	return err
//}

// @title    findChildrenAuthority
// @description   查询子角色
// @auth                     （2020/04/05  20:22）
// @param     auth            *model.SysAuthority
// @return                    error

//func findChildrenAuthority(authority *model.SysAuthority) (err error) {
//	err = global.G_DB.Preload("DataAuthorityId").Where("parent_id = ?", authority.AuthorityId).Find(&authority.Children).Error
//	if len(authority.Children) > 0 {
//		for k := range authority.Children {
//			err = findChildrenAuthority(&authority.Children[k])
//		}
//	}
//	return err
//}

func GetAllAuthority() (err error, u []model.AuthorityManager) {
	_ = global.GDb.Order("id asc").Find(&u)
	return err, u

}

func CreateNewAuthority(a model.AuthorityManager) (err error) {
	var auth model.AuthorityManager
	if !errors.Is(global.GDb.Where("authority_name = ?", a.AuthorityName).First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在相同角色名称")
	} else {
		err = global.GDb.Create(&a).Error
	}
	return err
}

func UpdateAuthority(a model.AuthorityManager) (err error) {
	err = global.GDb.Model(model.AuthorityManager{}).Where("authority_name = ?", a.AuthorityName).Updates(model.AuthorityManager{
		AuthorityName: a.AuthorityName,
		AuthorityState: a.AuthorityState,
		Describe: a.Describe,
	}).Error
	return err
}

func UpdateAuthorityPower(a model.AuthorityManager) (err error) {
	fmt.Println(a)
	err = global.GDb.Model(model.AuthorityManager{}).Where("authority_name = ?", a.AuthorityName).Updates(a).Error
	return err
}

func DeleteAuthority(a model.AuthorityManager) (err error) {
	var auth model.AuthorityManager
	err = global.GDb.Where("authority_name = ?", a.AuthorityName).Delete(&auth).Error
	return err
}

func QueryAuthorityPower(a model.AuthorityManager) (err error,auth model.AuthorityManager){
	err = global.GDb.Where("authority_name = ?", a.AuthorityName).First(&auth).Error
	return err,auth
}
