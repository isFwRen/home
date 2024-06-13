package service

import (
	"errors"
	"github.com/jinzhu/copier"
	"server/global"
	"server/module/pro_manager/model"
	"server/module/pro_manager/model/request"
	"server/utils"
)

func GetTransactionRule(info request.GetTransactionRule) (err error, list []map[string]interface{}, total int64) {
	limit := info.PageInfo.PageSize
	offset := info.PageInfo.PageSize * (info.PageInfo.PageIndex - 1)
	db := global.ProDbMap[info.ProCode]
	if db == nil {
		return global.ProDbErr, nil, 0
	}
	db = db.Model(&model.TransactionRule{})
	if info.RuleType != "" {
		db = db.Where("rule_type LIKE ? ", "%"+info.RuleType+"%")
	}
	var inf []model.TransactionRule
	err = db.Count(&total).Limit(limit).Offset(offset).Find(&inf).Error
	if err != nil {
		return err, nil, 0
	}

	list = []map[string]interface{}{}
	for _, rule := range inf {
		list = append(list, utils.Struct2Map(rule))
	}
	return nil, list, total
}

func AddTransactionRule(add request.AddTransactionRule) error {
	db := global.ProDbMap[add.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	db = db.Model(&model.TransactionRule{})
	var total int64
	err := db.Where("pro_code = ? AND rule_name = ? AND rule_type = ? ", add.ProCode, add.RuleName, add.RuleType).Count(&total).Error
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("存在重复数据")
	}
	var newData model.TransactionRule
	err = copier.Copy(&newData, &add)
	if err != nil {
		return err
	}
	err = global.ProDbMap[add.ProCode].Model(&model.TransactionRule{}).Debug().Create(&newData).Error
	return err
}

func EditTransactionRule(edit request.AddTransactionRule) error {
	db := global.ProDbMap[edit.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	var newData model.TransactionRule
	err := copier.Copy(&newData, &edit)
	if err != nil {
		return err
	}
	err = db.Model(&model.TransactionRule{}).Where("id = ? ", edit.Id).Updates(&newData).Error
	return err
}

func DeleteTransactionRule(rm request.Rm) (err error) {
	var R model.TransactionRule
	db := global.ProDbMap[rm.ProCode]
	if db == nil {
		return global.ProDbErr
	}
	for _, id := range rm.Ids {
		err = db.Where("id = ?", id).Delete(&R).Error
		if err != nil {
			return err
		}
	}
	return err
}
