/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 3:05 下午
 */

package model

import (
	"server/module/sys_base/model"
	//  "time"
	//  "github.com/lib/pq"
)

type SysSpotCheckWrong struct {
	model.Model
	CreatedCode string `json:"createdCode" from:"createdCode"`
	CreatedName string `json:"createdName" from:"createdName"`
	ProCode     string `json:"pro" from:"proCode"`
	Code        string `json:"cdoe" from:"cdoe"`
	Name        string `json:"name" from:"name"`
	Type        int    `json:"type" from:"type"`
	BillNum     string `json:"billNum" from:"billNum"`
	Path        string `json:"path" from:"path"`
	Picture     string `json:"picture" from:"picture"`
	Wrong       string `json:"wrong" from:"wrong"`
	Right       string `json:"right" from:"right"`
	Block       string `json:"block" from:"block"`
	BlockId     string `json:"blockId" from:"blockId"`
	FieldCode   string `json:"fieldCode" from:"fieldCode"`
	FieldName   string `json:"fieldName" from:"fieldName"`
	//  DeadlineUploadTime time.Time      `json:"deadlineUploadTime" from:"deadlineUploadTime"`                           //图片加密秘钥
	//  Images             pq.StringArray `json:"images" form:"images" gorm:"type:varchar(100)[] comment:'一系列校验的数组'"`     // 原始图片

}
