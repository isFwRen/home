package request

import (
	"github.com/lib/pq"
	"server/module/pro_conf/model"
)

type SysProjectBlocks struct {
	TempId         string              `json:"tempId" gorm:"comment:'模板id'"`
	SysProTempBArr []model.SysProTempB `json:"sysProTempBArr" gorm:"comment:'项目模板分块'"`
}

type SysReqBlockCoordinate struct {
	BlockId        string         `json:"blockId" gorm:"comment:'模板id'"`
	CoordinateType bool           `json:"coordinateType" gorm:"comment:'截图类型web:true,phone:false'"`
	PicPage        int            `json:"picPage" gorm:"图片页码，第几张图片"`
	Coordinate     pq.StringArray `json:"coordinate" gorm:"type:varchar(100)[];comment:'web截图位置'"`
}

type SysBlocksRelations struct {
	BlockId              string                    `json:"blockId" gorm:"comment:'分块id'"`
	MyType               int8                      `json:"myType" gorm:"类型（前置分块:0，一码前置分块:1，参考分块:2）"`
	TempBlockRelationArr []model.TempBlockRelation `json:"tempBlockRelationArr" gorm:"comment:'项目模板分块关系'"`
}

type RqTempBlockRelation struct {
	//AutoAddIdModel
	//Name     string `json:"name" gorm:"模板分块名字"`
	MyType   int8   `json:"myType" gorm:"类型（前置分块，一码前置分块，参考分块）"`
	TempBId  string `json:"tempBId" gorm:"模板分块id"`
	PreBId   string `json:"preBId" gorm:"前置or参考分块id"`
	PreBName string `json:"preBName" gorm:"前置or参考分块名字"`
	PreBCode string `json:"preBCode" gorm:"前置or参考分块编码"`
}

type SysBFRelations struct {
	BlockId           string                 `json:"blockId" gorm:"comment:'分块id'"`
	TempBFRelationArr []model.TempBFRelation `json:"tempBFRelationArr" gorm:"comment:'项目模板分块字段关系'"`
}
