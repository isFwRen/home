package service

import (
	"errors"
	"server/global"
	"server/module/pro_conf/model"
	proconf2 "server/module/pro_conf/model/request"
)

func CreateAgingConfigHoliday() (err error) {
	var SysProjectConfigAgingHoliday model.SysProjectConfigAgingHoliday
	err = global.GDb.Migrator().CreateTable(&SysProjectConfigAgingHoliday)
	return err
}

func InsertProjectConfigAgingHoliday(ProjectConfigAgingHolidayInsertDetail model.SysProjectConfigAgingHoliday) (err error) {
	err = global.GDb.Model(&model.SysProjectConfigAgingHoliday{}).Create(&ProjectConfigAgingHolidayInsertDetail).Error
	return err
}

func GetAgingConfigHoliday(configType proconf2.ProjectConfigAgingHoliday) (err error, configs []model.SysProjectConfigAgingHoliday, count int64) {
	db := global.GDb.Model(&model.SysProjectConfigAgingHoliday{})
	db = db.Where("date = ? ", configType.InquireStartDate).Order("updated_at desc").Limit(1)
	var configsRes []model.SysProjectConfigAgingHoliday
	err = db.Find(&configsRes).Error
	db.Count(&count)
	return err, configsRes, count
}

func UpdateAgingConfigHoliday(configType model.SysProjectConfigAgingHoliday) (err error) {
	//err = global.GDb.Where("id = ?", configType.ID).First(&model.SysProjectConfigAgingHoliday{}).Save(&configType).Error
	var count int64
	err = global.GDb.Model(&model.SysProjectConfigAgingHoliday{}).Where("id = ? And date = ? ", configType.ID, configType.Date).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("该" + configType.Date + "日期的节假日时效还没设置保存, 无法更新")
	}
	err = global.GDb.Model(&model.SysProjectConfigAgingHoliday{}).Where("id = ? ", configType.ID).Updates(map[string]interface{}{
		"date":  configType.Date,
		"day1":  configType.Day1,
		"day2":  configType.Day2,
		"day3":  configType.Day3,
		"day4":  configType.Day4,
		"day5":  configType.Day5,
		"day6":  configType.Day6,
		"day7":  configType.Day7,
		"day8":  configType.Day8,
		"day9":  configType.Day9,
		"day10": configType.Day10,
		"day11": configType.Day11,
		"day12": configType.Day12,
		"day13": configType.Day13,
		"day14": configType.Day14,
		"day15": configType.Day15,
		"day16": configType.Day16,
		"day17": configType.Day17,
		"day18": configType.Day18,
		"day19": configType.Day19,
		"day20": configType.Day20,
		"day21": configType.Day21,
		"day22": configType.Day22,
		"day23": configType.Day23,
		"day24": configType.Day24,
		"day25": configType.Day25,
		"day26": configType.Day26,
		"day27": configType.Day27,
		"day28": configType.Day28,
		"day29": configType.Day29,
		"day30": configType.Day30,
		"day31": configType.Day31,
	}).Error
	if err != nil {
		return err
	}
	return err
}
