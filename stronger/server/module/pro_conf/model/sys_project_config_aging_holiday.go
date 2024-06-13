package model

import "server/module/sys_base/model"

type ProjectConfigAgingHolidayArr struct {
	ProjectConfigAgingHoliday []SysProjectConfigAgingHoliday `json:"project_config_aging_holiday"`
}

type SysProjectConfigAgingHoliday struct {
	//model.AutoAddIdModel
	model.Model
	Date  string `json:"date" gorm:"comment:'日期'"`
	Day1  bool   `json:"day1" gorm:"comment:'1号'"`
	Day2  bool   `json:"day2" gorm:"comment:'2号'"`
	Day3  bool   `json:"day3" gorm:"comment:'3号'"`
	Day4  bool   `json:"day4" gorm:"comment:'4号'"`
	Day5  bool   `json:"day5" gorm:"comment:'5号'"`
	Day6  bool   `json:"day6" gorm:"comment:'6号'"`
	Day7  bool   `json:"day7" gorm:"comment:'7号'"`
	Day8  bool   `json:"day8" gorm:"comment:'8号'"`
	Day9  bool   `json:"day9" gorm:"comment:'9号'"`
	Day10 bool   `json:"day10" gorm:"comment:'10号'"`
	Day11 bool   `json:"day11" gorm:"comment:'11号'"`
	Day12 bool   `json:"day12" gorm:"comment:'12号'"`
	Day13 bool   `json:"day13" gorm:"comment:'13号'"`
	Day14 bool   `json:"day14" gorm:"comment:'14号'"`
	Day15 bool   `json:"day15" gorm:"comment:'15号'"`
	Day16 bool   `json:"day16" gorm:"comment:'16号'"`
	Day17 bool   `json:"day17" gorm:"comment:'17号'"`
	Day18 bool   `json:"day18" gorm:"comment:'18号'"`
	Day19 bool   `json:"day19" gorm:"comment:'19号'"`
	Day20 bool   `json:"day20" gorm:"comment:'20号'"`
	Day21 bool   `json:"day21" gorm:"comment:'21号'"`
	Day22 bool   `json:"day22" gorm:"comment:'22号'"`
	Day23 bool   `json:"day23" gorm:"comment:'23号'"`
	Day24 bool   `json:"day24" gorm:"comment:'24号'"`
	Day25 bool   `json:"day25" gorm:"comment:'25号'"`
	Day26 bool   `json:"day26" gorm:"comment:'26号'"`
	Day27 bool   `json:"day27" gorm:"comment:'27号'"`
	Day28 bool   `json:"day28" gorm:"comment:'28号'"`
	Day29 bool   `json:"day29" gorm:"comment:'29号'"`
	Day30 bool   `json:"day30" gorm:"comment:'30号'"`
	Day31 bool   `json:"day31" gorm:"comment:'31号'"`
}
