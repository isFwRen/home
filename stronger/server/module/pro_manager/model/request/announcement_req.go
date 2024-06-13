package request

import (
	"server/module/sys_base/model"
	"server/module/sys_base/model/request"
)

type AnnouncementPageReq struct {
	//Title           string    `json:"title"`
	model.BasePageInfo
	model.BaseTimeRange
	ReleaseType int    `json:"releaseType" form:"releaseType"`
	ProCode     string `json:"proCode" form:"proCode"`
	Status      int    `json:"status" form:"status"`
	//ReleaseDate     time.Time `json:"releaseDate"`
	//ReleaseUserCode string    `json:"releaseUserCode"`
	//VisitCount      int       `json:"visitCount"`
	//Content         string    `json:"content"`
	//ReleaseUserName string    `json:"releaseUserName"`
}

type AnnouncementChangeStatusReq struct {
	request.ReqIds
	Status int `json:"status" binding:"required"` //删除传：3，发布传：2，取消发布和恢复传：1
}

type AnnouncementPageHomeReq struct {
	model.BasePageInfo
	model.BaseTimeRange
	Title       string `json:"title" form:"title"`
	ProCode     string `json:"proCode" form:"proCode"`
	ReleaseType int    `json:"releaseType" form:"releaseType"`
	Status      int    `json:"status" form:"status"`
}
