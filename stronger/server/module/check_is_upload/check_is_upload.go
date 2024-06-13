/**
 * @Author: xingqiyi
 * @Description:
 * @Date: 2022/10/20 10:28
 */

package check_is_upload

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"server/global"
	"server/module/check_is_upload/service"
	"time"
)

const dayRange = 30 //检查单据天数范围

// CheckIsUpload 检查客户是否收到
func CheckIsUpload(proCode string) error {
	//获取时间范围
	startTime := time.Now().Add(dayRange * -24 * 60 * time.Minute)

	//获取时间范围的单据
	err, bills := service.GetCheckIsUploadHistoryBills(proCode, startTime)
	if err != nil {
		return err
	}
	if len(bills) == 0 {
		return errors.New("没有符合的单据")
	}

	for _, bill := range bills {
		//检查是否收到
		err, isUploadSuccess := CheckIsUploadAdapter(proCode, bill)
		if err != nil {
			global.GLog.Error("检查客户是否收到", zap.Error(err))
			continue
		}
		global.GLog.Info("bill  " + bill.BatchNum)
		fmt.Println(isUploadSuccess)
		if isUploadSuccess {
			//更新单据状态
			err = service.UpdateBillStage(proCode, bill)
			if err != nil {
				global.GLog.Error("更新单据状态", zap.Error(err))
			}
			global.GLog.Info("bill  " + bill.BillNum + "更新成功")
		}
	}

	return nil
}
