package guoshou

import (
	"go.uber.org/zap"
	"server/global"
	"server/module/download/service"
)

func DownloadImages(proCode string) {
	//查询待下载图片的单据
	err, bills := service.FetchBillByStage(proCode)
	global.GLog.Info(proCode, zap.Any("待下载图片单据", len(bills)))
	if err != nil {
		global.GLog.Error("查询待下载图片的单据", zap.Error(err))
		return
	}

	for _, b := range bills {
		global.GLog.Info(proCode, zap.Any("下载该单据图片", b.BillName))
		//下载图片
		bill, err := FetchFiles(b)
		if err != nil {
			global.GLog.Error("下载图片", zap.Error(err))
			return
		}

		//更新数据库单据流程状态
		err = service.UpdateBillStage(bill)
		if err != nil {
			global.GLog.Error("更新数据库单据流程状态", zap.Error(err))
			return
		}
	}
	return
}
