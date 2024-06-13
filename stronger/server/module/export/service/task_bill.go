/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/24 11:28 上午
 */

package service

import (
	"fmt"
	"server/global"
	"server/module/export/utils"
	model2 "server/module/load/model"
	"server/module/pro_manager/model"
	"sync"
)

//GetTaskBillList 获取所有可以导出的task单据
func GetTaskBillList(proCode string, wg *sync.WaitGroup, obj *[]model.ProjectBill) {
	defer utils.DealErr(wg)

	db := global.ProDbMap[proCode+"_task"]
	if db == nil {
		panic("没有该项目的连接")
	}
	err := db.Order("created_at desc").Where("stage = 6").Find(&obj).Error
	if err != nil {
		panic(err)
	}
}

//GetTaskBillFields 获取所有可以导出的task单据的字段
func GetTaskBillFields(bill model.ProjectBill) (err error, fields []model2.ProjectField) {
	db := global.ProDbMap[bill.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, fields
	}
	err = db.Order("block_index asc").Order("field_index asc").Where("bill_id = ?", bill.ID).Find(&fields).Error
	if err != nil {
		return err, fields
	}
	return nil, fields
}

//GetBlocks 获取所有可以导出的task单据的分块
func GetBlocks(bill model.ProjectBill) (err error, blocks []model2.ProjectBlock) {
	db := global.ProDbMap[bill.ProCode+"_task"]
	if db == nil {
		return global.ProDbErr, blocks
	}
	err = db.Order("created_at desc").Where("bill_id = ? ", bill.ID).Find(&blocks).Error
	if err != nil {
		return err, blocks
	}
	return nil, blocks
}

//DelBill 已导出删除内存库信息
func DelBill(billId string, proCode string) (err error) {
	db := global.ProDbMap[proCode]
	if db == nil {
		return global.ProDbErr
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			global.GLog.Error(fmt.Sprintf("%v", r))
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		return err
	}

	err = tx.Model(&model.ProjectBill{}).Where("id = ?", billId).Delete(model.ProjectBill{}).Error
	if err != nil {
		panic(err)
	}
	err = tx.Model(&model2.ProjectBlock{}).Where("bill_id = ?", billId).Delete(model2.ProjectBlock{}).Error
	if err != nil {
		panic(err)
	}
	err = tx.Model(&model2.ProjectField{}).Where("bill_id = ?", billId).Delete(model2.ProjectField{}).Error
	if err != nil {
		panic(err)
	}

	return tx.Commit().Error
}
