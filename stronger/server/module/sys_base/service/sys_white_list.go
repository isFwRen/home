package service

import (
	"errors"
	"fmt"
	"regexp"
	"server/global"
	m "server/module/pro_conf/model"
	"server/module/sys_base/model"
	"server/module/sys_base/model/request"
	u "server/utils"
	"strconv"
	"time"
	model2 "xingqiyi.com/gitlab-instance-09305a81/ums_server.git/module/sys_base/model"
)

func GetWhileList(info request.GetWhiteList) (err error, list interface{}, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	//err = global.GDb.Raw("select * from sys_pro_permissions where (user_id) in (select id from sys_users where status = true)").Count(&total).Error
	//if err != nil {
	//
	//}
	db := global.GUserDb.Model(&model2.SysUser{}).Where("status = ?", 2)
	if info.UserCode != "" {
		reg1 := regexp.MustCompile("^([A-Za-z0-9]+)$")
		if reg1.MatchString(info.UserCode) || info.UserCode[0] == 'P' {
			db = db.Where("code = ? ", info.UserCode)
		} else {
			db = db.Where("name = ? ", info.UserCode)
		}
		limit = 1
	}
	var t int64
	err = db.Count(&t).Error
	if err != nil {
		return err, nil, 0
	}
	fmt.Println("t", t)
	if t == 0 {
		return errors.New("没有查询到相关在职用户"), nil, 0
	}

	var code []string
	if info.UserCode != "" {
		err = db.Pluck("DISTINCT(code)", &code).Error
	} else {
		err = db.Limit(limit).Offset(offset).Pluck("DISTINCT(code)", &code).Error
	}

	if err != nil {
		return err, nil, 0
	}

	fmt.Println("code", code)

	proDb := global.ProDbMap[info.ProCode]
	var userCode []string
	err = proDb.Model(&model.SysWhiteList{}).Where("(user_code) in (?) AND temp_name = ? ", code, info.TempName).Pluck("DISTINCT(user_code)", &userCode).Error
	if err != nil {
		return err, nil, 0
	}
	fmt.Println("userCode", userCode)

	var wl []model.SysWhiteList
	if info.UserCode == "" {
		err = proDb.Model(&model.SysWhiteList{}).Limit(limit).Offset(offset).Where("pro_code = ? AND temp_name = ? ", info.ProCode, info.TempName).Find(&wl).Error
		if err != nil {
			return err, nil, 0
		}
	} else {
		reg1 := regexp.MustCompile("^([A-Za-z0-9]+)$")
		if reg1.MatchString(info.UserCode) || info.UserCode[0] == 'P' {
			err = proDb.Model(&model.SysWhiteList{}).Limit(limit).Offset(offset).Where("pro_code = ? AND temp_name = ? AND user_code = ? ", info.ProCode, info.TempName, info.UserCode).Find(&wl).Error
			if err != nil {
				return err, nil, 0
			}
		} else {
			err = proDb.Model(&model.SysWhiteList{}).Limit(limit).Offset(offset).Where("pro_code = ? AND temp_name = ? AND user_name = ? ", info.ProCode, info.TempName, info.UserCode).Find(&wl).Error
			if err != nil {
				return err, nil, 0
			}
		}
	}

	var proName m.SysProject
	_ = global.GDb.Model(&m.SysProject{}).Where("code = ? ", info.ProCode).Find(&proName).Error
	for _, v := range code {
		isnotexist := true
		for _, f := range userCode {
			if v == f {
				isnotexist = false
				break
			}
		}
		if isnotexist {
			var wlItem model.SysWhiteList
			wlItem.ProCode = info.ProCode
			wlItem.UserCode = v
			var c model.SysUser
			_ = global.GDb.Model(&model.SysUser{}).Where("code = ? ", v).Find(&c).Error
			wlItem.UserName = c.NickName
			wlItem.TempName = info.TempName
			wlItem.BlockPermissions = []string{}
			wlItem.ProName = proName.Name
			err = proDb.Model(&model.SysWhiteList{}).Create(&wlItem).Error
			if err != nil {
				return err, nil, 0
			}
			wl = append(wl, wlItem)
		}
	}

	return err, wl, t
}

