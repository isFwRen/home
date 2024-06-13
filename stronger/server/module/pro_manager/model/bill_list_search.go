/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 4:04 下午
 */

package model

import (
	"server/module/sys_base/model"
	"time"
)

type BillListSearch struct {
	model.BasePageInfo
	ProCode       string    `json:"proCode" form:"proCode"`                  //项目代码
	TimeStart     time.Time `json:"timeStart" form:"timeStart"`              //开始时间
	TimeEnd       time.Time `json:"timeEnd" form:"timeEnd"`                  //结束时间
	BillCode      string    `json:"billCode" form:"billCode"`                //案件号
	Status        int       `json:"status" form:"status,default=-1"`         //案件状态
	SaleChannel   string    `json:"saleChannel" form:"saleChannel"`          //销售渠道
	BatchNum      string    `json:"batchNum" form:"batchNum"`                //批次号
	Agency        string    `json:"agency" form:"agency"`                    //机构号
	InsuranceType string    `json:"insuranceType" form:"insuranceType"`      //医保类型
	ClaimType     int       `json:"claimType" form:"claimType,default=-1"`   //理赔类型
	StickLevel    int       `json:"stickLevel" form:"stickLevel,default=-1"` //加急件
	MinCountMoney float32   `json:"minCountMoney" form:"minCountMoney"`      //最小账单金额
	MaxCountMoney float32   `json:"maxCountMoney" form:"maxCountMoney"`      //最大账单金额
	IsQuestion    int       `json:"isQuestion" form:"isQuestion,default=-1"` //是否是问题件0:不筛选，1：true，2：false
	InvoiceNum    int       `json:"invoiceNum" form:"invoiceNum,default=-1"` //发票数量
	QualityUser   string    `json:"qualityUser" form:"qualityUser"`          //质检人
	Stage         string    `json:"stage" form:"stage"`                      //录入状态
}

type DelByIdAndProCode struct {
	ID         string `form:"id" binding:"required"`
	ProCode    string `json:"proCode" form:"proCode" binding:"required"` //项目代码
	DelRemarks string `json:"delRemarks" from:"delRemarks"`              //删除备注
}

type ProCodeAndId struct {
	ID      string `json:"id" form:"id" binding:"required"`
	ProCode string `json:"proCode" form:"proCode" binding:"required"` //项目代码
	FieldId string `json:"fieldId" form:"fieldId"`                    //字段id
}

type Remark struct {
	ID          string `json:"id" form:"id" binding:"required"`                //单据ID
	ProCode     string `json:"proCode" form:"proCode" binding:"required"`      //项目代码
	Remark      string `json:"remark" form:"remark"`                           //备注
	EditVersion int    `json:"editVersion" form:"editVersion" binding:"min=0"` //版本
}

type AutoUpload struct {
	ID           string `form:"id" binding:"required"`
	ProCode      string `json:"proCode" form:"proCode" binding:"required"` //项目代码
	IsAutoUpload bool   `json:"isAutoUpload" form:"isAutoUpload"`          //自动回传
}

type QingDanForm struct {
	ID         string `form:"id" binding:"required"`
	ProCode    string `json:"proCode" form:"proCode" binding:"required"` //项目代码
	ItemName   string `json:"itemName" form:"itemName"`                  //项目名称
	InvoiceNum string `json:"invoiceNum" form:"invoiceNum"`              //账单号
}

type DeductionDetailSearch struct {
	model.BasePageInfo
	ProCode   string    `json:"proCode" form:"proCode"`     //项目代码
	TimeStart time.Time `json:"timeStart" form:"timeStart"` //开始时间
	TimeEnd   time.Time `json:"timeEnd" form:"timeEnd"`     //结束时间
	BillCode  string    `json:"billCode" form:"billCode"`   //案件号
	Name      string    `json:"name" form:"name"`           //案件状态
}

type HistoryDelProCode struct {
	ProCode string `json:"proCode" form:"proCode" binding:"required"` //项目代码
}
