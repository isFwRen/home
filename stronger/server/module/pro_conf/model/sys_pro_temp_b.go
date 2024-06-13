package model

import (
	"server/module/sys_base/model"

	"github.com/lib/pq"
)

type SysProTempB struct {
	model.Model
	MyOrder            int                 `json:"myOrder" form:"myOrder" gorm:"模板分块排序"`
	Name               string              `json:"name" form:"name" gorm:"模板分块名字"`
	Code               string              `json:"code" form:"code" gorm:"模板分块编码"`
	FEight             bool                `json:"fEight" form:"fEight" gorm:"是否f8提交"`
	Ocr                string              `json:"ocr" form:"ocr" gorm:"ocr流程"`
	FormKey            string              `json:"formKey" form:"formKey" gorm:"初审哪个分块MB001-bc001，像外键,废弃"`
	FreeTime           int                 `json:"freeTime" form:"freeTime" gorm:"释放时间"`
	IsLoop             bool                `json:"isLoop" form:"isLoop" gorm:"是否循环分块"`
	PicPage            int                 `json:"picPage" form:"picPage" gorm:"图片页码，第几张图片"`
	IsMobile           bool                `json:"isMobile" form:"isMobile" gorm:"移动端可录入"`
	WCoordinate        pq.StringArray      `json:"wCoordinate" form:"wCoordinate" gorm:"type:varchar(100)[];comment:'web截图位置'"`
	MCoordinate        pq.StringArray      `json:"mCoordinate" form:"mCoordinate" gorm:"type:varchar(100)[];comment:'移动端截图位置'"`
	ProTempId          string              `json:"proTempId" form:"proTempId" gorm:"所属模板id"`
	MPicPage           int                 `json:"mPicPage" form:"mPicPage" gorm:"移动端图片页码，第几张图片"`
	Relation           string              `json:"relation" form:"relation" gorm:"关联,初审哪个分块MB001-bc001，像外键"`
	IsCompetitive      bool                `json:"isCompetitive" form:"isCompetitive" gorm:"流程(true:通用，一二码可抢)"`
	TempBFRelations    []TempBFRelation    `json:"tempBFRelations" form:"tempBFRelations" gorm:"foreignKey:TempBId;references:id;comment:'模板配置'"`
	TempBlockRelations []TempBlockRelation `json:"tempBlockRelations" form:"tempBlockRelations" gorm:"foreignKey:TempBId;references:id;comment:'模板配置'"`
}
