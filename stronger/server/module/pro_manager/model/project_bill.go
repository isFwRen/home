/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/3 3:05 下午
 */

package model

import (
	"server/module/sys_base/model"
	"time"

	"github.com/lib/pq"
)

type ProjectBill struct {
	model.Model
	BillName        string    `json:"billName" from:"billName"`                       //单据号来源内部
	BillNum         string    `json:"billNum" from:"billNum"`                         //单据号来源客户
	ProCode         string    `json:"proCode" from:"proCode"`                         //项目编号
	Stage           int       `json:"stage" from:"stage"`                             //录入状态
	DownloadPath    string    `json:"downloadPath" from:"downloadPath"`               //下载路径
	DownloadAt      time.Time `json:"downloadAt" from:"downloadAt"`                   //下载时间
	BatchNum        string    `json:"batchNum" from:"batchNum"`                       //批次号
	Agency          string    `json:"agency" from:"agency"`                           //机构号
	ScanAt          time.Time `json:"scanAt" from:"scanAt"`                           //扫描时间
	ExportAt        time.Time `json:"exportAt" from:"exportAt"`                       //导出时间
	FirstExportAt   time.Time `json:"firstExportAt" from:"firstExportAt"`             //初次导出时间
	UploadAt        time.Time `json:"uploadAt" from:"uploadAt"`                       //回传时间
	LastUploadAt    time.Time `json:"lastUploadAt" from:"lastUploadAt"`               //回传时间
	Status          int       `json:"status" from:"status" dict:"true"`               //案件状态
	SaleChannel     string    `json:"saleChannel" from:"saleChannel"`                 //销售渠道
	InsuranceType   string    `json:"insuranceType" from:"insuranceType" dict:"true"` //医保类型后改单证类型
	ClaimType       int       `json:"claimType" from:"claimType" dict:"true"`         //理赔类型
	CountMoney      float64   `json:"countMoney" from:"countMoney"`                   //账单金额
	InvoiceNum      int       `json:"invoiceNum" from:"invoiceNum"`                   //发票数量
	QuestionNum     int       `json:"questionNum" from:"questionNum"`                 //结果数据字段的问题件数量
	QualityUserCode string    `json:"qualityUserCode" from:"qualityUserCode"`         //质检人
	QualityUserName string    `json:"qualityUserName" from:"qualityUserName"`         //质检人
	IsAutoUpload    bool      `json:"isAutoUpload" from:"isAutoUpload"`               //是否自动回传
	StickLevel      int       `json:"stickLevel" from:"stickLevel"`                   //加急件
	PreStatus       int       `json:"preStatus" from:"preStatus"`                     //删除前的状态
	EditVersion     int       `json:"editVersion" from:"editVersion"`                 //版本
	DelRemarks      string    `json:"delRemarks" from:"delRemarks"`                   //删除备注
	WrongNote       string    `json:"wrongNote" from:"wrongNote"`                     //导出校验
	Template        string    `json:"template" from:"template"`                       //匹配模板
	//Crypto             string         `json:"crypto" from:"crypto"`                                    //图片加密秘钥
	IsTimeout          bool           `json:"isTimeout" from:"isTimeout"`                              //图片加密秘钥
	DeadlineUploadTime time.Time      `json:"deadlineUploadTime" from:"deadlineUploadTime"`            //图片加密秘钥
	Images             pq.StringArray `json:"images" form:"images" gorm:"type:varchar(100)[]"`         // 原始图片
	Pictures           pq.StringArray `json:"pictures" form:"pictures" gorm:"type:varchar(100)[]"`     // 转换后的图片
	PackCode           string         `json:"packCode" from:"packCode"`                                //客户推送通知函节点
	Files              pq.StringArray `json:"files" form:"files" gorm:"type:varchar(100)[]"`           //下载文件数组
	OtherInfo          string         `json:"otherInfo" form:"otherInfo"`                              //下载xml内容
	ImagesType         pq.StringArray `json:"imagesType" form:"imagesType" gorm:"type:varchar(100)[]"` //图片类型
	Remark             string         `json:"remark" form:"remark"`                                    //备注                                                   //下载xml内容
	ExportStage        int            `json:"exportStage" form:"exportStage"`                          //首次导出状态
	BillType           int            `json:"billType" form:"billType"`                                //单据状态
}
