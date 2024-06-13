package model

import "time"

// Model base model definition, including fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`, which could be embedded in your models
//
//	type User struct {
//	  gorm.Model
//	}
type Model struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt *time.Time `sql:"index"`
}

type BasePageInfo struct {
	PageIndex int    `json:"pageIndex" form:"pageIndex,default=1" binding:"required,min=1"`
	PageSize  int    `json:"pageSize" form:"pageSize,default=10" binding:"required,max=1500,min=1"`
	OrderBy   string `json:"orderBy" form:"orderBy"` //排序 JSON.stringify([["CreatedAt","desc"]])
}

type BaseTimeRange struct {
	StartTime time.Time `json:"startTime" form:"startTime"` //开始时间
	EndTime   time.Time `json:"endTime" form:"endTime"`     //结束时间
}

type BaseTimeRangeWithCode struct {
	ProCode string `json:"proCode" form:"proCode"` //项目编码
	BaseTimeRange
}

type BaseTimePageCode struct {
	BasePageInfo
	BaseTimeRangeWithCode
}

type AutoAddIdModel struct {
	ID        int `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
