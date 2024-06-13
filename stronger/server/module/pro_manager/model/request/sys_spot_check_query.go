package request

import (
	"server/module/sys_base/model"
	// "server/module/sys_base/model/request"
)

type SysSpotCheckQuery struct {
	//Title           string    `json:"title"`
	model.BasePageInfo
	model.BaseTimeRange
	// ReleaseType int    `json:"releaseType" form:"releaseType"`
	ProCode string `json:"proCode" form:"proCode"`
	Status  int    `json:"status" form:"status"`
	Type    int    `json:"type" form:"type"`
	//ReleaseDate     time.Time `json:"releaseDate"`
}

type SysSpotCheckDataQuery struct {
	//Title           string    `json:"title"`
	model.BasePageInfo
	model.BaseTimeRange
	// ReleaseType int    `json:"releaseType" form:"releaseType"`
	ProCode string `json:"proCode" form:"proCode"`
	Code    int    `json:"Code" form:"Code"`
	Name    int    `json:"Name" form:"Name"`
	Type    int    `json:"type" form:"type"`
	//ReleaseDate     time.Time `json:"releaseDate"`
}

type SysSpotCheckWrongQuery struct {
	//Title           string    `json:"title"`
	model.BasePageInfo
	model.BaseTimeRange
	// ReleaseType int    `json:"releaseType" form:"releaseType"`
	ProCode     string `json:"proCode" form:"proCode"`
	Code        string `json:"code" form:"code"`
	Name        string `json:"name" form:"name"`
	Type        string `json:"type" form:"type"`
	CreatedCode string `json:"createdCode" form:"createdCode"`
	CreatedName string `json:"createdName" form:"createdName"`
	FieldName   string `json:"fieldName" form:"fieldName"`
	//ReleaseDate     time.Time `json:"releaseDate"`
}

type SysSpotCheckStatisticQuery struct {
	//Title           string    `json:"title"`
	model.BasePageInfo
	model.BaseTimeRange
	// ReleaseType int    `json:"releaseType" form:"releaseType"`
	ProCode string `json:"proCode" form:"proCode"`
	Type    int    `json:"type" form:"type"`
	//ReleaseDate     time.Time `json:"releaseDate"`
}
