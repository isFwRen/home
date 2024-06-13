/**
 * @Author: 阿先
 * @Description:
 * @Date: 2020/11/2 13:44
 */

package service

import (
	"fmt"
	"server/global"
	"server/module/dingding/model"
	dingding2 "server/module/dingding/model/request"
)


func AddGroup(group model.DingdingGroups) (err error, groupInter model.DingdingGroups)  {
	fmt.Println(group)
	err = global.GDb.Model(group).Create(&group).Error
	return err, group
}

func SelectGroupListByPage(info dingding2.OddsDingdingGroupStruct) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.PageIndex - 1)
	db := global.GDb.Model(&model.DingdingGroups{})

	//t := reflect.TypeOf(info)
	if info.ProjectCode != "" {
		db = db.Where("project_code = ?", info.ProjectCode)
	}
	if info.GroupName != "" {
		db = db.Where("group_name = ?", info.GroupName)
	}
	if info.ProjectEnvironment != "" {
		db = db.Where("project_environment = ?", info.ProjectEnvironment)
	}
	//if _, ok := t.FieldByName("GroupName"); ok {
	//	db = db.Where("group_name = ?", info.GroupName)
	//}
	//if _, ok := t.FieldByName("ProjectEnvironment"); ok {
	//	db = db.Where("project_environment = ?", info.ProjectEnvironment)
	//}
	var groupList []model.DingdingGroups
	err = db.Count(&total).Error
	err = db.Limit(limit).Offset(offset).Find(&groupList).Error
	return err, groupList, total
}

func DelGroup(id string) (err error) {
	var group model.DingdingGroups
	err = global.GDb.Where("id = ?", id).Delete(&group).Error
	return err
}

func UpdateGroup(group model.DingdingGroups, id int) (err error) {
	err = global.GDb.Where("id = ?", id).First(&model.DingdingGroups{}).Save(&group).Error
	return err
}
