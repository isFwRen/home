package schedule

import (
	"runtime/debug"
	"server/global"
	model2 "server/module/load/model"
	"server/module/pro_manager/model"
	"server/module/pro_manager/project/B0108"
	"server/utils"
	"time"

	"go.uber.org/zap"
)

func UpdateTimeoutStage() {
	defer func() {
		if err := recover(); err != nil {
			global.GLog.Error("", zap.Any("", err))
			global.GLog.Error(string(debug.Stack()))
		}
	}()
	if !utils.HasItem(global.GConfig.System.ProArr, "B0108") {
		return
	}
	global.GLog.Info("录入中超时单据放优先单池", zap.Any("时间", time.Now()))
	db := global.ProDbMap["B0108"]
	if db == nil {
		global.GLog.Error("B0108", zap.Error(global.ProDbErr))
		return
	}
	var bills []model.ProjectBill
	db.Model(&model.ProjectBill{}).
		Where("stage = 2 and stick_level <> 2 and scan_at BETWEEN ? AND ? ",
			time.Now().Add(-24*60*time.Minute), time.Now()).Limit(100).Find(&bills)
	global.GLog.Info("录入中超时单据放优先单池", zap.Any("bills", len(bills)))
	for _, bill := range bills {
		backAtTheLatestStr, _, second := B0108.CalculateBackTimeAndTimeRemaining(bill)
		global.GLog.Info(bill.BillName, zap.Any("second", second))
		global.GLog.Info(bill.BillName, zap.Any("backAtTheLatestStr", backAtTheLatestStr))

		// lastTwoNum := bill.Agency[len(bill.Agency)-2 : len(bill.Agency)-1]
		if (bill.SaleChannel == "秒赔" && ((utils.RegIsMatch(`^(00083000|00083002|00083300|00083301)$`, bill.Agency) && second <= 960) || (!utils.RegIsMatch(`^(00083000|00083002|00083300|00083301)$`, bill.Agency) && second <= 360))) ||
			bill.SaleChannel == "理赔" && ((utils.RegIsMatch(`^(00083000|00083002|00083300|00083301)$`, bill.Agency) && second <= 360) || (!utils.RegIsMatch(`^(00083000|00083002|00083300|00083301)$`, bill.Agency) && second <= -960)) {
			dbTask := global.ProDbMap["B0108_task"]
			if dbTask == nil {
				global.GLog.Error("B0108_task", zap.Error(global.ProDbErr))
				return
			}
			txTask := dbTask.Begin()
			//更新任务库单据优先级
			err := txTask.Model(&model.ProjectBill{}).
				Where("stage = 2 and stick_level <> 2 and id = ?", bill.ID).
				Updates(&map[string]interface{}{
					"stick_level": 2,
				}).Error
			if err != nil {
				global.GLog.Error("B0108_task bill", zap.Error(err))
				txTask.Rollback()
				continue
			}
			//更新任务库分块优先级
			err = txTask.Model(&model2.ProjectBlock{}).
				Where("level <> 2 and bill_id = ?", bill.ID).
				Updates(&map[string]interface{}{
					"level": 2,
				}).Error
			if err != nil {
				txTask.Rollback()
				global.GLog.Error("B0108_task block", zap.Error(err))
				continue
			}
			err = txTask.Commit().Error
			if err != nil {
				global.GLog.Error("B0108_task commit", zap.Error(err))
				continue
			}

			//更新历史库单据优先级
			db.Model(&model.ProjectBill{}).
				Where("stage = 2 and stick_level <> 2 and id = ?", bill.ID).
				Updates(&map[string]interface{}{
					"stick_level": 2,
				})
			global.GLog.Info("设置优先单成功", zap.Any(bill.ID, bill.BillName))
		}
	}
}
