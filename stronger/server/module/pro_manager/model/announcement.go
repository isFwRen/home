/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/7/7 11:19
 */

package model

import (
	modelBase "server/module/sys_base/model"
	"time"
)

type Announcement struct {
	modelBase.Model
	Title           string    `json:"title"`
	ReleaseType     int       `json:"releaseType"`
	ProCode         string    `json:"proCode"`
	Status          int       `json:"status"`
	ReleaseDate     time.Time `json:"releaseDate"`
	ReleaseUserCode string    `json:"releaseUserCode"`
	VisitCount      int64     `json:"visitCount"`
	Content         string    `json:"content"`
	ReleaseUserName string    `json:"releaseUserName"`
}
