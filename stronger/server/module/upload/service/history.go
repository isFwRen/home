/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2021/11/30 4:28 下午
 */

package service

import (
	"server/global"
	"server/module/pro_manager/model"
	"strings"
	"time"
)

//GetIsAutoUploadBills 获取可以自动回传的单据
func GetIsAutoUploadBills(proCode string) (err error, bills []model.ProjectBill) {
	db, ok := global.ProDbMap[proCode]
	if !ok {
		return global.ProDbErr, bills
	}
	err = db.Model(&model.ProjectBill{}).Where("is_auto_upload = true and (stage = 3 or stage = 4)").
		Order("created_at desc").Limit(20).Find(&bills).Error

	return err, bills
}

func UpdateStageBill(req model.ProCodeAndId) (err error) {
	db, ok := global.ProDbMap[req.ProCode]
	if !ok {
		return global.ProDbErr
	}
	var updateMap = map[string]interface{}{
		"stage": 5,
	}
	err = db.Model(&model.ProjectBill{}).Where("id = ?", req.ID).
		Updates(updateMap).Error
	return err
}

//UpdateBill 回传更新单据
func UpdateBill(req model.ProCodeAndId, uploadTime time.Time) (err error) {
	db, ok := global.ProDbMap[req.ProCode]
	if !ok {
		return global.ProDbErr
	}
	var firstUploadTime time.Time
	err = db.Model(&model.ProjectBill{}).Where("id = ?", req.ID).Select("upload_at").Find(&firstUploadTime).Error
	if err != nil {
		return err
	}
	var updateMap = map[string]interface{}{
		"stage":          5,
		"last_upload_at": uploadTime,
	}
	if strings.Index(firstUploadTime.String(), "0001-01-01") == 0 {
		updateMap["upload_at"] = uploadTime
	}
	err = db.Model(&model.ProjectBill{}).Where("id = ?", req.ID).
		Updates(updateMap).Error
	return err
}

//UpdateBill 回传更新单据
func UpdateBillByBatch(proCode, batch_num string, uploadTime time.Time) (err error) {
	db, ok := global.ProDbMap[proCode]
	if !ok {
		return global.ProDbErr
	}
	var updateMap = map[string]interface{}{
		"stage":          5,
		"last_upload_at": uploadTime,
		"upload_at":      uploadTime,
	}
	err = db.Model(&model.ProjectBill{}).Where("batch_num = ?", batch_num).
		Updates(updateMap).Error
	return err
}

func GetNumByBatchNum(req model.ProCodeAndId, batchNum string) (err error, bill int64) {
	db := global.ProDbMap[req.ProCode]
	if db == nil {
		return global.ProDbErr, bill
	}
	err = db.Model(&model.ProjectBill{}).Where("batch_num = ? and id != ? and stage = 3", batchNum, req.ID).Count(&bill).
		Error
	return err, bill
}
