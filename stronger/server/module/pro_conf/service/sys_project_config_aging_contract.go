package service

import (
	"errors"
	"server/global"
	"server/module/pro_conf/model"
)

// InsertAgingContractConfig 添加合同时效
func InsertAgingContractConfig(contract model.SysProjectConfigAgingContract) (err error) {
	// TUDO 待补充 同个项目同个类型只有一个，需要先查询
	var contractCount int64
	global.GDb.Model(&model.SysProjectConfigAgingContract{}).Where("claim_type = ? And code = ?", contract.ClaimType, contract.Code).Count(&contractCount)
	if contractCount <= 0 {
		global.GDb.Model(&model.SysProjectConfigAgingContract{})
		db := global.GDb.Model(&model.SysProjectConfigAgingContract{})
		err = db.Create(&contract).Error
		return err
	}
	return errors.New("新增失败，该项目已存在该理赔类型。")
}

// UpdateAgingContractConfig 修改合同时效配置
func UpdateAgingContractConfig(contract model.SysProjectConfigAgingContract) (err error) {
	if contract.ContractEndTime == "00:00:00" {
		contract.ContractEndTime = "23:59:59"
	}
	if contract.ContractOutsideEndTime == "00:00:00" {
		contract.ContractOutsideEndTime = "23:59:59"
	}
	var contracts []model.SysProjectConfigAgingContract
	err = global.GDb.Model(&model.SysProjectConfigAgingContract{}).Where("code = ?", contract.Code).Find(&contracts).Error
	for _, agingContract := range contracts {
		if agingContract.ClaimType == contract.ClaimType {
			err = errors.New("理赔类型重复，请重新输入")
			return err
		}
	}
	err = global.GDb.Model(&model.SysProjectConfigAgingContract{}).Where("id = ?", contract.ID).Updates(contract).Error
	return err
}

// DelAgingContractConfig 删除合同时效配置
func DelAgingContractConfig(rmByIds []string) error {
	err := global.GDb.Delete(&[]model.SysProjectConfigAgingContract{}, "id in (?)", rmByIds).Error
	return err
}

// GetAgingContractConfig 查询合同时效配置
func GetAgingContractConfig(code string) (err error, req []model.SysProjectConfigAgingContract, total int64) {
	db := global.GDb.Model(&model.SysProjectConfigAgingContract{}).Order("created_at ASC")
	var configs []model.SysProjectConfigAgingContract
	if code != "" {
		err = db.Where("code = ? ", code).Count(&total).Find(&configs).Error
		return err, configs, total
	}
	return err, nil, 0
}
