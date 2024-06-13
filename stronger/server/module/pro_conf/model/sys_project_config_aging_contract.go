package model

import "server/module/sys_base/model"

type SysProjectConfigAgingContract struct {
	model.Model
	ProId                    string `json:"proId" binding:"required" gorm:"comment:'项目ID'"`
	ContractStartTime        string `json:"contractStartTime" binding:"required" gorm:"comment:'时效起始时间'"`
	ContractEndTime          string `json:"contractEndTime" binding:"required" gorm:"comment:'时效结束时间'"`
	ClaimType                int    `json:"claimType" binding:"required" gorm:"comment:'理赔类型'"`
	ContractOutsideStartTime string `json:"contractOutsideStartTime" binding:"required" gorm:"comment:'时效外开始时间'"`
	ContractOutsideEndTime   string `json:"contractOutsideEndTime" binding:"required" gorm:"comment:'时效外最晚时间'"`
	RequirementsTime         string `json:"requirementsTime" binding:"required" gorm:"comment:'考核要求(min)'"`
	Code                     string `json:"code" binding:"required" gorm:"comment:'项目编码'"`
}

type SysProjectConfigAgingContractReq struct {
	model.Model
	ProId                    string `json:"proId" binding:"required" gorm:"comment:'项目编码'"`
	ContractStartTime        string `json:"contractStartTime" binding:"required" gorm:"comment:'时效起始时间'"`
	ContractEndTime          string `json:"contractEndTime" binding:"required" gorm:"comment:'时效结束时间'"`
	ClaimType                string `json:"claimType" binding:"required" gorm:"comment:'理赔类型'"`
	ContractOutsideStartTime string `json:"contractOutsideStartTime" binding:"required" gorm:"comment:'时效外开始时间'"`
	ContractOutsideEndTime   string `json:"contractOutsideEndTime" binding:"required" gorm:"comment:'时效外最晚时间'"`
	RequirementsTime         string `json:"requirementsTime" binding:"required" gorm:"comment:'考核要求(min)'"`
	Code                     string `json:"code" binding:"required" gorm:"comment:'项目编码'"`
}
