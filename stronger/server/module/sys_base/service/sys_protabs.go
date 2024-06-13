package service

import (
	"errors"
	"gorm.io/gorm"
	"server/global"
	"server/module/sys_base/model"
)

// @title    SysProject
// @description   AddTabs, 新增Tabs
// @auth                     （2021年1月1日16:42:50）
// @param     sysTabs      model.TabsModel
// @return    err             error

func AddTabs(sysTabs model.TabsModel)(err error)  {
	var reSysTabs model.TabsModel
	// 判断是否有该记录
	err1 := global.GDb.Where("title = ? ", sysTabs.Title).First(&reSysTabs).Error
	//fmt.Println("result: ", reSysTabs)
	if !errors.Is(err1, gorm.ErrRecordNotFound) {
		return errors.New("已存在该记录")
	} else {
		err = global.GDb.Create(&sysTabs).Error
		if err != nil {
			return err
		}
		return err
	}
}

// @title    SysProject
// @description   GetTabsList, 获取所有Tabs
// @auth                     （2021年1月1日16:42:50）
// @param     TabsModel      model.TabsModel
// @return    err             error

func GetTabsList()(err error, list interface{}){
	var sysTabsList []model.TabsModel
	err = global.GDb.Order("id").Find(&sysTabsList).Error
	return err, sysTabsList
}

// @title    SysProject
// @description   GetTabsList, 根据name更新Tabs
// @auth                     （2021年1月1日16:42:50）
// @param     TabsModel      model.TabsModel
// @return    err             error

func UpdateTabs(tabs model.TabsModel)(rows int64){
	rows = global.GDb.Model(tabs).Where("name = ?",tabs.Name).Update("query_path",tabs.QueryPath).RowsAffected
	return rows
}

// @title    SysProject
// @description   GetTabsLast, 获取最后一条Tabs
// @auth                     （2021年1月1日16:42:50）
// @param     TabsModel      model.TabsModel
// @return    err             error

func GetTabsLast()(err error, list interface{}){
	var sysTabs model.TabsModel
	err = global.GDb.Last(&sysTabs).Error
	return err, sysTabs
}

// @title    SysProject
// @description   RemoveTabs, 根据name删除Tabs
// @auth                     （2021年1月1日16:42:50）
// @param     TabsModel      model.TabsModel
// @return    err             error

func RemoveTabs(tabs model.TabsModel) (reRows int64){
	rows := global.GDb.Where("name = ?", tabs.Name).Delete(&model.TabsModel{}).RowsAffected
	return rows
}