func EditWhiteList(info request.EditWhiteListArr) (err error) {
	db := global.ProDbMap[info.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	for _, v := range info.Arr {
		var total int64
		err = tx.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? ", info.ProCode, info.TempName).Where("user_code = ? AND user_name = ? ", v.UserCode, v.UserName).Count(&total).Error
		if err != nil {
			return err
		}
		if total > 0 {
			err = tx.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? ", info.ProCode, info.TempName).Where("user_code = ? AND user_name = ? ", v.UserCode, v.UserName).Updates(map[string]interface{}{
				"block_permissions": v.BlockPermissions,
			}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return errors.New("没有该用户")
		}
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func CopyWhileList(info request.CopyWhileList) (err error) {

	db := global.ProDbMap[info.ProCode]
	if db == nil {
		return global.ProDbErr
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var total int64
	var CodeWhiteList model.SysWhiteList
	reg1 := regexp.MustCompile("^([A-Za-z0-9]+)$")
	if reg1.MatchString(info.Code) || info.Code[0] == 'P' {
		err = tx.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? AND user_code = ? ", info.ProCode, info.TempName, info.Code).Find(&CodeWhiteList).Count(&total).Error
		if err != nil {
			return err
		}
	} else {
		err = tx.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? AND user_name = ? ", info.ProCode, info.TempName, info.Code).Find(&CodeWhiteList).Count(&total).Error
		if err != nil {
			return err
		}
	}

	if total > 1 {
		return errors.New("存在2条以上的同名或同工号的数据！")
	}

	for _, v := range info.CopyCode {
		if reg1.MatchString(v) || v[0] == 'P' {
			var total int64
			var CodeWhiteList2 model.SysWhiteList
			err = tx.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? AND user_code = ? ", info.ProCode, info.TempName, v).Find(&CodeWhiteList2).Count(&total).Error
			if err != nil {
				return err
			}
			if total > 0 {
				err = tx.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? AND user_code = ? ", info.ProCode, info.TempName, v).Updates(map[string]interface{}{
					"block_permissions": CodeWhiteList.BlockPermissions,
				}).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			} else {
				var newCodeWhiteList model.SysWhiteList
				newCodeWhiteList.UserCode = v
				newCodeWhiteList.CreatedAt = time.Now()
				newCodeWhiteList.UpdatedAt = time.Now()
				newCodeWhiteList.BlockPermissions = CodeWhiteList.BlockPermissions
				newCodeWhiteList.ProCode = CodeWhiteList.ProCode
				newCodeWhiteList.ProName = CodeWhiteList.ProName
				newCodeWhiteList.TempName = CodeWhiteList.TempName
				var c model.SysUser
				_ = global.GDb.Model(&model.SysUser{}).Where("code = ? ", v).Find(&c).Error
				newCodeWhiteList.UserName = c.NickName
				err = tx.Model(&model.SysWhiteList{}).Create(&newCodeWhiteList).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		} else {
			var total int64
			var CodeWhiteList2 model.SysWhiteList
			err = tx.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? AND user_name = ? ", info.ProCode, info.TempName, v).Find(&CodeWhiteList2).Count(&total).Error
			if err != nil {
				return err
			}
			if total > 0 {
				err = tx.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? AND user_name = ? ", info.ProCode, info.TempName, v).Updates(map[string]interface{}{
					"block_permissions": CodeWhiteList.BlockPermissions,
				}).Error
				if err != nil {
					return err
				}
			} else {
				var newCodeWhiteList model.SysWhiteList
				newCodeWhiteList.UserName = v
				newCodeWhiteList.CreatedAt = time.Now()
				newCodeWhiteList.UpdatedAt = time.Now()
				newCodeWhiteList.BlockPermissions = CodeWhiteList.BlockPermissions
				newCodeWhiteList.ProCode = CodeWhiteList.ProCode
				newCodeWhiteList.ProName = CodeWhiteList.ProName
				newCodeWhiteList.TempName = CodeWhiteList.TempName
				var c model.SysUser
				_ = global.GDb.Model(&model.SysUser{}).Where("nick_name = ? ", v).Find(&c).Error
				newCodeWhiteList.UserCode = c.Code
				err = tx.Model(&model.SysWhiteList{}).Create(&newCodeWhiteList).Error
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return err
}

func ExportWhileList(info request.ExportWhileList) (err error, path string) {
	err, sysTempB := GetBlockName(info.ProCode, info.TempName)
	if err != nil {
		return err, ""
	}
	var whiteList []model.SysWhiteList
	var total int64
	db := global.ProDbMap[info.ProCode]
	if db == nil {
		return global.ProDbErr, ""
	}
	err = db.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? ", info.ProCode, info.TempName).Find(&whiteList).Count(&total).Error
	fmt.Println("ExportWhileList", total)
	if err != nil {
		return err, ""
	}
	bookName := info.ProCode + "-" + info.TempName + "-白名单"
	excelLine := make([][]interface{}, 1)
	bm := make(map[string][]string, 0)
	//表头
	top := make([]interface{}, 0)
	top2 := make([]string, 0)
	for _, v := range sysTempB {
		top = append(top, v.Name)
		top2 = append(top2, v.Name)
		bm[v.Name] = make([]string, 0)
	}
	excelLine[0] = top
	for _, list := range whiteList {
		if len(list.BlockPermissions) != 0 {
			for _, v := range list.BlockPermissions {
				for _, v1 := range sysTempB {
					if v == v1.Code {
						bm[v1.Name] = append(bm[v1.Name], list.UserCode+list.UserName)
					}
				}
			}
		}
	}

	maxlen := 0
	for _, v := range bm {
		if len(v) > maxlen {
			maxlen = len(v)
		}
	}
	item := make([][]interface{}, maxlen+1)
	for i := 0; i <= maxlen; i++ {
		item[i] = make([]interface{}, len(bm)+1)
	}

	for i, v := range bm {
		for i1, v1 := range top2 {
			if v1 == i {
				for i2, v2 := range v {
					item[i2][i1] = v2
				}
			}
		}
	}
	excelLine = append(excelLine, item...)

	fmt.Println("ExportWhileList", excelLine)
	err = u.ExportExcelTheMainEntrance(excelLine, bookName, info.ProCode, "white-list/")
	return err, "files/white-list/" + info.ProCode + "/" + bookName + ".xlsx"
}

func GetBlockName(proCode, tempName string) (error, []m.SysProTempB) {
	var pro m.SysProject
	err := global.GDb.Model(&m.SysProject{}).Where("code = ? ", proCode).Find(&pro).Error
	if err != nil {
		return err, nil
	}

	var sysProTemplate m.SysProTemplate
	err = global.GDb.Model(&m.SysProTemplate{}).Where("pro_id = ? AND name = ? ", pro.ID, tempName).Find(&sysProTemplate).Error
	if err != nil {
		return err, nil
	}

	var sysProTemplateB []m.SysProTempB
	err = global.GDb.Model(&m.SysProTempB{}).Where("pro_temp_id = ? ", sysProTemplate.ID).Find(&sysProTemplateB).Error
	if err != nil {
		return err, nil
	}
	return nil, sysProTemplateB
}

func GetBlockPeopleSum(proCode, tempName string) (error, interface{}) {
	err, sysProTemp := GetBlockName(proCode, tempName)
	if err != nil {
		return err, nil
	}
	SysWhiteListPeopleSum := model.SysWhiteListPeopleSum{BlockSum: make(map[string]string, len(sysProTemp))}

	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr, nil
	}
	var whileList []model.SysWhiteList
	err = db.Model(&model.SysWhiteList{}).Where("pro_code = ? AND temp_name = ? ", proCode, tempName).Find(&whileList).Error
	if err != nil {
		return err, nil
	}
	b := make(map[string]int, 0)
	for _, v := range sysProTemp {
		SysWhiteListPeopleSum.BlockSum[v.Code+"-"+v.Name] = "0"
		b[v.Code] = 0
	}
	for _, v := range sysProTemp {
		for _, v1 := range whileList {
			for _, v2 := range v1.BlockPermissions {
				if v.Code == v2 {
					b[v.Code]++
				}
			}
		}
	}
	for _, v := range sysProTemp {
		sum := b[v.Code]
		SysWhiteListPeopleSum.BlockSum[v.Code+"-"+v.Name] = strconv.Itoa(sum)
	}
	return nil, SysWhiteListPeopleSum
}
