package model

import (
	"time"
)

type ProjectBill struct {
	ID          string    `json:"id" gorm:"primary_key"`
	Name        string    `json:"name" gorm:"comment:'字段名字1'"`
	Code        string    `json:"code" gorm:"comment:'字段名字2'"`
	CreateAt    time.Time `json:"create_at" gorm:"comment:'字段名字3'"`
	ProjectCode string    `json:"project_code" gorm:"comment:'字段名字4'"`
	Stage       string    `json:"stage" gorm:"comment:'字段名字5'"`
	Path        string    `json:"path" gorm:"comment:'字段名字6'"`
	DownloadAt  time.Time `json:"download_at" gorm:"comment:'字段名字7'"`
	Picture     []string  `json:"picture" gorm:"comment:'字段名字8'"`
	BillType    string    `json:"billType" gorm:"comment:'字段名字6'"`
	// blocks     []block
}
