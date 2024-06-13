/**
 * @Author: 星期一
 * @Description:
 * @Date: 2021/1/4 下午4:01
 */

package model

import (
	"github.com/lib/pq"
	"server/module/sys_base/model"
	"time"
)

type SysQuality struct {
	model.Model
	//ParentXmlNodeId   string    `json:"parentXmlNodeId" form:"parentXmlNodeId" gorm:"大导出节点id"`
	ParentXmlNodeName string         `json:"parentXmlNodeName" form:"parentXmlNodeName" gorm:"大导出节点name"`
	XmlNodeVal        pq.StringArray `json:"xmlNodeVal" form:"xmlNodeVal" gorm:"type:varchar(100)[];comment:'导出节点的值'"`
	XmlNodeName       string         `json:"xmlNodeName" form:"xmlNodeName" gorm:"导出节点name"`
	FieldName         string         `json:"fieldName" form:"fieldName" gorm:"字段名称"`
	FieldCode         string         `json:"fieldCode" form:"fieldCode" gorm:"字段编码"`
	InputType         int            `json:"inputType" form:"inputType" gorm:"输入方式(1:输入框,2:常量,3:下拉,)"`
	BelongType        int            `json:"belongType" form:"belongType" gorm:"所属信息类型（1:申请人信息,2:被保人信息,3:受托人信息,4:其他信息,5:受益人信息,6:领款人信息,7:账单信息,8:出险信息）"`
	BillInfo          int            `json:"billInfo" form:"billInfo" gorm:"账单信息（1：账单号，等）"`
	Beneficiary       int            `json:"beneficiary" form:"beneficiary" gorm:"受益人（1：受益人姓名，2：领款人姓名，3：其他）"`
	WidthPercent      int            `json:"widthPercent" form:"widthPercent" gorm:"排版(4:25%)"`
	MyOrder           int            `json:"myOrder" form:"myOrder" gorm:"排序"`
	UpdatedBy         string         `json:"updatedBy" form:"updatedBy" gorm:"更新人"`
	CreatedBy         string         `json:"createdBy" form:"createdBy" gorm:"创建人"`
	ProId             string         `json:"proId" form:"proId" gorm:"项目id"`
	ProName           string         `json:"proName" form:"proName" gorm:"项目名称"`
}

type SysQualityOption struct {
	Id                string    `json:"id" gorm:"唯一自增id"`
	ParentXmlNodeId   string    `json:"parentXmlNodeId" gorm:"大导出节点id"`
	ParentXmlNodeName string    `json:"parentXmlNodeName" gorm:"大导出节点name"`
	XmlNodeId         string    `json:"xmlNodeId" gorm:"导出节点id"`
	XmlNodeName       string    `json:"xmlNodeName" gorm:"导出节点name"`
	FieldName         string    `json:"fieldName" gorm:"字段名称"`
	InputType         int       `json:"inputType" gorm:"输入方式(1:输入框,2:常量,3:下拉,)"`
	BelongType        int       `json:"belongType" gorm:"所属信息类型（1:申请人信息,2:被保人信息,3:受托人信息,4:其他信息,5:受益人信息,6:领款人信息,7:账单信息,8:出险信息）"`
	BillInfo          int       `json:"billInfo" gorm:"账单信息（1：账单号，等）"`
	Beneficiary       int       `json:"beneficiary" gorm:"受益人（1：受益人姓名，2：领款人姓名，3：其他）"`
	WidthPercent      int       `json:"widthPercent" gorm:"排版(4:25%)"`
	MyOrder           int       `json:"myOrder" gorm:"排序"`
	CreatedAt         time.Time `json:"createdAt" gorm:"创建时间"`
	UpdatedAt         time.Time `json:"updatedAt" gorm:"更新时间"`
	UpdatedBy         string    `json:"updatedBy" gorm:"更新人"`
	CreatedBy         string    `json:"createdBy" gorm:"创建人"`
	ProId             string    `json:"proId" gorm:"项目id"`
	ProName           string    `json:"proName" gorm:"项目名称"`
	BelongTypeName    string    `json:"belongTypeName" gorm:"所属信息类型"`
	InputTypeName     string    `json:"inputTypeName" gorm:"输入方式"`
	BillInfoName      string    `json:"billInfoName" gorm:"账单信息"`
}
